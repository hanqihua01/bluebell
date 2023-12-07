package mysql

import (
	"bluebell/models"
)

func CreatePost(p *models.Post) (err error) {
	sqlStr := `insert into post
	(post_id, title, content, author_id, community_id)
	values (?, ?, ?, ?, ?)`
	_, err = db.Exec(sqlStr, p.ID, p.Title, p.Content, p.AuthorID, p.CommunityID)
	return
}

func GetPostByID(pid int64) (data *models.Post, err error) {
	post := new(models.Post)
	sqlStr := "select post_id, title, content, author_id, community_id, create_time from post where post_id = ?"
	err = db.Get(post, sqlStr, pid)
	return post, err
}

func GetPostList(pageNum, pageSize int64) (posts []*models.Post, err error) {
	sqlStr := "select post_id, title, content, author_id, community_id, create_time from post limit ?,?"
	posts = make([]*models.Post, 0, pageSize)
	err = db.Select(&posts, sqlStr, (pageNum-1)*pageSize, pageSize)
	return
}
