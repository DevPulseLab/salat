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

	repo.DB.
		Where(
			`(start_date >= @startDate AND start_date <= @endDate) 
			OR (end_date >= @startDate AND end_date <= @endDate) 
			OR (start_date <= @startDate AND end_date >= @startDate) 
			OR (start_date <= @endDate AND end_date >= @endDate)`,
			map[string]interface{}{
				"startDate": startDate,
				"endDate":   endDate,
			},
		).
		Find(&closeDateIntervals)

	return closeDateIntervals
}
