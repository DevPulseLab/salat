package repositories

import (
	"errors"
	"time"

	"github.com/DevPulseLab/salat/internal/db/models"
	"github.com/DevPulseLab/salat/internal/enum"
	"github.com/DevPulseLab/salat/internal/helper"
	"gorm.io/gorm"
)

type CalendarRepository struct {
	DB         *gorm.DB
	dateHelper *helper.DateHelper
}

func NewCalendarRepository(db *gorm.DB, dh *helper.DateHelper) *CalendarRepository {
	return &CalendarRepository{DB: db, dateHelper: dh}
}

func (repo *CalendarRepository) GetByIdForUserId(id, userId uint) (models.Calendar, error) {
	var calendarEntry models.Calendar
	result := repo.DB.Where("id = ? AND user_id = ? AND deleted_at IS NULL", id, userId).First(&calendarEntry)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return calendarEntry, result.Error
	}

	return calendarEntry, nil
}

func (repo *CalendarRepository) Remove(model *models.Calendar) {
	repo.DB.Delete(&model)
}

func (repo *CalendarRepository) ChangeEntryStatus(modelId uint, status string) error {
	var calendarEntry models.Calendar
	result := repo.DB.Where("id = ? AND deleted_at IS NULL", modelId).First(&calendarEntry)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error
	}

	calendarEntry.Status = status
	result = repo.DB.Save(&calendarEntry)

	return result.Error
}

func (repo *CalendarRepository) AddCalendarEntry(userId uint, startDate, endDate time.Time) (bool, []error) {
	currDate := startDate

	errors := []error{}

	for endDate.Sub(currDate).Hours() > 0 {
		status := enum.Approved
		if currDate.Before(time.Now()) {
			status = enum.Rejected
		} else if repo.dateHelper.IsDateInCurrentWeek(currDate) {
			status = enum.Reserved
		} else if repo.dateHelper.IsDateNextWeekAndNowBeforeFriday(currDate) {
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
