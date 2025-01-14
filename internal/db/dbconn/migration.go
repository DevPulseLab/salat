package dbconn

import (
	"github.com/DevPulseLab/salat/internal/db/models"
	"gorm.io/gorm"
)

func RunMigrate(db *gorm.DB) {
	err := db.AutoMigrate(
		&models.Calendar{},
		&models.User{},
		&models.RealDayStats{},
	)
	if err != nil {
		panic("migration failure")
	}
}
