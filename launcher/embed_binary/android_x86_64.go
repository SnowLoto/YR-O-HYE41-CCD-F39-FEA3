//go:build android && amd64
// +build android,amd64

package embed_binary

import (
	_ "embed"
)

//go:embed assets/go-cqhttp_android_amd64.brotli
var embedding_cqhttp []byte
var PLANTFORM = Android_x86_64
