package models

import (
	"time"

	"gorm.io/gorm"
)

type Calendar struct {
	gorm.Model
	UserId uint      `gorm:"not null" json:"user_id"`
	Date   time.Time `gorm:"not null" json:"date"`
	Status string    `gorm:"type:varchar(100);not null" json:"status"`
}
