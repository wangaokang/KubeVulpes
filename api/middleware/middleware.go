package middleware

import (
	"github.com/caoyingjunz/pixiu/pkg/util"
	"github.com/gin-contrib/requestid"
	"k8s.io/apimachinery/pkg/util/sets"

	option "kubevulpes/cmd/app/options"
)

var alwaysAllowPath sets.String

func init() {
	alwaysAllowPath = sets.NewString("/api/v1/users/login")
}

func InstallMiddlewares(o *option.Options) {
	o.HttpEngine.Use(
		requestid.New(requestid.WithGenerator(func() string {
			return util.GenerateRequestID()
		})),
		Cors(),
		Authentication(o),
		Audit(o))
}
