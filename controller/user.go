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
			"msg": "params wrong",
		})
		return
	}

	// 手动校验请求参数，c.ShouldBindJson()非常蠢，但是如果结合了validator就很牛
	// if len(p.Username) == 0 || len(p.Password) == 0 || len(p.RePassword) == 0 || p.Password != p.RePassword {
	// 	zap.L().Error("sign up with invalid params")
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"msg": "params wrong",
	// 	})
	// 	return
	// }

	// 业务处理
	logic.SignUp(&p)

	// 返回响应
	c.JSON(http.StatusOK, gin.H{
		"msg": "sign up success",
	})
}
