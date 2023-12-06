package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
)

func GetCommunityList() (data []*models.Community, err error) {
	return mysql.GetCommunityList()
}

func GetCommunityDetail(id int64) (communityDetail *models.CommunityDetail, err error) {
	return mysql.GetCommunityDetailById(id)
}
