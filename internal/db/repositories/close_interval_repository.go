package repositories

import (
	"time"

	"github.com/DevPulseLab/salat/internal/db/models"
	"gorm.io/gorm"
)

type CloseIntervalRepository struct {
	DB *gorm.DB
}

func NewCloseIntervalsRepository(db *gorm.DB) *CloseIntervalRepository {
	return &CloseIntervalRepository{DB: db}
}

func (repo *CloseIntervalRepository) SaveCloseInterval(startDate time.Time, endDate time.Time) error {
	closeIntervalEntry := models.CloseInterval{StartDate: startDate, EndDate: endDate}

	return repo.DB.Save(&closeIntervalEntry).Error
}

func (repo *CloseIntervalRepository) GetAllEntriesForInterval(startDate time.Time, endDate time.Time) []models.CloseInterval {
	var closeDateIntervals []models.CloseInterval

	repo.DB.Where("(start_date >= ? AND start_date <= ?) OR (end_date >= ? AND end_date <= ?)", startDate, endDate, startDate, endDate).Find(&closeDateIntervals)

	return closeDateIntervals
}
