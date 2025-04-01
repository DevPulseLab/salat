package forms

import "time"

type ToggleVisitForm struct {
	UserId    uint      `binding: "required" json:"userId"`
	VisitDate time.Time `binding: "required" json:"visitDate"`
}
