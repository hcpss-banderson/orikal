package model

import (
	"encoding/json"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"strconv"
	"time"
)

type StringInt int

type MigrationImportStatus struct {
	Id           string    `json:"id"`
	GroupId      string    `json:"group"`
	Status       string    `json:"status"`
	Total        StringInt `json:"total"`
	Unprocessed  StringInt `json:"unprocessed"`
	MessageCount StringInt `json:"message_count"`
	LastImported string    `json:"last_imported"`
	Acronym      string
}

func (m *MigrationImportStatus) ToRow() table.Row {
	red := text.Colors{text.FgRed}
	yellow := text.Colors{text.FgYellow}

	var status string
	switch m.Status {
	case "Importing":
		status = yellow.Sprint(m.Status)
	case "Idle":
		status = m.Status
	default:
		status = red.Sprint(m.Status)
	}

	var unprocessed string
	if m.Unprocessed != 0 {
		unprocessed = yellow.Sprint(m.Unprocessed)
	} else {
		unprocessed = strconv.Itoa(int(m.Unprocessed))
	}

	var messagCount string
	if m.MessageCount != 0 {
		messagCount = red.Sprint(m.MessageCount)
	} else {
		messagCount = strconv.Itoa(int(m.MessageCount))
	}

	var lastImported string
	if m.LastImported == "" {
		lastImported = red.Sprint("NOT RUN")
	} else {
		dateRun, err := time.Parse("2006-01-02 15:04:05", m.LastImported)
		if err != nil {
			panic(err)
		}
		threeDaysAgo := time.Now().Add(-24 * 3 * time.Hour)
		if dateRun.Before(threeDaysAgo) {
			lastImported = yellow.Sprint(m.LastImported)
		} else {
			lastImported = m.LastImported
		}
	}

	return table.Row{
		m.Acronym,
		m.GroupId,
		m.Id,
		status,
		strconv.Itoa(int(m.Total)),
		unprocessed,
		messagCount,
		lastImported,
	}
}

func (st *StringInt) UnmarshalJSON(b []byte) error {
	//convert the bytes into an interface
	//this will help us check the type of our value
	//if it is a string that can be converted into an int we convert it
	///otherwise we return an error
	var item interface{}
	if err := json.Unmarshal(b, &item); err != nil {
		return err
	}
	switch v := item.(type) {
	case int:
		*st = StringInt(v)
	case float64:
		*st = StringInt(int(v))
	case string:
		///here convert the string into
		///an integer
		if v == "" {
			*st = 0
		} else {
			i, err := strconv.Atoi(v)
			if err != nil {
				///the string might not be of integer type
				///so return an error
				return err

			}
			*st = StringInt(i)
		}
	}
	return nil
}
