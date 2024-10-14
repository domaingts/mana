/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/domaingts/mana/pkg/config"
	"github.com/spf13/cobra"
)

func NewDdnsGoCommand() *cobra.Command {
	cmd := "ddns-go"
	cfg := config.NewConfig(cmd, "domaingts", cmd)
	command := &cobra.Command{
		Use: cmd,
		PreRun: func(cmd *cobra.Command, args []string) {
			cfg.SetBinaryPath("/usr/local/bin")
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
