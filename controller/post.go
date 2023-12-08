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

func GetPostListHandler(c *gin.Context) {
	// 获取分页参数
	pageNumStr := c.Query("page")
	pageSizeStr := c.Query("size")
	pageNum, err := strconv.ParseInt(pageNumStr, 10, 64)
	if err != nil {
		pageNum = 1
	}
	pageSize, err := strconv.ParseInt(pageSizeStr, 10, 64)
	if err != nil {
		pageSize = 1
	}
	// 获取数据
	data, err := logic.GetPostList(pageNum, pageSize)
	if err != nil {
		zap.L().Error("post list getting failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSuccess(c, data)
}

// 根据前端传过来的参数（按分数、按创建时间）动态的获取帖子列表
func GetPostListHandler2(c *gin.Context) {
	// GET请求参数：/api/v1/posts2?page=1&size=10&order=time -> query string参数
	// 获取分页参数
	// 指定默认参数
	p := &models.ParamPostList{
		Page:  1,
		Size:  10,
		Order: models.OrderTime,
	}
	// c.ShouldBind() 动态的根据请求参数的位置自动绑定
	if err := c.ShouldBindQuery(p); err != nil {
		zap.L().Error("get post list with invalid parameters", zap.Error(err))
		ResponseError(c, CodeInvalidParam)
		return
	}
	// 获取数据
	data, err := logic.GetPostListNew(p)
	if err != nil {
		zap.L().Error("post list getting failed", zap.Error(err))
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSuccess(c, data)
}

// 根据社区查询帖子列表
// func GetCommunityPostListHandler(c *gin.Context) {
// 	// GET请求参数：/api/v1/posts2?page=1&size=10&order=time -> query string参数
// 	// 获取分页参数
// 	// 指定默认参数
// 	p := &models.ParamCommunityPostList{
// 		Page:  1,
// 		Size:  10,
// 		Order: models.OrderTime,
// 	}
// 	// c.ShouldBind() 动态的根据请求参数的位置自动绑定
// 	if err := c.ShouldBindQuery(p); err != nil {
// 		zap.L().Error("get post list with invalid parameters", zap.Error(err))
// 		ResponseError(c, CodeInvalidParam)
// 		return
// 	}
// 	// 获取数据
// 	data, err := logic.GetCommunityPostList(p)
// 	if err != nil {
// 		zap.L().Error("post list getting failed", zap.Error(err))
// 		ResponseError(c, CodeServerBusy)
// 		return
// 	}
// 	// 返回响应
// 	ResponseSuccess(c, data)
// }
