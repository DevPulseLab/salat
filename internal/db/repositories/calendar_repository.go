package repositories

import (
	"errors"
	"time"

	"github.com/DevPulseLab/salat/internal/db/models"
	"github.com/DevPulseLab/salat/internal/enum"
	"github.com/uniplaces/carbon"
	"gorm.io/gorm"
)

type CalendarRepository struct {
	DB *gorm.DB
}

func NewCalendarRepository(db *gorm.DB) *CalendarRepository {
	return &CalendarRepository{DB: db}
}

func (repo *CalendarRepository) FindByIdForUserId(id, userId uint) (models.Calendar, error) {
	var calendarEntry models.Calendar
	result := repo.DB.Where("id = ? AND user_id = ? AND deleted_at IS NULL", id, userId).First(&calendarEntry)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return calendarEntry, result.Error
	}

	return calendarEntry, nil
}

func (repo *CalendarRepository) SoftDelete(model *models.Calendar) {
	repo.DB.Delete(&model)
}

func (repo *CalendarRepository) UpdateStatus(modelId uint, status string) (models.Calendar, error) {
	var calendarEntry models.Calendar
	result := repo.DB.Where("id = ? AND deleted_at IS NULL", modelId).First(&calendarEntry)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return calendarEntry, result.Error
	}

	calendarEntry.Status = status
	result = repo.DB.Save(&calendarEntry)

	return calendarEntry, result.Error
}

func (repo *CalendarRepository) CountReservedByDate(date *carbon.Carbon) (int64, error) {
	var count int64
	result := repo.DB.Model(&models.Calendar{}).
		Where("status = ? AND DATE(date) = ? AND deleted_at IS NULL", string(enum.Reserved), date.Time.Format("2006-01-02")).
		Count(&count)
	return count, result.Error
}

func (repo *CalendarRepository) FindDeletedByUserIdAndDate(userID uint, date time.Time) (models.Calendar, error) {
	var calendarEntry models.Calendar
	result := repo.DB.Unscoped().Where("user_id = ? AND date = ? AND deleted_at IS NOT NULL", userID, date).First(&calendarEntry)
	return calendarEntry, result.Error
}

func (repo *CalendarRepository) RestoreAndUpdate(calendarEntry *models.Calendar, status string) error {
	return repo.DB.Unscoped().Model(&calendarEntry).Updates(map[string]interface{}{
		"deleted_at": nil,
		"status":     status,
	}).Error
}

func (repo *CalendarRepository) Create(calendarEntry *models.Calendar) error {
	return repo.DB.Create(calendarEntry).Error
}

func (repo *CalendarRepository) FindByUserIdAndDateRange(userID uint, startDate, endDate time.Time) ([]models.Calendar, error) {
	var calendars []models.Calendar
	result := repo.DB.Where("user_id = ? AND DATE(date) >= ? AND DATE(date) <= ?", userID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).Find(&calendars)
	return calendars, result.Error
}

func (repo *CalendarRepository) FindByDateRange(startDate, endDate time.Time) ([]models.Calendar, error) {
	var calendars []models.Calendar
	result := repo.DB.Where("DATE(date) >= ? AND DATE(date) <= ?", startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).Find(&calendars)
	return calendars, result.Error
}
