package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/util/snowflake"

	"go.uber.org/zap"
)

func CreatePost(p *models.Post) (err error) {
	// 生成post id
	p.ID = snowflake.GenID()
	// 保存到数据库
	return mysql.CreatePost(p)
}

func GetPostByID(pid int64) (data *models.ApiPostDetail, err error) {
	post, err := mysql.GetPostByID(pid)
	if err != nil {
		zap.L().Error("looking up post failed", zap.Int64("pid", pid), zap.Error(err))
		return
	}
	user, err := mysql.GetUserByID(post.AuthorID)
	if err != nil {
		zap.L().Error("looking up user failed", zap.Int64("uid", post.AuthorID), zap.Error(err))
		return
	}
	communityDetail, err := mysql.GetCommunityDetailById(post.CommunityID)
	if err != nil {
		zap.L().Error("looking up community failed", zap.Int64("cid", post.CommunityID), zap.Error(err))
		return
	}
	data = &models.ApiPostDetail{
		AuthorName:      user.UserName,
		Post:            post,
		CommunityDetail: communityDetail,
	}
	return
}
