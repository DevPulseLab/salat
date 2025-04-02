package dto

import "time"

type VisitStats struct {
	Id      uint      `json:"id"`
	UserId  uint      `json:"userId"`
	Date    time.Time `json:"date"`
	IsVisit bool      `json:"isVisit"`
}
