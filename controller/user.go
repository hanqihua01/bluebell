package controller

import (
	"bluebell/logic"
	"bluebell/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SignUpHandler(c *gin.Context) {
	// 获取参数并校验
	var p models.ParamSignUp
	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Error("sign up with invalid params", zap.Error(err))
		c.JSON(http.StatusOK, gin.H{
			"msg": err.Error(), // validator会将哪些有问题的字段写到err里
		})
		return
	}
	// 业务处理
	if err := logic.SignUp(&p); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": "sign up fail",
		})
		return
	}
	// 返回响应
	c.JSON(http.StatusOK, gin.H{
		"msg": "sign up success",
	})
}

func LoginHandler(c *gin.Context) {
	// 获取参数并校验
	var p models.ParamLogin
	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Error("log in with invalid params", zap.Error(err))
		c.JSON(http.StatusOK, gin.H{
			"msg": err.Error(), // validator会将哪些有问题的字段写到err里
		})
		return
	}
	// 业务处理
	if err := logic.Login(&p); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": "log in fail: username of password wrong",
		})
		return
	}
	// 返回响应
	c.JSON(http.StatusOK, gin.H{
		"msg": "log in success",
	})
}
