package dto

type User struct {
	Id          uint   `json:"id"`
	Username    string `json:"username"`
	Role        string `json:"role"`
	PenaltyCard string `json:"penaltyCard"`
}
