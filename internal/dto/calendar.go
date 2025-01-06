package dto

import "time"

type Calendar struct {
	Id     uint      `json:"id"`
	UserId uint      `json:"userId"`
	Date   time.Time `json:"date"`
	Status string    `json:"status"`
}
