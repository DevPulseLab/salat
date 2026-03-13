package repositories

import (
	"time"

	"github.com/DevPulseLab/salat/internal/db/models"
	"gorm.io/gorm"
)

type GuestStatsRepository struct {
	DB *gorm.DB
}

func NewGuestStatsRepository(db *gorm.DB) *GuestStatsRepository {
	return &GuestStatsRepository{DB: db}
}

func (repo *GuestStatsRepository) SaveStatsForDay(statsDay time.Time, numberOfGuests int) bool {
	var statsEntry models.GuestStats
	if err := repo.DB.Where("DATE(date) = ?", statsDay.Format("2006-01-02")).First(&statsEntry).Error; err == nil {
		statsEntry.NumberOfGuests = uint(numberOfGuests)
	} else {
		statsEntry = models.GuestStats{Date: statsDay, NumberOfGuests: uint(numberOfGuests)}
	}

	err := repo.DB.Save(&statsEntry).Error
	return err == nil
}

func (repo *GuestStatsRepository) GetStatsForDay(statsDay time.Time) uint {
	var statsEntry models.GuestStats
	if err := repo.DB.Where("DATE(date) = ?", statsDay.Format("2006-01-02")).First(&statsEntry).Error; err == nil {
		return statsEntry.NumberOfGuests
	} else {
		return 0
	}
}
