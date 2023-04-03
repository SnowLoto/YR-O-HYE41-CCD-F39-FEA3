//go:build windows && arm64
// +build windows,arm64

package plantform

import (
	_ "embed"
)

//go:embed cqhttp_brotli/go-cqhttp_windows_arm64.exe.brotli
var embedding_cqhttp []byte
var PLANTFORM = WINDOWS_arm64
