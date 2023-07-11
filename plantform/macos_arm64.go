//go:build darwin && arm64
// +build darwin,arm64

package plantform

import (
	_ "embed"
)

//go:embed cqhttp_brotli/go-cqhttp_darwin_arm64.brotli
var embedding_cqhttp []byte
var PLANTFORM = MACOS_arm64
var OriginExecName = "omega_launcher_darwin_arm64"
