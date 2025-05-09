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
	"github.com/gin-contrib/requestid"
	"k8s.io/apimachinery/pkg/util/sets"

	option "kubevulpes/cmd/app/options"
	"kubevulpes/pkg/util"
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
		Logger(&o.ComponentConfig.Default.LogOptions),
		UserRateLimiter(),
		Limiter(),
		Authentication(o),
		Authorization(o),
		Audit(o))
}
