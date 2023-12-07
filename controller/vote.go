package controller

import (
	"bluebell/logic"
	"bluebell/models"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func PostVoteHandler(c *gin.Context) {
	// 参数校验
	p := new(models.ParamVoteData)
	if err := c.ShouldBindJSON(p); err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	userID, err := GetCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
		return
	}
	// 业务处理
	if err := logic.VoteForPost(userID, p); err != nil {
		zap.L().Error("voting for post failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回
	ResponseSuccess(c, nil)
}
