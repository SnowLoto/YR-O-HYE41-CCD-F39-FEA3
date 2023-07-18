//go:build linux && !android && arm64
// +build linux,!android,arm64

package plantform

import (
	_ "embed"
)

var (
	PLANTFORM       = Linux_arm64
	launcherName    = "omega_launcher_linux_arm64"
	cqhttpName      = "go-cqhttp_linux_arm64"
	fastBuilderName = "phoenixbuilder-aarch64"
	jdkDownloadName = "jdk-20_linux-aarch64_bin.tar.gz"
)
