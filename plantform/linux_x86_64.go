//go:build linux && !android && amd64 && !tag_android_x86_64
// +build linux,!android,amd64,!tag_android_x86_64

package plantform

import (
	_ "embed"
)

//go:embed cqhttp_brotli/go-cqhttp_linux_amd64.brotli
var embedding_cqhttp []byte
var PLANTFORM = Linux_x86_64
var OriginExecName = "omega_launcher_linux_amd64"
