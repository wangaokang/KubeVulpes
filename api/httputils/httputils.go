/*
Copyright 2024 The Vuples Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package httputils

import (
	"context"
	goerrors "errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"kubevulpes/api/errors"
	validatorutil "kubevulpes/api/validator"
	"kubevulpes/pkg/db/model"
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

func (r *Response) IsSuccessful() bool {
	return r.Code == http.StatusOK ||
		r.Code == http.StatusCreated ||
		r.Code == http.StatusAccepted
}

func (r *Response) SetMessageWithCode(m interface{}, c int) {
	r.SetCode(c)
	r.SetMessage(m)
}

func (r *Response) Error() string {
	return r.Message
}

func (r *Response) String() string {
	//data, _ := json.Marshal(r)
	//return string(data)
	return ""
}

// NewResponse 构造 http 返回值
// SetSuccess 时设置 code 为 200 并追加 success 的标识
// SetFailed 时设置 code 为 400，也可以自定义设置错误码，并追加报错信息
func NewResponse() *Response {
	return &Response{}
}

// SetSuccess 设置成功返回值
func SetSuccess(c *gin.Context, r *Response) {
	_ = contextBind(c).withResponseCode(http.StatusOK)
	r.SetMessageWithCode("success", http.StatusOK)
	c.JSON(http.StatusOK, r)
}

// SetFailed 设置错误返回值
func SetFailed(c *gin.Context, r *Response, err error) {
	switch e := err.(type) {
	case errors.Error:
		SetFailedWithCode(c, r, e.Code, e)
	case validator.ValidationErrors:
		SetFailedWithValidationError(c, r, validatorutil.TranslateError(e))
	default:
		SetFailedWithCode(c, r, http.StatusBadRequest, err)
	}
}

// SetFailedWithCode 设置错误返回值
func SetFailedWithCode(c *gin.Context, r *Response, code int, err error) {
	_ = contextBind(c).withResponseCode(code).withRawError(err)
	r.SetMessageWithCode(err, code)
	c.JSON(http.StatusOK, r)
}

func SetFailedWithValidationError(c *gin.Context, r *Response, e string) {
	_ = contextBind(c).withResponseCode(http.StatusBadRequest).withRawError(goerrors.New(e))
	r.SetMessageWithCode(e, http.StatusBadRequest)
	c.JSON(http.StatusOK, r)
}

//func SetLicenseFileToContext(c *gin.Context, saltTS string) {
//	c.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment;filename=.ctr"))
//	c.Writer.Write([]byte(saltTS))
//}

// AbortFailedWithCode 设置错误，code 返回值并终止请求
func AbortFailedWithCode(c *gin.Context, code int, err error) {
	r := NewResponse()
	_ = contextBind(c).withResponseCode(code).withRawError(err)
	r.SetMessageWithCode(err, code)
	c.JSON(http.StatusOK, r)
	c.Abort()
}

func ShouldBindAny(c *gin.Context, jsonObject interface{}, uriObject interface{}, queryObject interface{}) error {
	var err error
	if jsonObject != nil {
		if err = c.ShouldBindJSON(jsonObject); err != nil {
			return err
		}
	}
	if uriObject != nil {
		if err = c.ShouldBindUri(uriObject); err != nil {
			return err
		}
	}
	if queryObject != nil {
		if err = c.ShouldBindQuery(queryObject); err != nil {
			return err
		}
	}
	return nil
}

const (
	userKey  = "user"
	userName = "userName"
	project  = "project"
	company  = "company"
	serialId = "serialId"
	status   = "status"
)

func GetUserFromRequest(ctx context.Context) (*model.User, error) {
	val := ctx.Value(userKey)
	if val == nil {
		return nil, fmt.Errorf("get nil user")
	}

	user, ok := val.(*model.User)
	if !ok {
		return nil, fmt.Errorf("failed to assert user")
	}
	return user, nil
}

func GetUserIdFromContext(ctx context.Context) (int64, error) {
	user, err := GetUserFromRequest(ctx)
	if err != nil {
		return 0, err
	}
	return user.Id, nil
}

func SetUserToContext(c *gin.Context, user *model.User) {
	c.Set(userKey, user)
}

func SetStatusToContext(c *gin.Context, us *model.UserStatus) {
	c.Set(status, us)
}

func getStatusFromContext(ctx context.Context) (*model.UserStatus, error) {
	val := ctx.Value(status)
	if val == nil {
		return nil, fmt.Errorf("get nil status")
	}

	us, ok := val.(*model.UserStatus)
	if !ok {
		return nil, fmt.Errorf("failed to assert status")
	}

	return us, nil
}

func GetUserNameFromContext(ctx context.Context) (string, error) {
	val := ctx.Value(userName)
	if val == nil {
		return "", fmt.Errorf("get nil user")
	}

	user, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("failed to assert user")
	}
	return user, nil
}

func SetUserNameToContext(c *gin.Context, user string) {
	c.Set(userName, user)
}

//func SetProjectToContext(c *gin.Context, pName, comp string) {
//	c.Set(project, pName)
//	c.Set(company, comp)
//}

func GetProjectFromContext(ctx context.Context) (string, string, error) {
	val := ctx.Value(project)
	if val == nil {
		return "", "", fmt.Errorf("get nil project")
	}
	projectName, ok := val.(string)
	if !ok {
		return "", "", fmt.Errorf("failed to assert project")
	}
	com, ok := val.(string)
	if !ok {
		return "", "", fmt.Errorf("failed to assert company")
	}

	return projectName, com, nil
}

//func SetSerialIdToContext(c *gin.Context, sid string) {
//	c.Set(serialId, sid)
//}

func GetSerialIdFromContext(ctx context.Context) (string, error) {
	val := ctx.Value(serialId)
	if val == nil {
		return "", fmt.Errorf("get nil serialId")
	}
	serial, ok := val.(string)
	if !ok {
		return "", fmt.Errorf("failed to assert serialId")
	}
	return serial, nil
}

func GetObjectFromRequest(c *gin.Context) (string, string, bool) {
	return getObjectFromRequest(c.Request.URL.Path)
}

// getObjectFromRequest cuts and returns the object from the request path.
// e.g. /license/clusters/1 -> "clusters" "1" true

// todo 后续优化
func getObjectFromRequest(path string) (obj, sid string, ok bool) {
	// must start with /
	l := len(path)
	if l == 0 || path[0] != '/' {
		return
	}
	subs := strings.Split(path[1:l], "/")
	l = len(subs)
	if l < 2 || subs[0] != "license" {
		return subs[0], "", true
	}
	if l == 2 {
		// e.g. /license/clusters -> "clusters" "" true
		if subs[1] != "users" && subs[1] != "projects" && subs[1] != "audits" {
			return subs[0], subs[1], true
		}
		return subs[1], "", true
	}
	return subs[1], subs[2], true
}

// TranslateUrl todo 翻译 url 后续优化
func TranslateUrl(c *gin.Context) (model, action, content string, err error) {
	object, _, ok := getObjectFromRequest(c.Request.URL.Path)
	if !ok {
		return "", "", "", errors.ErrInvalidRequest
	}

	switch object {
	case "license":
		model = "license管理"
		serial, err := GetSerialIdFromContext(c)
		if err != nil {
			return "", "", "", err
		}
		content = fmt.Sprintf("证书序列号: %s", serial)
		switch c.Request.Method {
		case http.MethodPost:
			action = "新增授权证书"
		case http.MethodDelete:
			action = "删除证书"
		case http.MethodGet:
			if strings.HasSuffix(c.Request.URL.Path, "download") {
				action = "下载授权证书"
			}
		}
	case "users":
		model = "用户管理"
		email, err := GetUserNameFromContext(c)
		content = fmt.Sprintf("账号: %s", email)
		if err != nil {
			return "", "", "", err
		}
		switch c.Request.Method {
		case http.MethodPost:
			action = "创建账号"
			if strings.HasSuffix(c.Request.URL.Path, "password") {
				action = "重置密码"
			}
		case http.MethodPut:
			action = "编辑账号"
			if strings.HasSuffix(c.Request.URL.Path, "permissions") {
				if usersOperateMap(c) {
					action = "启用账号"
				} else {
					action = "禁用账号"
				}
			}
		case http.MethodDelete:
			action = "删除账号"
		case http.MethodGet:

		}
		if strings.HasSuffix(c.Request.URL.Path, "login") {
			model = "平台登录"
			action = "进入license管理平台"
		}
	case "projects":
		model = "license管理"
		projectName, comp, err := GetProjectFromContext(c)
		if err != nil {
			return "", "", "", err
		}
		content = fmt.Sprintf("企业: %s, 项目: %s", projectName, comp)
		switch c.Request.Method {
		case http.MethodPost:
			action = "新增项目"
		case http.MethodPut:
			action = "编辑项目"
		case http.MethodDelete:
			action = "删除项目"
		}
	default:
		return "", "", "", fmt.Errorf("unknown module")
	}
	return
}

func usersOperateMap(c *gin.Context) bool {
	if strings.HasSuffix(c.Request.URL.Path, "permissions") {
		s, err := getStatusFromContext(c)
		if err != nil {
			return false
		}
		if *s == model.UserStatusDisabled {
			return false
		}

	}
	return true
}

const (
	objIDsKey = "objIDs"
)

func SetIdRangeContext(c *gin.Context, ids []int64) {
	c.Set(objIDsKey, ids)
}

func GetIdRangeFromListReq(ctx context.Context) (exists bool, ids []int64) {
	val := ctx.Value(objIDsKey)
	if val == nil {
		return
	}

	ids, exists = val.([]int64)
	return
}

const (
	ResponseCodeKey = "response_code"
	RawErrorKey     = "raw_error"
)

type ctxBind struct {
	*gin.Context
}

func contextBind(c *gin.Context) *ctxBind {
	return &ctxBind{c}
}

// withResponseCode puts the response code into the HTTP context.
func (cb *ctxBind) withResponseCode(code int) *ctxBind {
	cb.Set(ResponseCodeKey, code)
	return cb
}

// withRawError puts the raw error into the HTTP context.
func (cb *ctxBind) withRawError(err error) *ctxBind {
	cb.Set(RawErrorKey, err)
	return cb
}

// GetResponseCode gets the response code from the HTTP context.
func GetResponseCode(ctx context.Context) (code int) {
	val := ctx.Value(ResponseCodeKey)
	if val == nil {
		return
	}

	code = val.(int)
	return
}

// GetRawError gets the raw error from the HTTP context.
func GetRawError(ctx context.Context) (err error) {
	val := ctx.Value(RawErrorKey)
	if val == nil {
		return
	}

	err = val.(error)
	return
}
