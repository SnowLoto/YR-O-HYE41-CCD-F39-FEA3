package cqhttp

import (
	_ "embed"
	"encoding/json"
	"io/fs"
	"omega_launcher/defines"
	"omega_launcher/fastbuilder"
	"omega_launcher/utils"
	"path"
	"path/filepath"
	"strings"
)

//go:embed assets/组件-群服互通.json
var defaultQGroupLinkConfigByte []byte

//go:embed assets/组件-第三方__Liliya233__频服互通.json
var defaultQGuildLinkConfigByte []byte

func getOmegaQGroupLinkConfig() (string, *defines.QGroupLinkComponentConfig) {
	// 默认的空配置
	cfg := &defines.QGroupLinkComponentConfig{}
	// 默认配置文件路径
	fp := path.Join(fastbuilder.GetOmegaStorageDir(), "配置", "群服互通", "组件-群服互通.json")
	// 尝试从配置文件夹下寻找全部群服互通配置文件
	if err := filepath.Walk(path.Join(fastbuilder.GetOmegaStorageDir(), "配置"), func(filePath string, info fs.FileInfo, err error) error {
		// 跳过目录
		if info.IsDir() {
			return nil
		}
		// 识别非json组件文件并跳过
		fileName := info.Name()
		if !strings.HasPrefix(fileName, "组件") || !strings.HasSuffix(fileName, ".json") {
			return nil
		}
		// 对配置文件进行解析
		currentCfg := &defines.QGroupLinkComponentConfig{}
		if parseErr := utils.GetJsonData(filePath, currentCfg); parseErr != nil {
			return nil
		}
		// 如果不是群服互通组件, 则跳过
		if currentCfg.Name != "群服互通" {
			return nil
		}
		// 如果存在多个群服互通组件, 则报错
		if cfg.Configs != nil {
			panic("当前存在多个群服互通组件, 请自行删除多余的群服互通组件")
		}
		// 更新配置与路径信息
		cfg = currentCfg
		fp = filePath
		return nil
	}); err != nil {
		panic(err)
	}
	// 未找到配置时, 使用默认配置
	if cfg.Name != "群服互通" {
		err := json.Unmarshal(defaultQGroupLinkConfigByte, cfg)
		if err != nil {
			panic(err)
		}
	}
	return fp, cfg
}

func getOmegaQGuildLinkConfig() (string, *defines.QGuildLinkComponentConfig) {
	// 默认的空配置
	cfg := &defines.QGuildLinkComponentConfig{}
	// 默认配置文件路径
	fp := path.Join(fastbuilder.GetOmegaStorageDir(), "配置", "第三方", "Liliya233", "频服互通", "组件-第三方__Liliya233__频服互通.json")
	// 尝试从配置文件夹下寻找全部频服互通配置文件
	if err := filepath.Walk(path.Join(fastbuilder.GetOmegaStorageDir(), "配置"), func(filePath string, info fs.FileInfo, err error) error {
		// 跳过目录
		if info.IsDir() {
			return nil
		}
		// 识别非json组件文件并跳过
		fileName := info.Name()
		if !strings.HasPrefix(fileName, "组件") || !strings.HasSuffix(fileName, ".json") {
			return nil
		}
		// 对配置文件进行解析
		currentCfg := &defines.QGuildLinkComponentConfig{}
		if parseErr := utils.GetJsonData(filePath, currentCfg); parseErr != nil {
			return nil
		}
		// 如果不是频服互通组件, 则跳过
		if currentCfg.Name != "频服互通" {
			return nil
		}
		// 如果存在多个频服互通组件, 则报错
		if cfg.Configs != nil {
			panic("当前存在多个频服互通组件, 请自行删除多余的频服互通组件")
		}
		// 更新配置与路径信息
		cfg = currentCfg
		fp = filePath
		return nil
	}); err != nil {
		panic(err)
	}
	// 未找到配置时, 使用默认配置
	if cfg.Name != "频服互通" {
		err := json.Unmarshal(defaultQGuildLinkConfigByte, cfg)
		if err != nil {
			panic(err)
		}
	}
	return fp, cfg
}

// 将 Omega 配置内容写入到文件
func updateOmegaConfig(fp string, cfg any) {
	err := utils.WriteJsonData(fp, cfg)
	if err != nil {
		panic(err)
	}
}
