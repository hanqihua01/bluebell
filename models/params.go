package models

type ParamSignUp struct {
	Username   string `json:"username" binding:"required"` // binding标签是validator包读取的，gin集成了validator
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}
