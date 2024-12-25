package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"kubevulpes/api/httputils"
	option "kubevulpes/cmd/app/options"
	utilToken "kubevulpes/pkg/util/token"
)

// Authentication 身份认证
func Authentication(o *option.Options) gin.HandlerFunc {
	keyBytes := []byte(o.ComponentConfig.Default.JWTKey)

	return func(c *gin.Context) {
		if o.ComponentConfig.Default.Mode == "debug" {
			// Considered all as root user when running in debug mode.
			// todo GetRoot
			//root, err := o.Factory.User().GetRoot(c)
			//if err != nil {
			//	httputils.AbortFailedWithCode(c, http.StatusInternalServerError, err)
			//	return
			//}

			//httputils.SetUserToContext(c, root)
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
	_, err = utilToken.ParseToken(token, keyBytes)
	if err != nil {
		return err
	}

	// todo 优化多人登录问题
	//_, err = o.Controller.User().GetLoginToken(c, claim.Id)
	//if err != nil {
	//	return fmt.Errorf("未登陆或者密码被修改，请重新登陆")
	//}
	//
	//user, err := o.Factory.User().Get(c, claim.Id)
	//if err != nil {
	//	return err
	//}
	//if user == nil {
	//	return errors.ErrUnauthorized
	//}
	//httputils.SetUserToContext(c, user)

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
