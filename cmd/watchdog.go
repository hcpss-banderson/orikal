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

var watchdogCmd = &cobra.Command{
	Use:   "watchdog",
	Short: "Show watchdog entries",
	Long:  `Show watchdog entries`,
	Run: func(cmd *cobra.Command, args []string) {
		dir := cmd.Flag("projectDir").Value.String()
		config := cmd.Flag("configFile").Value.String()
		severity := cmd.Flag("severity-min").Value.String()
		count := cmd.Flag("count").Value.String()
		kamal := service.NewKamalService(dir, config)

		roles := kamal.GetRoles()
		bar := progressbar.NewOptions(len(roles),
			progressbar.OptionSetWidth(32),
		)

		acronymChan := make(chan string)
		dataChan := make(chan model.Payload)
		go func() {
			kamal.AppExec("drush watchdog:show --format=json --severity-min="+severity+" --count="+count, acronymChan, dataChan)
		}()

		for acronym := range acronymChan {
			bar.Describe("Received " + acronym + "...")
			bar.Add(1)
		}
		fmt.Println()

		var report []model.WatchdogEntry
		for value := range dataChan {
			acronym := value.Acronym
			data := value.Data
			//var dat []model.WatchdogEntry
			var gen map[string]model.WatchdogEntry

			if err := json.Unmarshal([]byte(data), &gen); err != nil {
				panic(err)
			}
			//fmt.Println("START")
			//spew.Dump(gen)
			//fmt.Println("STOP")
			//
			//if err := json.Unmarshal([]byte(data), &dat); err != nil {
			//	panic(err)
			//}

			for _, d := range gen {
				d.Acronym = acronym
				report = append(report, d)
			}
		}

		t := table.NewWriter()
		t.AppendHeader(table.Row{"Acronym", "ID", "Date", "Type", "Severity", "Message"})
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
	rootCmd.AddCommand(watchdogCmd)
	watchdogCmd.Flags().Int("severity-min", 4, "Set the minimum event severity level: Emergency(0), Alert(1), Critical(2), Error(3), Warning(4), Notice(5), Info(6), Debug(7)")
	watchdogCmd.Flags().Int("count", 25, "Number of records to show.")
}
