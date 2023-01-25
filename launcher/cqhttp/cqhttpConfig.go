package cqhttp

import (
	"fmt"
	"omega_launcher/defines"
	"omega_launcher/utils"
	"os"
	"path"
	"strings"
	"syscall"

	_ "embed"

	"github.com/pterm/pterm"
	"golang.org/x/term"
	v2 "gopkg.in/yaml.v2"
)

//go:embed assets/config.yml
var defaultConfigBytes []byte

// 从cqhttp配置里读取QQ账密信息
func getCQConfig() *defines.CQHttpConfig {
	data, err := os.ReadFile(path.Join(GetCQHttpDir(), "config.yml"))
	if err != nil {
		return nil
	}
	cfg := &defines.CQHttpConfig{}
	if err := v2.Unmarshal(data, &cfg); err != nil {
		return nil
	}
	if cfg.Account.Uin == "" {
		return nil
	}
	return cfg
}

// 写入cqhttp配置
func writeCQConfig(cfgStr string) {
	err := utils.WriteFileData(path.Join(GetCQHttpDir(), "config.yml"), []byte(cfgStr))
	if err != nil {
		pterm.Fatal.WithFatal(false).Println("更新 go-cqhttp 配置文件时遇到问题")
		panic(err)
	}
}

// 更新cqhttp配置文件的地址
func updateCQConfigAddress(address string) {
	cqCfg := getCQConfig()
	cfgStr := strings.Replace(string(defaultConfigBytes), "[地址]", address, 1)
	// 保留账密信息
	if cqCfg != nil {
		cfgStr = strings.Replace(cfgStr, "[QQ账号]", cqCfg.Account.Uin, 1)
		cfgStr = strings.Replace(cfgStr, "[QQ密码]", cqCfg.Account.Password, 1)
	} else {
		cfgStr = strings.Replace(cfgStr, "[QQ账号]", "null", 1)
		cfgStr = strings.Replace(cfgStr, "[QQ密码]", "null", 1)
	}
	// 写入新配置
	writeCQConfig(cfgStr)
}

// 初始化cqhttp配置文件
func initCQConfig() {
	if utils.IsDocker() {
		panic("非本地环境只能通过上传文件的方式来配置 go-cqhttp")
	}
	// 要求输入cqhttp配置信息
	pterm.Info.Printf("请输入QQ账号: ")
	cfgStr := strings.Replace(string(defaultConfigBytes), "[QQ账号]", utils.GetValidInput(), 1)
	pterm.Info.Printf("请输入QQ密码 (不会回显): ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Print("\n")
	if err != nil {
		panic(err)
	}
	cfgStr = strings.Replace(cfgStr, "[QQ密码]", string(bytePassword), 1)
	cfgStr = strings.Replace(cfgStr, "[地址]", "null", 1)
	// 写入新配置
	writeCQConfig(cfgStr)
}
