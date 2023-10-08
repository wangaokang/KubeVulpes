package httputils

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code    int         `json:"code"`              // 返回的状态码
	Result  interface{} `json:"result,omitempty"`  // 正常返回时的数据，可以为任意数据结构
	Message string      `json:"message,omitempty"` // 异常返回时的错误信息
}

func (r *Response) SetCode(c int) {
	r.Code = c
}

func (r *Response) SetMessage(m interface{}) {
	switch msg := m.(type) {
	case error:
		r.Message = msg.Error()
	case string:
		r.Message = msg
	}
}

func (r *Response) SetMessageWithCode(m interface{}, c int) {
	r.SetCode(c)
	r.SetMessage(m)
}

func NewResponse() *Response {
	return &Response{}
}

func SetFailed(c *gin.Context, r *Response, err error) {
	SetFailedWithCode(c, r, http.StatusBadRequest, err)
}

// SetFailedWithCode 设置错误返回值
func SetFailedWithCode(c *gin.Context, r *Response, code int, err error) {
	r.SetMessageWithCode(err, code)
	c.JSON(http.StatusOK, r)
}

// SetSuccess 设置成功返回值
func SetSuccess(c *gin.Context, r *Response) {
	r.SetMessageWithCode("success", http.StatusOK)
	c.JSON(http.StatusOK, r)
}

// AbortFailedWithCode 设置错误，code 返回值并终止请求
func AbortFailedWithCode(c *gin.Context, code int, err error) {
	r := NewResponse()
	r.SetMessageWithCode(err, code)
	c.JSON(http.StatusOK, r)
	c.Abort()
}
