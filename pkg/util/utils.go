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

package util

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

// 定义一个常量数组，存储单位
const (
	B  = "B"
	KB = "KB"
	MB = "MB"
	GB = "GB"
	TB = "TB"
	PB = "PB"
)

// sizeUnitArray 包含字节大小的单位，用于格式化输出
var sizeUnitArray = [...]string{B, KB, MB, GB, TB, PB}

// EncryptUserPassword 生成加密密码
// 前端传的密码为明文，需要加密存储
func EncryptUserPassword(origin string) (string, error) {
	pwd, err := bcrypt.GenerateFromPassword([]byte(origin), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(pwd), nil
}

// ValidateUserPassword 验证用户的密码是否正确
func ValidateUserPassword(old, new string) error {
	return bcrypt.CompareHashAndPassword([]byte(old), []byte(new))
}

// formatSize 将字节大小转换为更友好的格式
func formatSize(size int64, unitArray [6]string) string {
	sizeF := float64(size)
	for _, unit := range unitArray {
		if sizeF < 1024 {
			// 使用 "%.2f" 格式化浮点数，保留两位小数
			return fmt.Sprintf("%.2f %s", sizeF, unit)
		}
		sizeF = sizeF / 1024
	}
	// 如果超出预定义单位，使用最高单位并保留四位小数
	return fmt.Sprintf("%.4f %s", sizeF, unitArray[len(unitArray)-1])
}

// MultiSizeConvert 将两个字节大小的整数值转换为可读格式
func MultiSizeConvert(size1, size2 int64) (string, string) {
	return formatSize(size1, sizeUnitArray), formatSize(size2, sizeUnitArray)
}

// DeduplicateIntSlice returns a new slice with duplicated elements removed.
func DeduplicateIntSlice(s []int64) (ret []int64) {
	ret = make([]int64, 0)
	m := make(map[int64]struct{})
	for _, v := range s {
		if _, ok := m[v]; ok {
			continue
		}
		m[v] = struct{}{}
		ret = append(ret, v)
	}

	return
}

// Less returns the smxaller one.
func Less(a, b int) int {
	if a < b {
		return a
	}
	return b
}
