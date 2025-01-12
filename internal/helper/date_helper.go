package helper

import "time"

type DateHelper struct{}

func NewDateHelper() *DateHelper {
	return &DateHelper{}
}

func (helper *DateHelper) IsDateInCurrentWeek(t time.Time) bool {
	year, week := time.Now().ISOWeek()
	targetYear, targetWeek := t.ISOWeek()
	return year == targetYear && week == targetWeek
}

func (helper *DateHelper) IsDateNextWeekAndNowAfterFriday(t time.Time) bool {
	fridayThisWeek := helper.getFridayOfWeek(time.Now())

	year, week := time.Now().AddDate(0, 0, 7).ISOWeek()
	targetYear, targetWeek := t.ISOWeek()
	return year == targetYear && week == targetWeek && time.Now().After(fridayThisWeek)
}

func (helper *DateHelper) getFridayOfWeek(inputDate time.Time) time.Time {
	mondayOffset := (int(inputDate.Weekday()) + 6) % 7
	monday := inputDate.AddDate(0, 0, -mondayOffset) // Get the Monday of the current week
	return monday.AddDate(0, 0, 4)
}

func (helper *DateHelper) IsWeekend(inputDate time.Time) bool {
	weekday := inputDate.Weekday()
	return weekday == time.Saturday || weekday == time.Sunday
}
