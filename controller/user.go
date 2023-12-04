package controller

import (
	"bluebell/logic"
	"bluebell/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SignUpHandler(c *gin.Context) {
	// 获取参数并初步校验
	var p models.ParamSignUp
	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Error("sign up with invalid params", zap.Error(err))
		c.JSON(http.StatusOK, gin.H{
			"msg": err.Error(), // validator会将哪些有问题的字段写到err里
		})
		return
	}

	// 业务处理
	logic.SignUp(&p)

	// 返回响应
	c.JSON(http.StatusOK, gin.H{
		"msg": "sign up success",
	})
}
