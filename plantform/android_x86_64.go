//go:build android && amd64
// +build android,amd64

package plantform

import (
	_ "embed"
)

//go:embed cqhttp_brotli/go-cqhttp_android_amd64.brotli
var embedding_cqhttp []byte
var PLANTFORM = Android_x86_64
