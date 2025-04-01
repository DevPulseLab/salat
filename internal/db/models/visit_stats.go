package models

import (
	"time"

	"gorm.io/gorm"
)

type VisitStats struct {
	gorm.Model
	UserId  uint      `gorm:"not null;uniqueIndex:unq_usr_visit" json:"user_id"`
	Date    time.Time `gorm:"not null;uniqueIndex:unq_usr_visit" json:"date"`
	IsVisit bool      `gorm:"type:bool;not null" json:"is_visit"`
}
