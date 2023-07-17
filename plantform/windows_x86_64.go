//go:build windows && amd64
// +build windows,amd64

package plantform

import (
	_ "embed"
)

var (
	PLANTFORM      = WINDOWS_x86_64
	OriginExecName = "omega_launcher_windows_amd64.exe"
	CQHttpName     = "go-cqhttp_windows_amd64.exe"
)
