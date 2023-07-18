package plantform

const (
	Android_arm64 = "android_arm64"
	Android_amd64 = "android_amd64"
	Linux_arm64   = "linux_arm64"
	Linux_amd64   = "linux_amd64"
	MACOS_arm64   = "macos_arm64"
	MACOS_amd64   = "macos_amd64"
	WINDOWS_arm64 = "windows_arm64"
	WINDOWS_amd64 = "windows_amd64"
)

func GetPlantform() string {
	return PLANTFORM
}

func GetLauncherName() string {
	return launcherName
}

func GetCQHttpName() string {
	return cqhttpName
}

func GetFastBuilderName() string {
	if fastBuilderName == "" {
		panic("未存在该预构建版本的FastBuilder: " + GetPlantform())
	}
	return fastBuilderName
}

func GetJDKDownloadName() string {
	if jdkDownloadName == "" {
		panic("未存在可用的JDK, 请尝试自行安装: " + GetPlantform())
	}
	return jdkDownloadName
}
