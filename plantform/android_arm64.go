//go:build android && arm64
// +build android,arm64

package plantform

import (
	_ "embed"
)

//go:embed cqhttp_brotli/go-cqhttp_android_arm64.brotli
var embedding_cqhttp []byte
var PLANTFORM = Android_arm64
