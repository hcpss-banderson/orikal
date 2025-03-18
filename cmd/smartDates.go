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

// smartDatesCmd represents the smartDates command
var smartDatesCmd = &cobra.Command{
	Use:   "smartDates",
	Short: "An overview of smart date usage",
	Long:  `An overview of smart date usage`,
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
			kamal.AppExec("drush hcpss_event_importer:event-report --format=json", acronymChan, dataChan)
		}()

		for acronym := range acronymChan {
			bar.Describe("Received " + acronym + "...")
			bar.Add(1)
		}
		fmt.Println()

		var report []model.SmartDateStatus
		for value := range dataChan {
			acronym := value.Acronym
			data := value.Data
			var dat []model.SmartDateStatus
			if err := json.Unmarshal([]byte(data), &dat); err != nil {
				panic(err)
			}

			for _, d := range dat {
				d.Acronym = acronym
				report = append(report, d)
			}
		}

		t := table.NewWriter()
		t.AppendHeader(table.Row{"Acronym", "Labels", "Count"})
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
	rootCmd.AddCommand(smartDatesCmd)
}
