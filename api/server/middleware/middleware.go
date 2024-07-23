package middleware

import (
	"fmt"
	"k8s.io/apimachinery/pkg/util/sets"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/caoyingjunz/pixiu/cmd/app/config"
	"kubevulpes/api/server/httputils"
	tokenutil "kubevulpes/pkg/util/token"
)

const (
	// DebugMode indicates gin mode is debug.
	DebugMode = "debug"
)

var alwaysAllowPath sets.String

// Authentication 身份认证
func Authentication(cfg config.DefaultOptions) gin.HandlerFunc {
	keyBytes := []byte(cfg.JWTKey)
	mode := cfg.Mode

	return func(c *gin.Context) {
		if mode == DebugMode || alwaysAllowPath.Has(c.Request.URL.Path) || initAdminUser(c) {
			return
		}

		claim, err := validate(c, keyBytes)
		if err != nil {
			httputils.AbortFailedWithCode(c, http.StatusUnauthorized, err)
			return
		}

		c.Set("userId", claim.Id)
		c.Set("userName", claim.Name)
	}
}

func validate(c *gin.Context, keyBytes []byte) (*tokenutil.Claims, error) {
	token, err := extractToken(c, false)
	if err != nil {
		return nil, err
	}

	return tokenutil.ParseToken(token, keyBytes)
}

// 从请求头中获取 token
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
		return "", fmt.Errorf("unsupported authorization type")
	}

	return fields[1], nil
}

func initAdminUser(c *gin.Context) bool {
	return c.Request.Method == http.MethodPost && strings.HasPrefix(c.Request.URL.Path, "/pixiu/users") && c.Query("initAdmin") == "true"
}
