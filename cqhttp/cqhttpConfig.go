package cqhttp

import (
	"fmt"
	"omega_launcher/utils"
	"os"
	"path"
	"strings"

	_ "embed"

	"github.com/pterm/pterm"
	v2 "gopkg.in/yaml.v2"
)

//go:embed assets/config.yml
var defaultConfigBytes []byte

// Copy from go-cqhttp
// Account 账号配置
type Account struct {
	Uin        int64  `yaml:"uin"`
	Password   string `yaml:"password"`
	SignServer string `yaml:"sign-server"`
}

// Config 总配置文件
type CQHttpConfig struct {
	Account *Account `yaml:"account"`
}

// 从cqhttp配置里读取QQ账密信息
func getCQConfig() *CQHttpConfig {
	cfg := &CQHttpConfig{}
	data, err := os.ReadFile(path.Join(GetCQHttpDir(), "config.yml"))
	if err != nil {
		return nil
	}
	if err := v2.Unmarshal(data, cfg); err != nil {
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
	cfgStr := strings.Replace(string(defaultConfigBytes), "[WS地址]", address, 1)
	// 保留账密信息
	if cqCfg != nil {
		cfgStr = strings.Replace(cfgStr, "[QQ账号]", fmt.Sprint(cqCfg.Account.Uin), 1)
		cfgStr = strings.Replace(cfgStr, "[QQ密码]", cqCfg.Account.Password, 1)
		cfgStr = strings.Replace(cfgStr, "[SignServer地址]", cqCfg.Account.SignServer, 1)
	} else {
		cfgStr = strings.Replace(cfgStr, "[QQ账号]", "null", 1)
		cfgStr = strings.Replace(cfgStr, "[QQ密码]", "", 1)
		cfgStr = strings.Replace(cfgStr, "[SignServer地址]", "-", 1)
	}
	// 写入新配置
	writeCQConfig(cfgStr)
}

// 初始化cqhttp
func initCQConfig() {
	// 移除token等文件
	utils.RemoveFile(path.Join(GetCQHttpDir(), "device.json"))
	utils.RemoveFile(path.Join(GetCQHttpDir(), "session.token"))
	// 要求输入cqhttp配置信息
	cfgStr := strings.Replace(string(defaultConfigBytes), "[QQ账号]", fmt.Sprint(utils.GetInt64Input("请输入QQ账号")), 1)
	cfgStr = strings.Replace(cfgStr, "[QQ密码]", utils.GetPswInput("请输入QQ密码"), 1)
	SignServer := utils.GetInput("请输入SignServer地址, 没有可留空")
	if SignServer == "" {
		SignServer = "-"
	}
	cfgStr = strings.Replace(cfgStr, "[SignServer地址]", SignServer, 1)
	cfgStr = strings.Replace(cfgStr, "[WS地址]", "null", 1)
	// 写入新配置
	writeCQConfig(cfgStr)
}
