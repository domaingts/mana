package config

type FilenameGetter interface {
	Filename(version string) string
}

func NewGetter(c *Config) FilenameGetter {
	switch c.cmd {
	case "ddns-go":
		return &DDNSGOGetter{}
	}
	return nil
}

type DDNSGOGetter struct{}

func (d *DDNSGOGetter) Filename(_ string) string {
	return "ddns-go-linux-amd64-v3.tar.gz"
}
