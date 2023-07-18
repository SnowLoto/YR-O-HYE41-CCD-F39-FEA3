//go:build windows && arm64
// +build windows,arm64

package plantform

import (
	_ "embed"
)

var (
	PLANTFORM       = WINDOWS_arm64
	launcherName    = "omega_launcher_windows_arm64.exe"
	cqhttpName      = "go-cqhttp_windows_arm64.exe"
	fastBuilderName = ""
	jdkDownloadName = ""
)
