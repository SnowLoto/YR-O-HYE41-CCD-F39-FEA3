//go:build android && amd64
// +build android,amd64

package plantform

import (
	_ "embed"
)

var (
	PLANTFORM       = Android_amd64
	launcherName    = "omega_launcher_android_amd64"
	cqhttpName      = "go-cqhttp_android_amd64"
	fastBuilderName = "phoenixbuilder-android-termux-shared-executable-x86_64"
	jdkDownloadName = ""
)
