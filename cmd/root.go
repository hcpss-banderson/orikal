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
	Short: "See what is happening on the school sites",
	Long: `Orikal sees all with his Infinite Eye. A fully functional school site 
deployment directory is required. You can either run orikal from that
directory, or by passing the --projectDir flag.
`,
}

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
