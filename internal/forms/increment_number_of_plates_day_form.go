package forms

import "time"

type IncrementNumberOfPlatesForm struct {
	StatsDay time.Time `binding: "required" json:"statsDay"`
}
