package dto

import "time"

type CloseInterval struct {
	Id        uint      `json:"id"`
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
}
