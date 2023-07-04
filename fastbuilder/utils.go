package fastbuilder

import (
	"encoding/json"
	"omega_launcher/plantform"
	"omega_launcher/utils"
	"path"
	"path/filepath"

	"github.com/pterm/pterm"
)

// 获取FB文件名
func GetFBExecName() (name string) {
	switch plantform.GetPlantform() {
	case plantform.WINDOWS_arm64:
		// 不存在该构建
	case plantform.WINDOWS_x86_64:
		name = "phoenixbuilder-windows-executable-x86_64.exe"
	case plantform.Linux_arm64:
		name = "phoenixbuilder-aarch64"
	case plantform.Linux_x86_64:
		name = "phoenixbuilder"
	case plantform.MACOS_arm64:
		name = "phoenixbuilder-macos-arm64"
	case plantform.MACOS_x86_64:
		name = "phoenixbuilder-macos-x86_64"
	case plantform.Android_arm64:
		name = "phoenixbuilder-android-termux-shared-executable-arm64"
	case plantform.Android_x86_64:
		name = "phoenixbuilder-android-termux-shared-executable-x86_64"
	}
	if name == "" {
		panic("尚未支持该平台" + plantform.GetPlantform())
	}
	return name
}

// 获取FB文件路径
func getFBExecPath() string {
	path := path.Join(utils.GetCurrentDir(), GetFBExecName())
	result, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	return result
}

// 获取本地FB文件Hash
func getCurrentFBHash() string {
	exec := getFBExecPath()
	return utils.GetFileHash(exec)
}

// 获取远程仓库的Hash
func getRemoteFBHash(url string) string {
	// 获取文件内容
	jsonData := utils.DownloadSmallContent(url + "hashes.json")
	// 解析文件内容
	var hash string
	hashMap := make(map[string]string, 0)
	if err := json.Unmarshal([]byte(jsonData), &hashMap); err != nil {
		panic(err)
	}
	hash = hashMap[GetFBExecName()]
	if hash == "" {
		pterm.Error.Printfln("未能从远程仓库获取 Hash")
	}
	return hash
}

// 检查当前目录是否存在FB执行文件
func CheckExecFile() bool {
	return utils.IsFile(getFBExecPath())
}

func GetOmegaStorageDir() string {
	return path.Join(utils.GetCurrentDir(), "omega_storage")
}
