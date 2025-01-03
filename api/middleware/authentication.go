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

package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"kubevulpes/api/errors"
	"kubevulpes/api/httputils"
	option "kubevulpes/cmd/app/options"
	ctrlutil "kubevulpes/pkg/controller/util"
	"kubevulpes/pkg/db/model"
	utilToken "kubevulpes/pkg/util/token"
)

// Authentication 身份认证
func Authentication(o *option.Options) gin.HandlerFunc {
	keyBytes := []byte(o.ComponentConfig.Default.JWTKey)

	return func(c *gin.Context) {
		if o.ComponentConfig.Default.Mode == "debug" {

			root, err := o.Factory.User().GetRoot(c)
			if err != nil {
				httputils.AbortFailedWithCode(c, http.StatusInternalServerError, err)
				return
			}

			httputils.SetUserToContext(c, root)
			return
		}

		if alwaysAllowPath.Has(c.Request.URL.Path) {
			return
		}

		if err := validate(c, o, keyBytes); err != nil {
			httputils.AbortFailedWithCode(c, http.StatusUnauthorized, err)
			return
		}
	}
}

func validate(c *gin.Context, o *option.Options, keyBytes []byte) error {
	token, err := extractToken(c, false)
	if err != nil {
		return err
	}

	claim, err := utilToken.ParseToken(token, keyBytes)
	if err != nil {
		return err
	}

	_, err = o.Controller.User().GetLoginToken(c, claim.Id)
	if err != nil {
		return fmt.Errorf("未登陆或者密码被修改，请重新登陆")
	}

	user, err := o.Factory.User().Get(c, claim.Id)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.ErrUnauthorized
	}
	httputils.SetUserToContext(c, user)

	return nil
}

// 从请求头中获取 tokxzxen
func extractToken(c *gin.Context, ws bool) (string, error) {
	emptyFunc := func(t string) bool { return len(t) == 0 }
	if ws {
		wsToken := c.GetHeader("Sec-WebSocket-Protocol")
		if emptyFunc(wsToken) {
			return "", fmt.Errorf("authorization header is not provided")
		}
		return wsToken, nil
	}

	token := c.GetHeader("Authorization")
	if emptyFunc(token) {
		return "", fmt.Errorf("authorization header is not provided")
	}
	fields := strings.Fields(token)
	if len(fields) != 2 {
		return "", fmt.Errorf("invalid authorization header format")
	}
	if fields[0] != "Bearer" {
		return "", fmt.Errorf("unsupported authorization types")
	}

	return fields[1], nil
}

// HTTP method to operation
var operationsMap = map[string]model.Operation{
	http.MethodGet:    model.OpRead,
	http.MethodPost:   model.OpCreate,
	http.MethodPatch:  model.OpUpdate,
	http.MethodPut:    model.OpUpdate,
	http.MethodDelete: model.OpDelete,
}

// Authorization 鉴权
func Authorization(o *option.Options) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 允许请求直接通过
		if o.ComponentConfig.Default.InDebug() || alwaysAllowPath.Has(c.Request.URL.Path) {
			return
		}

		user, err := httputils.GetUserFromRequest(c)
		if err != nil {
			httputils.AbortFailedWithCode(c, http.StatusMethodNotAllowed, err)
			return
		}

		switch user.Status {
		case 1:
			// 禁用用户无法进行任何操作
			httputils.AbortFailedWithCode(c, http.StatusForbidden, fmt.Errorf("用户已被禁用"))
			return
		}

		obj, id, ok := httputils.GetObjectFromRequest(c)
		if !ok {
			return
		}

		op := operationsMap[c.Request.Method]
		// load policy for consistency
		// ref: https://github.com/casbin/casbin/issues/679#issuecomment-761525328
		if err := o.Enforcer.LoadPolicy(); err != nil {
			httputils.AbortFailedWithCode(c, http.StatusInternalServerError, err)
			return
		}
		ok, err = o.Enforcer.Enforce(user.Name, obj, id, op.String())
		if err != nil {
			httputils.AbortFailedWithCode(c, http.StatusMethodNotAllowed, err)
			return
		}
		if !ok {
			httputils.AbortFailedWithCode(c, http.StatusForbidden, fmt.Errorf("无操作权限"))
		}
		if id != "" {
			return
		}
		// this is a list API
		if err := ctrlutil.SetIdRangeContext(c, o.Enforcer, user, obj); err != nil {
			httputils.AbortFailedWithCode(c, http.StatusInternalServerError, err)
		}
	}
}
