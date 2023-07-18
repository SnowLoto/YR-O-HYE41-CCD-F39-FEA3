//go:build linux && !android && amd64 && !tag_android_x86_64
// +build linux,!android,amd64,!tag_android_x86_64

package plantform

import (
	_ "embed"
)

var (
	PLANTFORM       = Linux_amd64
	launcherName    = "omega_launcher_linux_amd64"
	cqhttpName      = "go-cqhttp_linux_amd64"
	fastBuilderName = "phoenixbuilder"
	jdkDownloadName = "jdk-20_linux-x64_bin.tar.gz"
)
