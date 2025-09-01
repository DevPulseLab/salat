package repositories

import (
	"time"

	"github.com/DevPulseLab/salat/internal/db/models"
	"gorm.io/gorm"
)

type VisitStatsRepository struct {
	DB *gorm.DB
}

func NewVisitStatsRepository(db *gorm.DB) *VisitStatsRepository {
	return &VisitStatsRepository{DB: db}
}

func (repo *VisitStatsRepository) ToggleVisit(userId uint, visitDate time.Time) (*models.VisitStats, error) {
	var statsEntry models.VisitStats
	if err := repo.DB.Where("user_id = ? AND DATE(date) = DATE(?)", userId, visitDate).First(&statsEntry).Error; err == nil {
		statsEntry.IsVisit = !statsEntry.IsVisit
	} else {
		statsEntry = models.VisitStats{UserId: userId, Date: visitDate, IsVisit: true}
	}

	err := repo.DB.Save(&statsEntry).Error
	return &statsEntry, err
}

func (repo *VisitStatsRepository) GetVisitVisit(startDate, endDate time.Time) []models.VisitStats {
	var visitStatsList []models.VisitStats

	repo.DB.Where("DATE(date) >= (?) AND DATE(date) <= DATE(?)", startDate, endDate).Find(&visitStatsList)

	return visitStatsList
}
