package logic

import (
	"bluebell/dao/mysql"
	"bluebell/dao/redis"
	"bluebell/models"
	"bluebell/util/snowflake"

	"go.uber.org/zap"
)

func CreatePost(p *models.Post) (err error) {
	// 生成post id
	p.ID = snowflake.GenID()
	// 保存到数据库
	err = mysql.CreatePost(p)
	if err != nil {
		return err
	}
	err = redis.CreatePost(p.ID, p.CommunityID)
	return
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

func GetPostList(pageNum, pageSize int64) (data []*models.ApiPostDetail, err error) {
	posts, err := mysql.GetPostList(pageNum, pageSize)
	if err != nil {
		return nil, err
	}
	data = make([]*models.ApiPostDetail, 0, len(posts))
	for _, post := range posts {
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("looking up user failed", zap.Int64("author_id", post.AuthorID), zap.Error(err))
			continue
		}
		communityDetail, err := mysql.GetCommunityDetailById(post.CommunityID)
		if err != nil {
			zap.L().Error("looking up community failed", zap.Int64("community_id", post.CommunityID), zap.Error(err))
			continue
		}
		postdetail := &models.ApiPostDetail{
			AuthorName:      user.UserName,
			Post:            post,
			CommunityDetail: communityDetail,
		}
		data = append(data, postdetail)
	}
	return
}

func GetPostList2(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// redis查询id列表
	ids, err := redis.GetPostIdsInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		return
	}
	// 根据id去mysql查询帖子详情信息
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}
	// 查询每篇帖子的赞成票数量
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}
	data = make([]*models.ApiPostDetail, 0, len(posts))
	for idx, post := range posts {
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("looking up user failed", zap.Int64("author_id", post.AuthorID), zap.Error(err))
			continue
		}
		communityDetail, err := mysql.GetCommunityDetailById(post.CommunityID)
		if err != nil {
			zap.L().Error("looking up community failed", zap.Int64("community_id", post.CommunityID), zap.Error(err))
			continue
		}
		postdetail := &models.ApiPostDetail{
			AuthorName:      user.UserName,
			Post:            post,
			CommunityDetail: communityDetail,
		}
		postdetail.VoteNum = voteData[idx]
		data = append(data, postdetail)
	}
	return
}

func GetCommunityPostList(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// redis查询id列表
	ids, err := redis.GetCommunityPostIdsInOrder(p)
	if err != nil {
		return
	}
	if len(ids) == 0 {
		return
	}
	// 根据id去mysql查询帖子详情信息
	posts, err := mysql.GetPostListByIDs(ids)
	if err != nil {
		return
	}
	// 查询每篇帖子的赞成票数量
	voteData, err := redis.GetPostVoteData(ids)
	if err != nil {
		return
	}
	data = make([]*models.ApiPostDetail, 0, len(posts))
	for idx, post := range posts {
		user, err := mysql.GetUserByID(post.AuthorID)
		if err != nil {
			zap.L().Error("looking up user failed", zap.Int64("author_id", post.AuthorID), zap.Error(err))
			continue
		}
		communityDetail, err := mysql.GetCommunityDetailById(post.CommunityID)
		if err != nil {
			zap.L().Error("looking up community failed", zap.Int64("community_id", post.CommunityID), zap.Error(err))
			continue
		}
		postdetail := &models.ApiPostDetail{
			AuthorName:      user.UserName,
			Post:            post,
			CommunityDetail: communityDetail,
		}
		postdetail.VoteNum = voteData[idx]
		data = append(data, postdetail)
	}
	return
}

func GetPostListNew(p *models.ParamPostList) (data []*models.ApiPostDetail, err error) {
	// 如果传来了communityid，则说明仅查看该community的posts
	if p.CommunityID == 0 {
		data, err = GetPostList2(p)
	} else {
		data, err = GetCommunityPostList(p)
	}
	if err != nil {
		return nil, err
	}
	return
}
