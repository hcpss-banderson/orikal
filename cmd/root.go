/*
Copyright © 2025 Brendan Anderson brendan_anderson@hcpss.org
*/
package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var projectDir string
var configFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "orikal",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&projectDir, "projectDir", "", "Project directory (default is the current directory)")
	rootCmd.PersistentFlags().StringVar(&configFile, "configFile", "config/config.yml", "Kamal config file (relative to projectDir)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if projectDir == "" {
		dir, err := os.Getwd()
		cobra.CheckErr(err)
		rootCmd.PersistentFlags().Set("projectDir", dir)
	} else {
		absPath, err := filepath.Abs(projectDir)
		if err != nil {
			panic(err)
		}
		rootCmd.PersistentFlags().Set("projectDir", absPath)
	}
}
