package forms

import "time"

type GetNumberOfPlatesForDayForm struct {
	StatsDay time.Time `binding: "required" json:"statsDay"`
}
