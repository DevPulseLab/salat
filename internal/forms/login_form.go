package forms

type LoginForm struct {
	Username string `binding: "required" json:"username"`
	Password string `binding: "required" json:"password"`
}
