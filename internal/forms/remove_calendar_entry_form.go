package forms

type RemoveCalendarEntryForm struct {
	CalendarEntryId uint `binding: "required" json:"calendarEntryId"`
}
