package utils

import (
	"github.com/hashicorp/go-version"
)

// 是否有更高的版本
func HasGreaterVer(current, online string) bool {
	// 解析当前版本
	v1, err := version.NewVersion(current)
	if err != nil || v1 == nil {
		return false
	}
	// 解析在线版本
	v2, _ := version.NewVersion(online)
	if err != nil || v2 == nil {
		return false
	}
	// 版本对比
	return v1.LessThan(v2)
}
