package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	CodeSuccess        = 200
	CodeBadRequest     = 40001
	CodeUnauthorized   = 40101
	CodeTokenInvalid   = 40102
	CodeUserNotFound   = 40002
	CodeUserForbidden  = 40301
	CodeUsernameExists = 40901
)

var msgMap = map[int]string{
	CodeSuccess:        "success",
	CodeBadRequest:     "参数错误",
	CodeUnauthorized:   "未登录",
	CodeTokenInvalid:   "Token 无效或过期",
	CodeUserNotFound:   "用户不存在或密码错误",
	CodeUserForbidden:  "账号被禁用",
	CodeUsernameExists: "用户名已存在",
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func JSON(c *gin.Context, code int, data interface{}) {
	msg := msgMap[code]
	if msg == "" {
		msg = "unknown error"
	}
	c.JSON(http.StatusOK, Response{Code: code, Message: msg, Data: data})
}

func OK(c *gin.Context, data interface{}) {
	JSON(c, CodeSuccess, data)
}

func Error(c *gin.Context, code int) {
	JSON(c, code, nil)
}
