package controller

import (
	"bluebell/logic"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SignUpHandler(c *gin.Context) {
	// 参数校验
	// 业务处理
	logic.SignUp()
	// 返回响应
	c.JSON(http.StatusOK, "ok")
}
