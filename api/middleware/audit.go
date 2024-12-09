package middleware

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"kubevulpes/api/httputils"
	option "kubevulpes/cmd/app/options"
	"kubevulpes/pkg/db/model"
)

// 自定义 ResponseWriter 用于捕获写入的数据
type auditWriter struct {
	gin.ResponseWriter
	resp *httputils.Response
	opts *option.Options
}

func newResponseWriter(w gin.ResponseWriter, o *option.Options) *auditWriter {
	return &auditWriter{
		ResponseWriter: w,
		resp:           httputils.NewResponse(),
		opts:           o,
	}
}

func (w *auditWriter) Write(b []byte) (int, error) {
	_ = json.NewDecoder(bytes.NewReader(b)).Decode(w.resp)
	return w.ResponseWriter.Write(b)
}

func Audit(o *option.Options) gin.HandlerFunc {
	return func(c *gin.Context) {
		auditor := newResponseWriter(c.Writer, o)
		c.Writer = auditor
		c.Next()

		// do audit asynchronously
		go auditor.asyncAudit(c)
	}
}

// asyncAudit audits the request asynchronously.
// It should be called in a goroutine.
func (w *auditWriter) asyncAudit(c *gin.Context) {
	if c.Request.Method == http.MethodGet && (allowRequest(c) ||
		c.Writer.Status() != http.StatusUnauthorized) {
		return
	}
	user, _ := httputils.GetUserFromRequest(c)
	// request with invalid token,here cant get user,just return
	if user == nil {
		return
	}

	_, action, _, err := httputils.TranslateUrl(c)
	if err != nil {
		return
	}

	var status model.AuditOperationStatus = 0
	if w.resp != nil && w.resp.IsSuccessful() {
		status = 1
	}

	_ = &model.Audit{
		Operator: user.Name,
		//Email:         user.Email,
		Action: action,
		//IP:            c.ClientIP(),
		Status: status,
		//Content:       content,
		//ResourceModel: resourceModel,
	}
	//if _, err := w.opts.Factory.Audit().Create(context.TODO(), audit); err != nil {
	//	klog.Errorf("failed to create audit record [%s]: %v", audit.String(), err)
	//}
}

// 允许特定请求不经过验证
func allowRequest(c *gin.Context) bool {
	// 用户请求
	if strings.HasSuffix(c.Request.URL.Path, "download") {
		return true
	}
	return false
}
