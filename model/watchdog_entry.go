package model

import "github.com/jedib0t/go-pretty/v6/table"

type WatchdogEntry struct {
	Acronym  string
	Id       string `json:"wid"`
	Date     string `json:"date"`
	Type     string `json:"type"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
}

func (w *WatchdogEntry) ToRow() table.Row {
	return table.Row{
		w.Acronym,
		w.Id,
		w.Date,
		w.Type,
		w.Severity,
		w.Message,
	}
}
