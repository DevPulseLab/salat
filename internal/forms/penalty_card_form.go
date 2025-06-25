package forms

type PenaltyCardForm struct {
	UserId   uint   `binding: "required" json:"userId"`
	CardType string `binding: "required" json:"cardType"`
}
