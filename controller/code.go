package controller

type ResCode int64

const (
	CodeSuccess ResCode = 1000 + iota
	CodeInvalidParam
	CodeUserExist
	CodeUserNotExist
	CodeInvalidPassword
	CodeServerBusy

	CodeNeedLogin
	CodeInvalidToken
)

var codeMsgMap = map[ResCode]string{
	CodeSuccess:         "success",
	CodeInvalidParam:    "invalid parameter",
	CodeUserExist:       "user already exists",
	CodeUserNotExist:    "user not exists",
	CodeInvalidPassword: "invalid password",
	CodeServerBusy:      "server is busy",

	CodeNeedLogin:    "need log in",
	CodeInvalidToken: "invalid token",
}

func (c ResCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[CodeServerBusy]
	}
	return msg
}
