package models

import "gorm.io/gorm"

const (
	RoleGuest = "guest"
	RoleUser  = "user"
	RoleAdmin = "admin"
)

type User struct {
	gorm.Model
	Username    string `gorm:"size:255;not null;unique" json:"username"` // Email
	Password    string `gorm:"size:255;not null" json:"password"`
	Role        string `gorm:"default:guest;not null"`
	PenaltyCard string `gorm:"null"`

	Data []Calendar `gorm:"foreignKey:UserId" json:"calendars"`
}
