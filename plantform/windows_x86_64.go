//go:build windows && amd64
// +build windows,amd64

package plantform

import (
	_ "embed"
)

//go:embed cqhttp_brotli/go-cqhttp_windows_amd64.exe.brotli
var embedding_cqhttp []byte
var PLANTFORM = WINDOWS_x86_64
