/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/domaingts/mana/pkg/config"
	ddnsgo "github.com/domaingts/mana/pkg/static/ddns-go"
	"github.com/spf13/cobra"
)

func NewDdnsGoCommand() *cobra.Command {
	cmd := "ddns-go"
	cfg := config.NewConfig(cmd, "domaingts", cmd)
	cfg.SetBinaryPath("/usr/local/bin")
	cfg.SetConfigPath("/etc/ddns-go")
	command := &cobra.Command{
		Use: cmd,
		PreRun: func(cmd *cobra.Command, args []string) {
			cfg.InitConfig()
			if err := cfg.CreateService(ddnsgo.Service); err != nil {
				fmt.Println(err)
				os.Exit(0)
			}
			if err := cfg.CreateStartConfig(ddnsgo.Config); err != nil {
				fmt.Println(err)
				os.Exit(0)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {
			if err := cfg.Run(); err != nil {
				fmt.Println(err)
				os.Exit(0)
			}
		},
	}
	return command
}

func init() {
	rootCmd.AddCommand(NewDdnsGoCommand())
}
