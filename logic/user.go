package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/util/jwt"
	"bluebell/util/snowflake"
)

func SignUp(p *models.ParamSignUp) (err error) {
	// 判断用户存不存在
	if err := mysql.CheckUserExist(p.Username); err != nil {
		return err
	}
	// 生成UID
	userID := snowflake.GenID()
	// 构造一个User实例
	user := models.User{
		UserID:   userID,
		UserName: p.Username,
		Password: p.Password,
	}
	// 保存进数据库
	return mysql.InsertUser(&user)
}

func Login(p *models.ParamLogin) (user *models.User, err error) {
	user = &models.User{
		UserName: p.Username,
		Password: p.Password,
	}
	// 传递user指针，查询成功会user.ID将会有值
	if err := mysql.Login(user); err != nil {
		return nil, err
	}
	// 生成JWT
	token, err := jwt.GenToken(user.UserID, user.UserName)
	if err != nil {
		return nil, err
	}
	user.Token = token
	return
}
