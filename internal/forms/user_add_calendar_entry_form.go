package forms

import "time"

type UserAddCalendarEntryForm struct {
	StartDate time.Time `binding: "required" json:"startDate"`
	EndDate   time.Time `binding: "required" json:"endDate"`
}
