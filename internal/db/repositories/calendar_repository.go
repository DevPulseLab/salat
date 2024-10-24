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

func (repo *CalendarRepository) AddCalendarEntry(userId uint, startDate, endDate time.Time, status enum.CalendarStatus) (bool, []error) {
	currDate := startDate

	errors := []error{}

	for endDate.Sub(currDate).Hours() > 0 {
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

func (repo *CalendarRepository) GetCalendarEntriesByUserId(userId uint) []models.Calendar {
	var calendars []models.Calendar

	repo.DB.Find(&calendars).Where("user_id = ?", userId)

	return calendars
}
