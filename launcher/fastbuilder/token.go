package fastbuilder

import (
	"encoding/json"
	"omega_launcher/utils"
	"os"
	"path/filepath"
	"strings"

	"github.com/pterm/pterm"
)

// 判断是否为token
func isToken(token string) bool {
	if strings.HasPrefix(token, "w9/BeLNV/") || strings.HasPrefix(token, "w9/abqFz/") {
		return true
	}
	return false
}

// 加载现有的token
func loadCurrentFBToken() string {
	// 获取目录
	homedir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	fbconfigdir := filepath.Join(homedir, ".config/fastbuilder")
	token := filepath.Join(fbconfigdir, "fbtoken")
	// 尝试读取token文件
	if utils.IsFile(token) {
		if data, err := utils.GetFileData(token); err == nil {
			return string(data)
		}
	}
	return ""
}

// 请求token
func requestToken() string {
	// 尝试加载现有的token
	currentFbToken := loadCurrentFBToken()
	// 读取成功, 提示是否使用
	if currentFbToken != "" && isToken(currentFbToken) {
		if utils.GetInputYN("要使用当前设备储存的 Token 吗?") {
			return currentFbToken
		}
	}
	// 获取新的token
	Code := utils.GetValidInput("请输入 Fastbuilder 账号, 或者输入 Token")
	// 输入token则直接返回
	if isToken(Code) {
		pterm.Success.Println("输入内容为 FBToken")
		return Code
	}
	// 根据输入信息构建新token
	tokenstruct := &map[string]interface{}{
		"encrypt_token": true,
		"username":      Code,
		"password":      utils.GetPswInput("请输入 Fastbuilder 密码"),
	}
	token, err := json.Marshal(tokenstruct)
	if err != nil {
		panic(err)
	}
	return string(token)
}
