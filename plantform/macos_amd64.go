//go:build darwin && amd64
// +build darwin,amd64

package plantform

import (
	_ "embed"
)

var (
	PLANTFORM       = MACOS_amd64
	launcherName    = "omega_launcher_darwin_amd64"
	cqhttpName      = "go-cqhttp_darwin_amd64"
	fastBuilderName = "phoenixbuilder-macos-x86_64"
	jdkDownloadName = "jdk-20_macos-x64_bin.tar.gz"
)
