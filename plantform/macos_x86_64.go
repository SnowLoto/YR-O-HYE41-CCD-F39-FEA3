//go:build darwin && amd64
// +build darwin,amd64

package plantform

import (
	_ "embed"
)

var (
	PLANTFORM      = MACOS_x86_64
	OriginExecName = "omega_launcher_darwin_amd64"
	CQHttpName     = "go-cqhttp_darwin_amd64"
)
