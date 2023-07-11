//go:build linux && !android && arm64
// +build linux,!android,arm64

package plantform

import (
	_ "embed"
)

//go:embed cqhttp_brotli/go-cqhttp_linux_arm64.brotli
var embedding_cqhttp []byte
var PLANTFORM = Linux_arm64
var OriginExecName = "omega_launcher_linux_arm64"
