//go:build darwin && arm64
// +build darwin,arm64

package plantform

import (
	_ "embed"
)

var (
	PLANTFORM       = MACOS_arm64
	launcherName    = "omega_launcher_darwin_arm64"
	cqhttpName      = "go-cqhttp_darwin_arm64"
	fastBuilderName = "phoenixbuilder-macos-arm64"
	jdkDownloadName = "jdk-20_macos-aarch64_bin.tar.gz"
)
