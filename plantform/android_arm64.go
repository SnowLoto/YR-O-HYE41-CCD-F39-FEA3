//go:build android && arm64
// +build android,arm64

package plantform

import (
	_ "embed"
)

var (
	PLANTFORM      = Android_arm64
	OriginExecName = "omega_launcher_android_arm64"
	CQHttpName     = "go-cqhttp_android_arm64"
)
