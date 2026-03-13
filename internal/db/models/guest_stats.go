package models

import (
	"time"

	"gorm.io/gorm"
)

type GuestStats struct {
	gorm.Model
	Date           time.Time `gorm:"not null;uniqueIndex:unq_guest_stats_date" json:"date"`
	NumberOfGuests uint      `gorm:"not null;default:0" json:"number_of_guests"`
}
