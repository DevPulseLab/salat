package forms

type LoginForm struct {
	Username string `form:"username" binding: "required" json:"username"`
	Password string `form:"password" binding: "required" json:"password"`
}
