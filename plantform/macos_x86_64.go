//go:build darwin && amd64
// +build darwin,amd64

package plantform

import (
	_ "embed"
)

//go:embed cqhttp_brotli/go-cqhttp_darwin_amd64.brotli
var embedding_cqhttp []byte
var PLANTFORM = MACOS_x86_64
