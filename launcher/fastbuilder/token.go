package fastbuilder

import (
	"encoding/json"
	"fmt"
	"omega_launcher/utils"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/pterm/pterm"
	"golang.org/x/term"
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
		pterm.Info.Printf("要使用当前设备储存的 Token 吗? 要请输入 y, 不要请输入 n: ")
		if utils.GetInputYN() {
			return currentFbToken
		}
	}
	// 获取新的token
	pterm.Info.Printf("请输入 Fastbuilder 账号, 或者输入 Token: ")
	Code := utils.GetValidInput()
	// 输入token则直接返回
	if isToken(Code) {
		pterm.Success.Println("输入内容为 FBToken")
		return Code
	}
	pterm.Info.Printf("请输入 Fastbuilder 密码 (不会回显): ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Print("\n")
	if err != nil {
		panic(err)
	}
	Passwd := string(bytePassword)
	// 根据输入信息构建新token
	tokenstruct := &map[string]interface{}{
		"encrypt_token": true,
		"username":      Code,
		"password":      Passwd,
	}
	token, err := json.Marshal(tokenstruct)
	if err != nil {
		panic(err)
	}
	return string(token)
}
