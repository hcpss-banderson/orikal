/*
Copyright © 2025 Brendan Anderson <brendan_anderson@hcpss.org>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/hcpss-banderson/orikal/model"
	"github.com/hcpss-banderson/orikal/service"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/schollz/progressbar/v3"
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

		roles := kamal.GetRoles()
		bar := progressbar.NewOptions(len(roles),
			progressbar.OptionSetWidth(32),
		)

		acronymChan := make(chan string)
		dataChan := make(chan model.Payload)
		go func() {
			kamal.AppExec("drush ms --format=json", acronymChan, dataChan)
		}()

		for acronym := range acronymChan {
			bar.Add(1)
			bar.Describe("Receiving " + acronym + "...")
		}
		fmt.Println()

		var report []model.MigrationImportStatus
		for value := range dataChan {
			acronym := value.Acronym
			data := value.Data
			var dat []model.MigrationImportStatus
			if err := json.Unmarshal([]byte(data), &dat); err != nil {
				panic(err)
			}

			for _, d := range dat {
				if d.Id != "" {
					d.Acronym = acronym
					report = append(report, d)
				}
			}
		}

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
