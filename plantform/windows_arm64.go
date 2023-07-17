//go:build windows && arm64
// +build windows,arm64

package plantform

import (
	_ "embed"
)

var (
	PLANTFORM      = WINDOWS_arm64
	OriginExecName = "omega_launcher_windows_arm64.exe"
	CQHttpName     = "go-cqhttp_windows_arm64.exe"
)
