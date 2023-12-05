package controller

import (
	"bluebell/dao/mysql"
	"bluebell/logic"
	"bluebell/models"
	"errors"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 用户注册
func SignUpHandler(c *gin.Context) {
	// 获取参数并校验
	var p models.ParamSignUp
	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Error("sign up with invalid parameters", zap.Error(err))
		ResponseErrorWithMsg(c, CodeInvalidParam, err.Error())
		return
	}
	// 业务处理
	if err := logic.SignUp(&p); err != nil {
		zap.L().Error("signing up failed", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserExist) {
			ResponseError(c, CodeUserExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSuccess(c, nil)
}

// 用户登录
func LoginHandler(c *gin.Context) {
	// 获取参数并校验
	var p models.ParamLogin
	if err := c.ShouldBindJSON(&p); err != nil {
		zap.L().Error("log in with invalid params", zap.Error(err))
		ResponseErrorWithMsg(c, CodeInvalidPassword, err.Error())
		return
	}
	// 业务处理
	if err := logic.Login(&p); err != nil {
		zap.L().Error("logging in failed", zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(c, CodeUserNotExist)
			return
		}
		if errors.Is(err, mysql.ErrorInvalidPassword) {
			ResponseError(c, CodeInvalidPassword)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}
	// 返回响应
	ResponseSuccess(c, nil)
}
