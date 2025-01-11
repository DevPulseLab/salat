package forms

type ChangeCalendarEntryForm struct {
	CalendarEntryId uint   `binding: "required" json:"calendarEntryId"`
	NewStatus       string `binding: "required" json:"newStatus"`
}
