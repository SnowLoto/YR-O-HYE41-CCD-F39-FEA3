//go:build android && arm64
// +build android,arm64

package plantform

import (
	_ "embed"
)

var (
	PLANTFORM       = Android_arm64
	launcherName    = "omega_launcher_android_arm64"
	cqhttpName      = "go-cqhttp_android_arm64"
	fastBuilderName = "phoenixbuilder-android-termux-shared-executable-arm64"
	jdkDownloadName = ""
)
