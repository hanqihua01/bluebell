package mysql

import (
	"bluebell/models"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
)

var (
	ErrorUserExist       = errors.New("user already exists")
	ErrorUserNotExist    = errors.New("user not exists")
	ErrorInvalidPassword = errors.New("password wrong")
)

// CheckUserExist 检查指定用户名的用户是否存在
func CheckUserExist(username string) (err error) {
	sqlStr := "select count(user_id) from user where username = ?"
	var count int
	if err := db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	if count > 0 {
		return ErrorUserExist
	}
	return nil
}

// InsertUser 向数据库中插入一条新的用户记录
func InsertUser(user *models.User) (err error) {
	// 对密码进行加密
	user.Password = encryptPassword(user.Password)
	sqlStr := "insert into user(user_id, username, password) values(?,?,?)"
	if _, err := db.Exec(sqlStr, user.UserID, user.UserName, user.Password); err != nil {
		return err
	}
	return nil
}

const secret = "hanqihua"

func encryptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}

func Login(user *models.User) (err error) {
	oPassword := user.Password // 保存传入的password，因为db.Get会把查询结果覆盖掉user
	sqlStr := "select user_id, username, password from user where username = ?"
	err = db.Get(user, sqlStr, user.UserName)
	if err == sql.ErrNoRows { // 用户不存在
		return ErrorUserNotExist
	}
	if err != nil { // 查询数据库错误
		return err
	}
	password := encryptPassword(oPassword)
	if password != user.Password {
		return ErrorInvalidPassword
	}
	return nil
}

func GetUserByID(uid int64) (user *models.User, err error) {
	user = new(models.User)
	sqlStr := "select user_id, username from user where user_id = ?"
	err = db.Get(user, sqlStr, uid)
	return
}
