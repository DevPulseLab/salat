package forms

import "time"

type AddCloseCalendarEntryForm struct {
	StartDate time.Time `binding: "required" json:"startDate"`
	EndDate   time.Time `binding: "required" json:"endDate"`
}
