/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "globus",
	Short: "CLI app for the globus transfer library of the OpenEM project",
	Long: `
This CLI app demonstrates all aspects of the Globus library that it 
is distributed with. It can be used on its own to interact with the Globus
API, and is specifically geared towards transfer requests and associated
task monitoring and management.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {}
