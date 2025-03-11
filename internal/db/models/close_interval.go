package models

import (
	"time"

	"gorm.io/gorm"
)

type CloseInterval struct {
	gorm.Model
	StartDate time.Time `gorm:"not null" json:"start_date"`
	EndDate   time.Time `gorm:"not null" json:"end_date"`
}
