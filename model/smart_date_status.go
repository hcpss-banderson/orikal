package model

import (
	"github.com/jedib0t/go-pretty/v6/table"
)

type SmartDateStatus struct {
	Label   string `json:"label"`
	Count   int    `json:"count"`
	Acronym string
}

func (s *SmartDateStatus) ToRow() table.Row {
	return table.Row{
		s.Acronym,
		s.Label,
		s.Count,
	}
}
