package router

import (
	"bluebell/controller"
	"bluebell/util/logger"
	"bluebell/util/ratelimit"
	"net/http"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func Setup(mode string) *gin.Engine {
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New() // 不使用gin.Default()，因为要使用自定义的logger和recovery中间件
	r.Use(logger.GinLogger(), logger.GinRecovery(true), ratelimit.RateLimitMiddleware(2*time.Second, 1))

	// 注册业务路由
	v1 := r.Group("/api/v1")
	v1.POST("/signup", controller.SignUpHandler)
	v1.POST("/login", controller.LoginHandler)

	v1.Use(controller.JWTAuthMiddleware()) // 登录后才需要JWT认证

	{
		v1.GET("/community", controller.CommunityHandler)
		v1.GET("/community/:id", controller.CommunityDetailHandler)

		v1.POST("/post", controller.CreatePostHandler)
		v1.GET("/posts", controller.GetPostListHandler)
		v1.GET("/posts2", controller.GetPostListHandler2)
		v1.GET("/post/:id", controller.GetPostDetailHandler)

		v1.POST("/vote", controller.PostVoteHandler)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "no such route",
		})
	})

	pprof.Register(r)

	return r
}
