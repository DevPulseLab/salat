package models

import (
	"time"

	"gorm.io/gorm"
)

type RealDayStats struct {
	gorm.Model
	Date           time.Time `gorm:"not null;uniqueIndex:unq_real_stats_date" json:"date"`
	NumberOfPlates uint      `gorm:"not null;default:0" json:"number_of_plates"`
}
