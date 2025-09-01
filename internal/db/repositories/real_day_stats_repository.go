package repositories

import (
	"time"

	"github.com/DevPulseLab/salat/internal/db/models"
	"gorm.io/gorm"
)

type RealDayStatsRepository struct {
	DB *gorm.DB
}

func NewRealDayStatsRepository(db *gorm.DB) *RealDayStatsRepository {
	return &RealDayStatsRepository{DB: db}
}

func (repo *RealDayStatsRepository) IncrementStatsForDay(statsDay time.Time) bool {
	var statsEntry models.RealDayStats
	if err := repo.DB.Where("DATE(date) = ?", statsDay.Format("2006-01-02")).First(&statsEntry).Error; err == nil {
		statsEntry.NumberOfPlates++
	} else {
		statsEntry = models.RealDayStats{Date: statsDay, NumberOfPlates: 1}
	}

	err := repo.DB.Save(&statsEntry).Error
	return err == nil
}

func (repo *RealDayStatsRepository) SaveStatsForDay(statsDay time.Time, numberOfPlates int) bool {
	var statsEntry models.RealDayStats
	if err := repo.DB.Where("DATE(date) = ?", statsDay.Format("2006-01-02")).First(&statsEntry).Error; err == nil {
		statsEntry.NumberOfPlates = uint(numberOfPlates)
	} else {
		statsEntry = models.RealDayStats{Date: statsDay, NumberOfPlates: uint(numberOfPlates)}
	}

	err := repo.DB.Save(&statsEntry).Error
	return err == nil
}

func (repo *RealDayStatsRepository) GetStatsForDay(statsDay time.Time) uint {
	var statsEntry models.RealDayStats
	if err := repo.DB.Where("DATE(date) = ?", statsDay.Format("2006-01-02")).First(&statsEntry).Error; err == nil {
		return statsEntry.NumberOfPlates
	} else {
		return 0
	}
}
