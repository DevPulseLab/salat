package dto

import "time"

type Calendar struct {
	Id     uint      `json:"id"`
	Date   time.Time `json:"date"`
	Status string    `json:"status"`
}
