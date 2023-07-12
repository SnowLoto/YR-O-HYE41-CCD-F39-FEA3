package fastbuilder

import (
	"encoding/json"
	"omega_launcher/launcher"
	"omega_launcher/utils"
	"os"
	"path/filepath"
	"strings"

	"github.com/pterm/pterm"
)

// 判断是否为token
func IsToken(token string) bool {
	return strings.HasPrefix(token, "w9/")
}

// 加载现有的token
func loadCurrentFBToken() string {
	// 获取目录
	homedir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	tokenPath := filepath.Join(homedir, ".config", "fastbuilder", "fbtoken")
	// 尝试读取token文件
	if utils.IsFile(tokenPath) {
		if data, err := utils.GetFileData(tokenPath); err == nil {
			return string(data)
		}
	}
	return ""
}

// 删除现有的token
func deleteCurrentFBToken() bool {
	// 获取目录
	homedir, err := os.UserHomeDir()
	if err != nil {
		return false
	}
	tokenPath := filepath.Join(homedir, ".config", "fastbuilder", "fbtoken")
	// 尝试删除token文件
	return utils.RemoveFile(tokenPath)
}

// 请求token
func requestToken() string {
	// 尝试加载现有的token
	currentFbToken := loadCurrentFBToken()
	// 读取成功, 提示是否使用
	if currentFbToken != "" && IsToken(currentFbToken) {
		if utils.GetInputYN("要使用当前设备储存的 Token 吗?") {
			return currentFbToken
		}
		deleteCurrentFBToken()
	}
	// 获取新的token
	Code := utils.GetValidInput("请输入 Fastbuilder 账号, 或者输入 Token")
	// 输入token则直接返回
	if IsToken(Code) {
		pterm.Success.Println("输入内容为 Token")
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

// 配置Token
func FBTokenSetup(cfg *launcher.Config) {
	if cfg.FBToken != "" {
		if utils.GetInputYN("要使用上次的 Fastbuilder 账号登录吗?") {
			return
		}
	}
	cfg.FBToken = requestToken()
}
