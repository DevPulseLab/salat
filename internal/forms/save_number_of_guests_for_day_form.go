package forms

import "time"

type SaveNumberOfGuestsForDayForm struct {
	StatsDay       time.Time `binding: "required" json:"statsDay"`
	NumberOfGuests int       `binding: "required" json:"numberOfGuests"`
}
