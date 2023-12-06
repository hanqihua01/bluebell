package controller

import (
	"bluebell/logic"
	"bluebell/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CreatePostHandler(c *gin.Context) {
	// 获取参数及参数校验
	p := new(models.Post)
	if err := c.ShouldBindJSON(p); err != nil {
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 从ctx里取得用户id
	userID, err := GetCurrentUserID(c)
	if err != nil {
		ResponseError(c, CodeNeedLogin)
	}
	p.AuthorID = userID
	// 创建帖子
	if err := logic.CreatePost(p); err != nil {
		zap.L().Error("post creating failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回
	ResponseSuccess(c, nil)
}

func GetPostDetailHandler(c *gin.Context) {
	// 获取帖子id
	pidStr := c.Param("id")
	pid, err := strconv.ParseInt(pidStr, 10, 64)
	if err != nil {
		zap.L().Error("get post detail with invalid parameters", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
	}
	// 根据id查找帖子详情数据
	data, err := logic.GetPostByID(pid)
	if err != nil {
		zap.L().Error("post getting failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
	}
	// 返回相应
	ResponseSuccess(c, data)
}
