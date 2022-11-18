package simka

type UserLoginInput struct {
	NIM          string `json:"nim" binding:"required"`
	Password     string `json:"password" binding:"required"`
	UserAnswer   string `json:"user_answer"`
	SystemAnswer string `json:"system_answer"`
}
