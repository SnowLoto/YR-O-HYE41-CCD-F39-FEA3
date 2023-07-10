package cqhttp

import (
	_ "embed"
	"encoding/json"
	"io/fs"
	"omega_launcher/fastbuilder"
	"omega_launcher/utils"
	"path"
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

// Omega 群服互通配置文件结构
type QGroupLink struct {
	Address                   string                        `json:"CQHTTP正向Websocket代理地址"`
	GameMessageFormat         string                        `json:"游戏消息格式化模版"`
	QQMessageFormat           string                        `json:"Q群消息格式化模版"`
	Groups                    map[string]int64              `json:"链接的Q群"`
	Selector                  string                        `json:"游戏内可以听到QQ消息的玩家的选择器"`
	NoBotMsg                  bool                          `json:"不要转发机器人的消息"`
	ChatOnly                  bool                          `json:"只转发聊天消息"`
	MuteIgnored               bool                          `json:"屏蔽其他群的消息"`
	FilterQQToServerMsgByHead string                        `json:"仅仅转发开头为以下特定字符的消息到服务器"`
	FilterServerToQQMsgByHead string                        `json:"仅仅转发开头为以下特定字符的消息到QQ"`
	AllowedCmdExecutor        map[int64]bool                `json:"允许这些人透过QQ执行命令"`
	AllowdFakeCmdExecutor     map[int64]map[string][]string `json:"允许这些人透过QQ执行伪命令"`
	DenyCmds                  map[string]string             `json:"屏蔽这些指令"`
	AllowCmds                 []string                      `json:"允许所有人使用这些指令"`
	SendJoinAndLeaveMsg       bool                          `json:"向Q群发送玩家进出消息"`
	ShowExchangeDetail        bool                          `json:"在控制台显示消息转发详情"`
}

// Omega 频服互通配置文件结构
type QGuildLink struct {
	ChatOnly            bool                `json:"只转发聊天消息"`
	NoBotMsg            bool                `json:"不要转发机器人的消息"`
	SendJoinAndLeaveMsg bool                `json:"向频道发送玩家进出消息"`
	ShowExchangeDetail  bool                `json:"在控制台显示消息转发详情"`
	GameMessageFormat   string              `json:"游戏消息格式化模版"`
	QGuildMessageFormat string              `json:"频道消息格式化模版"`
	Address             string              `json:"CQHTTP正向Websocket代理地址"`
	Selector            string              `json:"游戏内可以听到QQ消息的玩家的选择器"`
	ServerToQQMsgFilter string              `json:"仅仅转发开头为以下特定字符的消息到QQ"`
	QQToServerMsgFilter string              `json:"仅仅转发开头为以下特定字符的消息到服务器"`
	DenyCmds            []string            `json:"不允许执行这些指令"`
	PublicCmds          []string            `json:"允许所有频道成员使用这些指令"`
	CmdExecutor         []string            `json:"允许这些身份组的频道成员透过QQ执行指令"`
	LinkChannelNames    map[string][]string `json:"链接的子频道"`
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
	utils.MkDir(path.Join(fastbuilder.GetOmegaStorageDir(), "配置"))
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
		currentCfg := &OmegaComponentConfig{}
		if parseErr := utils.GetJsonData(filePath, currentCfg); parseErr != nil {
			return nil
		}
		// 如果不是频服互通组件, 则跳过
		if currentCfg.Name != "群服互通" && currentCfg.Name != "第三方::Liliya233::频服互通" {
			return nil
		}
		// 更新并写入IP地址
		currentCfg.Configs["CQHTTP正向Websocket代理地址"] = address
		writeOmegaConfig(filePath, currentCfg)
		// 记录信息
		if currentCfg.Name == "群服互通" {
			hasQGroupCfg = true
			qGroupCfgFp = filePath
		} else {
			hasQGuildCfg = true
			qGuildCfgFp = filePath
		}
		return nil
	}); err != nil {
		panic(err)
	}
	// 未找到配置时, 使用默认配置
	if !hasQGroupCfg {
		utils.MkDir(path.Join(fastbuilder.GetOmegaStorageDir(), "配置", "群服互通"))
		qGroupCfgFp = path.Join(fastbuilder.GetOmegaStorageDir(), "配置", "群服互通", "组件-群服互通.json")
		newQGroupCfg := &OmegaComponentConfig{}
		if err := json.Unmarshal(defaultQGroupLinkConfigByte, newQGroupCfg); err == nil {
			newQGroupCfg.Configs["CQHTTP正向Websocket代理地址"] = address
			writeOmegaConfig(qGroupCfgFp, newQGroupCfg)
		}
	}
	if !hasQGuildCfg {
		utils.MkDir(path.Join(fastbuilder.GetOmegaStorageDir(), "配置", "第三方_by_Liliya233", "频服互通"))
		qGuildCfgFp = path.Join(fastbuilder.GetOmegaStorageDir(), "配置", "第三方_by_Liliya233", "频服互通", "组件-第三方__Liliya233__频服互通.json")
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
