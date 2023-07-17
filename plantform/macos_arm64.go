//go:build darwin && arm64
// +build darwin,arm64

package plantform

import (
	_ "embed"
)

var (
	PLANTFORM      = MACOS_arm64
	OriginExecName = "omega_launcher_darwin_arm64"
	CQHttpName     = "go-cqhttp_darwin_arm64"
)
