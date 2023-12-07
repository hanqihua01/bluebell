package models

type ParamSignUp struct {
	Username   string `json:"username" binding:"required"` // binding标签是validator包读取的，gin集成了validator
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

type ParamLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ParamVoteData struct {
	//UserID 从当前登录用户里可以直接获取
	PostID    string `json:"post_id" binding:"required"`                       // 帖子id
	Direction int8   `json:"direction,string" binding:"required,oneof=1 0 -1"` // 投票类型，赞成(1)反对-(1)取消(0)
}
