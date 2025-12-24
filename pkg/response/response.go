package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 定义统一的 JSON 返回结构
type Response struct {
	Code int         `json:"code"` // 业务码：0成功，非0失败
	Msg  string      `json:"msg"`  // 提示信息
	Data interface{} `json:"data"` // 数据
}

const (
	SUCCESS = 0    // 成功码
	ERROR   = 7    // 通用错误码
)

// Result 基础方法
func Result(c *gin.Context, httpCode int, code int, msg string, data interface{}) {
	c.JSON(httpCode, Response{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}

// Success 成功返回 (最常用)
// 用法: response.Success(c, map[string]interface{}{"name": "zy"})
func Success(c *gin.Context, data interface{}) {
	Result(c, http.StatusOK, SUCCESS, "success", data)
}

// Fail 失败返回
// 用法: response.Fail(c, "参数错误")
func Fail(c *gin.Context, msg string) {
	Result(c, http.StatusOK, ERROR, msg, nil)
}

// FailWithCode 失败返回（自定义错误码）
// 用法: response.FailWithCode(c, 10001, "余额不足")
func FailWithCode(c *gin.Context, code int, msg string) {
	Result(c, http.StatusOK, code, msg, nil)
}