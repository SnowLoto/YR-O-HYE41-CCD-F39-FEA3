package cqhttp

import (
	_ "embed"
	"encoding/json"
	"io/fs"
	"omega_launcher/fastbuilder"
	"omega_launcher/utils"
	"path/filepath"
	"strings"
)

// Omega 配置文件结构
type OmegaComponentConfig struct {
	Name        string         `json:"名称"`
	Description string         `json:"描述"`
	Disabled    bool           `json:"是否禁用"`
	Version     string         `json:"版本"`
	Source      string         `json:"来源"`
	Configs     map[string]any `json:"配置"`
}

//go:embed assets/组件-群服互通.json
var defaultQGroupLinkConfigByte []byte

//go:embed assets/组件-第三方__Liliya233__频服互通.json
var defaultQGuildLinkConfigByte []byte

// 更新全部互通配置文件的地址
func updateOmegaConfigAddress(address string) (qGroupCfgFp, qGuildCfgFp string) {
	// 记录是否存在配置
	hasQGroupCfg, hasQGuildCfg := false, false
	// 尝试从配置文件夹下寻找全部频服互通配置文件
	utils.MkDir(filepath.Join(fastbuilder.GetOmegaStorageDir(), "配置"))
	if err := filepath.Walk(filepath.Join(fastbuilder.GetOmegaStorageDir(), "配置"), func(path string, info fs.FileInfo, err error) error {
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
		currentCfg := &OmegaComponentConfig{}
		if parseErr := utils.GetJsonData(path, currentCfg); parseErr != nil {
			return nil
		}
		// 如果不是频服互通组件, 则跳过
		if currentCfg.Name != "群服互通" && currentCfg.Name != "第三方::Liliya233::频服互通" {
			return nil
		}
		// 更新并写入IP地址
		currentCfg.Configs["CQHTTP正向Websocket代理地址"] = address
		writeOmegaConfig(path, currentCfg)
		// 记录信息
		if currentCfg.Name == "群服互通" {
			hasQGroupCfg = true
			qGroupCfgFp = path
		} else {
			hasQGuildCfg = true
			qGuildCfgFp = path
		}
		return nil
	}); err != nil {
		panic(err)
	}
	// 未找到配置时, 使用默认配置
	if !hasQGroupCfg {
		utils.MkDir(filepath.Join(fastbuilder.GetOmegaStorageDir(), "配置", "群服互通"))
		qGroupCfgFp = filepath.Join(fastbuilder.GetOmegaStorageDir(), "配置", "群服互通", "组件-群服互通.json")
		newQGroupCfg := &OmegaComponentConfig{}
		if err := json.Unmarshal(defaultQGroupLinkConfigByte, newQGroupCfg); err == nil {
			newQGroupCfg.Configs["CQHTTP正向Websocket代理地址"] = address
			writeOmegaConfig(qGroupCfgFp, newQGroupCfg)
		}
	}
	if !hasQGuildCfg {
		utils.MkDir(filepath.Join(fastbuilder.GetOmegaStorageDir(), "配置", "第三方_by_Liliya233", "频服互通"))
		qGuildCfgFp = filepath.Join(fastbuilder.GetOmegaStorageDir(), "配置", "第三方_by_Liliya233", "频服互通", "组件-第三方__Liliya233__频服互通.json")
		newQGuildCfg := &OmegaComponentConfig{}
		if err := json.Unmarshal(defaultQGuildLinkConfigByte, newQGuildCfg); err == nil {
			newQGuildCfg.Configs["CQHTTP正向Websocket代理地址"] = address
			writeOmegaConfig(qGuildCfgFp, newQGuildCfg)
		}
	}
	// 返回路径
	return qGroupCfgFp, qGuildCfgFp
}

// 将 Omega 配置内容写入到文件
func writeOmegaConfig(fp string, cfg any) {
	err := utils.WriteJsonData(fp, cfg)
	if err != nil {
		panic(err)
	}
}
