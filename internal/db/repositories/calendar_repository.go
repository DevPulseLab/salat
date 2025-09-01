package repositories

import (
	"errors"
	"log"
	"time"

	"github.com/DevPulseLab/salat/internal/db/models"
	"github.com/DevPulseLab/salat/internal/dto"
	"github.com/DevPulseLab/salat/internal/enum"
	"github.com/DevPulseLab/salat/internal/helper"
	"github.com/uniplaces/carbon"
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

func (repo *CalendarRepository) ChangeEntryStatus(modelId uint, status string) (models.Calendar, error) {
	var calendarEntry models.Calendar
	result := repo.DB.Where("id = ? AND deleted_at IS NULL", modelId).First(&calendarEntry)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return calendarEntry, result.Error
	}

	calendarEntry.Status = status
	result = repo.DB.Save(&calendarEntry)

	return calendarEntry, result.Error
}

func (repo *CalendarRepository) CountReservedForDate(date *carbon.Carbon) int64 {
	var count int64
	result := repo.DB.Model(&models.Calendar{}).
		Where("status = ? AND DATE(date) = DATE(?) AND deleted_at IS NULL", string(enum.Reserved), date.Time).
		Count(&count)
	if result.Error != nil {
		log.Printf("Error: %v", result.Error)
		return 0
	}

	return count
}

func (repo *CalendarRepository) AddCalendarEntry(user *models.User, startDate, endDate time.Time, closeIntervals []dto.CloseInterval) ([]models.Calendar, []error) {
	currDate := time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, time.Local)
	now := carbon.Now().Time
	nowPlus30Days := carbon.Now().AddDate(0, 0, 30)

	errors := []error{}
	addedDays := []models.Calendar{}

	for endDate.Sub(currDate).Hours() >= 0 {
		if repo.dateHelper.IsWeekend(currDate) {
			currDate = currDate.AddDate(0, 0, 1)
			continue
		}

		if repo.dateHelper.IsDateInCloseIntervals(currDate, closeIntervals) {
			currDate = currDate.AddDate(0, 0, 1)
			continue
		}

		status := enum.Approved
		if user.PenaltyCard == string(enum.Yellow) {
			status = enum.Reserved
		} else if user.PenaltyCard == string(enum.Red) {
			status = enum.Rejected
		} else if currDate.Equal(now) || currDate.Before(now) || currDate.After(nowPlus30Days) {
			status = enum.Rejected
		} else if repo.dateHelper.IsDateInCurrentWeek(currDate) {
			status = enum.Reserved
		} else if repo.dateHelper.IsDateNextWeekAndNowAfterFriday(currDate) {
			status = enum.Reserved
		}

		var deletedCalendarEntry models.Calendar
		if err := repo.DB.Unscoped().Where("user_id = ? AND DATE(date) = DATE(?) AND deleted_at IS NOT NULL", user.ID, currDate).First(&deletedCalendarEntry).Error; err == nil {
			repo.DB.Unscoped().Model(&deletedCalendarEntry).Update("deleted_at", nil)
			deletedCalendarEntry.Status = string(status)
			saveErr := repo.DB.Save(deletedCalendarEntry).Error
			if saveErr != nil {
				errors = append(errors, saveErr)
			} else {
				addedDays = append(addedDays, deletedCalendarEntry)
			}
			continue
		}

		calendarModel := models.Calendar{UserId: user.ID, Date: currDate, Status: string(status)}
		insertErr := repo.DB.Create(&calendarModel).Error
		if insertErr != nil {
			errors = append(errors, insertErr)
		} else {
			addedDays = append(addedDays, calendarModel)
		}

		currDate = currDate.AddDate(0, 0, 1)
	}

	if len(errors) > 0 {
		return addedDays, errors
	}

	return addedDays, nil
}

func (repo *CalendarRepository) GetCalendarEntriesByUserId(userId uint, startDate, endDate time.Time) []models.Calendar {
	var calendars []models.Calendar

	repo.DB.Where("user_id = ? AND DATE(date) >= DATE(?) AND DATE(date) <= DATE(?)", userId, startDate, endDate).Find(&calendars)

	return calendars
}

func (repo *CalendarRepository) GetCalendarEntriesForAllUsers(startDate, endDate time.Time) []models.Calendar {
	var calendars []models.Calendar

	repo.DB.Where("DATE(date) >= DATE(?) AND DATE(date) <= DATE(?)", startDate, endDate).Find(&calendars)

	return calendars
}
