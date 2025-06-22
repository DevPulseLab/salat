package helper

import (
	"time"

	"github.com/DevPulseLab/salat/internal/dto"
	"github.com/uniplaces/carbon"
)

type DateHelper struct{}

func NewDateHelper() *DateHelper {
	return &DateHelper{}
}

func (helper *DateHelper) IsDateInCurrentWeek(t time.Time) bool {
	year, week := carbon.Now().ISOWeek()
	targetYear, targetWeek := t.ISOWeek()
	return year == targetYear && week == targetWeek
}

func (helper *DateHelper) IsDateNextWeekAndNowAfterFriday(t time.Time) bool {
	fridayThisWeek := helper.getFridayOfWeek(carbon.Now().Time)
	fridayThisWeekLunch := time.Date(fridayThisWeek.Year(), fridayThisWeek.Month(), fridayThisWeek.Day(), 12, 0, 0, 0, time.Local)

	year, week := carbon.Now().AddDate(0, 0, 7).ISOWeek()
	targetYear, targetWeek := t.ISOWeek()

	return year == targetYear && week == targetWeek && carbon.Now().After(fridayThisWeekLunch)
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

func (helper *DateHelper) IsDateInCloseIntervals(inputDate time.Time, closeIntervals []dto.CloseInterval) bool {
	for _, closeInterval := range closeIntervals {
		if inputDate.Unix() >= closeInterval.StartDate.Unix() && inputDate.Unix() <= closeInterval.EndDate.Unix() {
			return true
		}
	}

	return false
}

func ParseDate(dateString string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02", dateString, time.Local)
}
