package router

import (
	"bluebell/controller"
	"bluebell/util/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Setup(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New() // 不使用gin.Default()，因为要使用自定义的logger和recovery中间件
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	// 注册业务路由
	r.POST("/signup", controller.SignUpHandler)
	r.POST("/login", controller.LoginHandler)

	r.GET("/ping", controller.JWTAuthMiddleware(), func(c *gin.Context) {
		isLogin := true
		if isLogin {
			c.String(http.StatusOK, "pong")
		} else {
			c.String(http.StatusOK, "please log in")
		}
	})

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "no such route",
		})
	})

	return r
}
