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

package validator

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zt "github.com/go-playground/validator/v10/translations/zh"
)

type customValidator interface {
	getTag() string
	translateError(ut ut.Translator) error
	translate(ut ut.Translator, fe validator.FieldError) string

	// Should be implemented by the custom validator.
	validate(fl validator.FieldLevel) bool
}

var tran ut.Translator
var customValidators []customValidator

func init() {
	_zh := zh.New() // default is Chinese
	uni := ut.New(_zh, _zh)
	tran, _ = uni.GetTranslator("zh")

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		_ = zt.RegisterDefaultTranslations(v, tran)

		for _, c := range customValidators {
			_ = v.RegisterValidation(c.getTag(), c.validate)
			_ = v.RegisterTranslation(c.getTag(), tran, c.translateError, c.translate)
		}
	}
}
