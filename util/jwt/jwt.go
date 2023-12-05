package jwt

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const TokenExpireDuration = time.Hour * 2

var mySecret = []byte("来者Can追") // 私钥

type MyClaims struct {
	UserID             int64  `json:"user_id"`  //自定义user_id字段
	Username           string `json:"username"` // 自定义username字段
	jwt.StandardClaims        // JWT官方字段
}

func GenToken(userID int64, username string) (string, error) {
	// 创建自己的JWT负载
	c := MyClaims{
		userID,
		username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TokenExpireDuration).Unix(),
			Issuer:    "bluebell",
		},
	}
	// 使用指定的签名方法创建签名对象，即生成JWT头部和负载
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 对JWT头部和负载签名，生成完整Token
	return token.SignedString(mySecret)
}

func ParseToken(tokenString string) (*MyClaims, error) {
	var mc = new(MyClaims)
	token, err := jwt.ParseWithClaims(tokenString, mc, func(token *jwt.Token) (i interface{}, err error) {
		return mySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if token.Valid {
		return mc, nil
	}
	return nil, errors.New("invalid token")
}
