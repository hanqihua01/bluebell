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
	PostID    string `json:"post_id" binding:"required"`              // 帖子id
	Direction int8   `json:"direction,string" binding:"oneof=1 0 -1"` // 投票类型，赞成(1)反对-(1)取消(0)
}

type ParamPostList struct {
	Page        int64  `json:"page" form:"page"` // form: url ?key=val 参数的解析
	Size        int64  `json:"size" form:"size"`
	Order       string `json:"order" form:"order"`
	CommunityID int64  `json:"community_id" form:"community_id"` // 可以为空
}

const (
	OrderTime  = "time"
	OrderScore = "score"
)
