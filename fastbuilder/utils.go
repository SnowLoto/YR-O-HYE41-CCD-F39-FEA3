package fastbuilder

import (
	"encoding/json"
	"omega_launcher/plantform"
	"omega_launcher/utils"
	"path/filepath"

	"github.com/pterm/pterm"
)

func GetOmegaStorageDir() string {
	return filepath.Join(utils.GetCurrentDir(), "omega_storage")
}

// 获取FB文件路径
func getFBExecPath() string {
	path := filepath.Join(utils.GetCurrentDir(), plantform.GetFastBuilderName())
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
	jsonData, err := utils.DownloadBytes(url + "hashes.json")
	if err != nil {
		panic(err)
	}
	// 解析文件内容
	var hash string
	hashMap := make(map[string]string, 0)
	if err := json.Unmarshal([]byte(jsonData), &hashMap); err != nil {
		panic(err)
	}
	hash = hashMap[plantform.GetFastBuilderName()]
	if hash == "" {
		pterm.Error.Printfln("未能从远程仓库获取 Hash")
	}
	return hash
}

// 检查当前目录是否存在FB执行文件
func CheckExecFile() bool {
	return utils.IsFile(getFBExecPath())
}
