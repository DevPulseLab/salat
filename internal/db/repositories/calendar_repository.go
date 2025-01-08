package repositories

import (
	"time"

	"github.com/DevPulseLab/salat/internal/db/models"
	"github.com/DevPulseLab/salat/internal/enum"
	"gorm.io/gorm"
)

type CalendarRepository struct {
	DB *gorm.DB
}

func NewCalendarRepository(db *gorm.DB) *CalendarRepository {
	return &CalendarRepository{DB: db}
}

func (repo *CalendarRepository) AddCalendarEntry(userId uint, startDate, endDate time.Time) (bool, []error) {
	currDate := startDate

	errors := []error{}

	for endDate.Sub(currDate).Hours() > 0 {
		status := enum.Approved
		if currDate.Before(time.Now()) {
			status = enum.Rejected
		} else if isDateInCurrentWeek(currDate) {
			status = enum.Reserved
		} else if isDateNextWeekAndNowBeforeFriday(currDate) {
			status = enum.Reserved
		}

		insertErr := repo.DB.Create(&models.Calendar{UserId: userId, Date: currDate, Status: string(status)}).Error
		if insertErr != nil {
			errors = append(errors, insertErr)
		}

		currDate = currDate.AddDate(0, 0, 1)
	}

	if len(errors) > 0 {
		return false, errors
	}

	return true, nil
}

func (repo *CalendarRepository) GetCalendarEntriesByUserId(userId uint, startDate, endDate time.Time) []models.Calendar {
	var calendars []models.Calendar

	repo.DB.Where("user_id = ? AND date >= ? AND date <= ?", userId, startDate, endDate).Find(&calendars)

	return calendars
}

func (repo *CalendarRepository) GetCalendarEntriesForAllUsers(startDate, endDate time.Time) []models.Calendar {
	var calendars []models.Calendar

	repo.DB.Where("date >= ? AND date <= ?", startDate, endDate).Find(&calendars)

	return calendars
}

func isDateInCurrentWeek(t time.Time) bool {
	year, week := time.Now().ISOWeek()
	targetYear, targetWeek := t.ISOWeek()
	return year == targetYear && week == targetWeek
}

func isDateNextWeekAndNowBeforeFriday(t time.Time) bool {
	fridayThisWeek := getFridayOfWeek(time.Now())

	year, week := time.Now().AddDate(0, 0, 7).ISOWeek()
	targetYear, targetWeek := t.ISOWeek()
	return year == targetYear && week == targetWeek && time.Now().After(fridayThisWeek)
}

func getFridayOfWeek(inputDate time.Time) time.Time {
	weekday := inputDate.Weekday()
	daysUntilFriday := (time.Friday - weekday + 7) % 7
	friday := inputDate.AddDate(0, 0, int(daysUntilFriday))
	return time.Date(friday.Year(), friday.Month(), friday.Day(), 12, 0, 0, 0, friday.Location())
}
