package enum

type CalendarStatus string

const (
	Approved CalendarStatus = "approved"
	Rejected CalendarStatus = "rejected"
	Reserved CalendarStatus = "reserved"
)
