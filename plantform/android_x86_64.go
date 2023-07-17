//go:build android && amd64
// +build android,amd64

package plantform

import (
	_ "embed"
)

var (
	PLANTFORM      = Android_x86_64
	OriginExecName = "omega_launcher_android_amd64"
	CQHttpName     = "go-cqhttp_android_amd64"
)
