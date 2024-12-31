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
	"bytes"
	"context"
	"encoding/json"
	"k8s.io/klog/v2"
	"net/http"
	"strings"

	"github.com/gin-contrib/requestid"
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
	if c.Request.Method == http.MethodGet &&
		c.Writer.Status() != http.StatusUnauthorized {
		return
	}

	userName := "unknown"
	if user, err := httputils.GetUserFromRequest(c); err == nil && user != nil {
		userName = user.Name
	}

	obj, _, ok := httputils.GetObjectFromRequest(c)
	if !ok {
		return
	}

	audit := &model.Audit{
		RequestId:  requestid.Get(c),
		Action:     c.Request.Method,
		IP:         c.ClientIP(),
		Operator:   userName,
		Path:       c.Request.RequestURI,
		ObjectType: model.ObjectType(obj),
		Status:     getAuditStatus(c),
	}
	if err := w.opts.Factory.Audit().Create(context.TODO(), audit); err != nil {
		klog.Errorf("failed to create audit record [%s]: %v", audit.String(), err)
	}
}

// 允许特定请求不经过验证
func allowRequest(c *gin.Context) bool {
	// 用户请求
	if strings.HasSuffix(c.Request.URL.Path, "download") {
		return true
	}
	return false
}

// getAuditStatus returns the status of operation.
func getAuditStatus(c *gin.Context) model.AuditOperationStatus {
	respCode := httputils.GetResponseCode(c)
	if respCode == 0 {
		return model.AuditOpUnknown
	}

	if responseOK(respCode) {
		return model.AuditOpSuccess
	}

	return model.AuditOpFail
}

func responseOK(code int) bool {
	return code == http.StatusOK ||
		code == http.StatusCreated ||
		code == http.StatusAccepted
}
