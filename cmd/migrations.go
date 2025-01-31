/*
Copyright © 2025 Brendan Anderson <brendan_anderson@hcpss.org>
*/
package cmd

import (
	"fmt"
	"github.com/hcpss-banderson/orikal/service"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

// migrationsCmd represents the migrations command
var migrationsCmd = &cobra.Command{
	Use:   "migrations",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		dir := cmd.Flag("projectDir").Value.String()
		config := cmd.Flag("configFile").Value.String()
		kamal := service.NewKamalService(dir, config)
		report := kamal.AppExec("drush ms --format=json")

		t := table.NewWriter()
		t.AppendHeader(table.Row{"Acronym", "Group", "Id", "Status", "Total", "Unprocessed", "MessageCount", "LastImported"})
		prev := ""
		for _, r := range report {
			if prev != r.Acronym {
				t.AppendSeparator()
				prev = r.Acronym
			}

			t.AppendRow(r.ToRow())

		}
		t.SetColumnConfigs([]table.ColumnConfig{
			{Number: 1, AutoMerge: true},
		})
		fmt.Println(t.Render())
	},
}

func init() {
	rootCmd.AddCommand(migrationsCmd)
}
