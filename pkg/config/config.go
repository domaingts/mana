package config

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/domaingts/mana/constant"
	"github.com/tidwall/gjson"
	"golang.org/x/net/http2"
)

var (
	defaultServicePath = "/etc/systemd/system"
)

type Config struct {
	user       string
	repo       string
	cmd        string
	path       string
	binaryPath string
	configPath string
	client     *http.Client
	canceled   bool
}

func NewConfig(cmd, user, repo string) *Config {
	return &Config{
		cmd:  cmd,
		user: user,
		repo: repo,
	}
}

func (c *Config) InitConfig() {
	_, err := os.Stat(c.configPath)
	if err != nil && os.IsNotExist(err) {
		err = os.Mkdir(c.configPath, 0755)
		if err != nil {
			panic(err)
		}
	}
}

func (c *Config) SetBinaryPath(path string) {
	c.binaryPath = path
}

func (c *Config) SetConfigPath(path string) {
	c.configPath = path
}

func (c *Config) CreateService(input []byte) error {
	path := fmt.Sprintf("%s/%s.service", defaultServicePath, c.cmd)
	_, err := os.Stat(path)
	if err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	file.Chmod(0755)
	_, err = file.Write(input)
	return err
}

func (c *Config) CreateStartConfig(input []byte) error {
	path := fmt.Sprintf("%s/%s", c.configPath, "config.yaml")
	_, err := os.Stat(path)
	if err == nil {
		return nil
	} else if !os.IsNotExist(err) {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	file.Chmod(0755)
	_, err = file.Write(input)
	return err
}

func (c *Config) Run() error {
	c.newClient()
	version, err := c.getLatestVersion()
	if err != nil {
		return err
	}
	body, err := c.getReader(version)
	if err != nil {
		return err
	}
	defer body.Close()
	return c.untarTargetFile(body)
}

func (c *Config) newClient() {
	c.client = &http.Client{
		Transport: &http2.Transport{},
		Timeout:   time.Minute,
	}
}

func (c *Config) getLatestVersion() (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", c.user, c.repo)
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	request.Header.Add("Accept", "application/vnd.github.v3+json")
	response, err := c.client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get version from github: %d", response.StatusCode)
	}
	result, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	version := gjson.GetBytes(result, constant.TAG_NAME).String()
	if version == "" {
		return "", fmt.Errorf("empty version: %s", version)
	}
	fmt.Println("get latest version", version)
	return version, nil
}

func (c *Config) getReader(version string) (io.ReadCloser, error) {
	getter := NewGetter(c)
	url := fmt.Sprintf("https://github.com/%s/%s/releases/download/%s/%s", c.user, c.repo, version, getter.Filename(version))
	response, err := c.client.Get(url)
	if err != nil {
		return nil, err
	}
	return response.Body, nil
}

func (c *Config) untarTargetFile(in io.Reader) error {
	gzr, err := gzip.NewReader(in)
	if err != nil {
		return err
	}
	defer gzr.Close()
	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if header.Name == c.cmd {
			target := path.Join(c.binaryPath, c.cmd)
			fmt.Println(header.Mode)
			outfile, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755)
			if err != nil {
				return err
			}
			defer outfile.Close()
			_, err = io.Copy(outfile, tr)
			if err != nil {
				return err
			}
			fmt.Println("write", header.Name, "to", target)
			return nil
		}
	}
	return errors.New("target file doesn't include in the zip file.")
}
