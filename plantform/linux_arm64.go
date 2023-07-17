//go:build linux && !android && arm64
// +build linux,!android,arm64

package plantform

import (
	_ "embed"
)

var (
	PLANTFORM      = Linux_arm64
	OriginExecName = "omega_launcher_linux_arm64"
	CQHttpName     = "go-cqhttp_linux_arm64"
)
