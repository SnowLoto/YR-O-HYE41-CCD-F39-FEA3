//go:build linux && !android && amd64 && !tag_android_x86_64
// +build linux,!android,amd64,!tag_android_x86_64

package plantform

import (
	_ "embed"
)

var (
	PLANTFORM      = Linux_x86_64
	OriginExecName = "omega_launcher_linux_amd64"
	CQHttpName     = "go-cqhttp_linux_amd64"
)
