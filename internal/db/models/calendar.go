package models

import (
	"time"

	"github.com/DevPulseLab/salat/internal/enum"
	"gorm.io/gorm"
)

type Calendar struct {
	gorm.Model
	UserId uint      `gorm:"not null;uniqueIndex:unq_usr_date" json:"user_id"`
	Date   time.Time `gorm:"not null;uniqueIndex:unq_usr_date" json:"date"`
	Status string    `gorm:"type:varchar(100);not null" json:"status"`
}

func (model *Calendar) IsApproved() bool {
	return model.Status == string(enum.Approved)
}

func (model *Calendar) IsRejected() bool {
	return model.Status == string(enum.Rejected)
}
