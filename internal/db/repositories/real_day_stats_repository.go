package repositories

import (
	"time"

	"github.com/DevPulseLab/salat/internal/db/models"
	"gorm.io/gorm"
)

type NewRealDayStatsRepository struct {
	DB *gorm.DB
}

func NewNewRealDayStatsRepository(db *gorm.DB) *NewRealDayStatsRepository {
	return &NewRealDayStatsRepository{DB: db}
}

func (repo *NewRealDayStatsRepository) IncrementStatsForDay(statsDay time.Time) bool {
	var statsEntry models.RealDayStats
	if err := repo.DB.Where("date = ?", statsDay).First(&statsEntry).Error; err == nil {
		statsEntry.NumberOfPlates++
	} else {
		statsEntry = models.RealDayStats{Date: statsDay, NumberOfPlates: 1}
	}

	err := repo.DB.Save(&statsEntry).Error
	return err == nil
}

func (repo *NewRealDayStatsRepository) SaveStatsForDay(statsDay time.Time, numberOfPlates int) bool {
	var statsEntry models.RealDayStats
	if err := repo.DB.Where("date = ?", statsDay).First(&statsEntry).Error; err == nil {
		statsEntry.NumberOfPlates = uint(numberOfPlates)
	} else {
		statsEntry = models.RealDayStats{Date: statsDay, NumberOfPlates: uint(numberOfPlates)}
	}

	err := repo.DB.Save(&statsEntry).Error
	return err == nil
}

func (repo *NewRealDayStatsRepository) GetStatsForDay(statsDay time.Time) uint {
	var statsEntry models.RealDayStats
	if err := repo.DB.Where("date = ?", statsDay).First(&statsEntry).Error; err == nil {
		return statsEntry.NumberOfPlates
	} else {
		return 0
	}
}
