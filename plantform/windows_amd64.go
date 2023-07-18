//go:build windows && amd64
// +build windows,amd64

package plantform

import (
	_ "embed"
)

var (
	PLANTFORM       = WINDOWS_amd64
	launcherName    = "omega_launcher_windows_amd64.exe"
	cqhttpName      = "go-cqhttp_windows_amd64.exe"
	fastBuilderName = "phoenixbuilder-windows-executable-x86_64.exe"
	jdkDownloadName = "jdk-20_windows-x64_bin.zip"
)
