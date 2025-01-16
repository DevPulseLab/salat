package forms

import "time"

type SaveNumberOfPlatesForDayForm struct {
	StatsDay       time.Time `binding: "required" json:"statsDay"`
	NumberOfPlates int       `binding: "required" json:"numberOfPlates"`
}
