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

func (helper *DateHelper) IsDateNextWeekAndNowBeforeFriday(t time.Time) bool {
	fridayThisWeek := helper.getFridayOfWeek(time.Now())

	year, week := time.Now().AddDate(0, 0, 7).ISOWeek()
	targetYear, targetWeek := t.ISOWeek()
	return year == targetYear && week == targetWeek && time.Now().After(fridayThisWeek)
}

func (helper *DateHelper) getFridayOfWeek(inputDate time.Time) time.Time {
	weekday := inputDate.Weekday()
	daysUntilFriday := (time.Friday - weekday + 7) % 7
	friday := inputDate.AddDate(0, 0, int(daysUntilFriday))
	return time.Date(friday.Year(), friday.Month(), friday.Day(), 12, 0, 0, 0, friday.Location())
}
