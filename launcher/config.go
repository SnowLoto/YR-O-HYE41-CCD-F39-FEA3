package launcher

import (
	"omega_launcher/utils"
	"path/filepath"
)

// 启动器配置文件结构
type Config struct {
	Repo                int    `json:"仓库序号"`
	RentalCode          string `json:"租赁服号"`
	RentalPasswd        string `json:"租赁服密码"`
	FBToken             string `json:"FBToken"`
	EnableCQHttp        bool   `json:"是否开启go-cqhttp"`
	StartOmega          bool   `json:"是否启动Omega"`
	UpdateFB            bool   `json:"是否更新FB"`
	BlockCQHttpOutput   bool   `json:"是否屏蔽go-cqhttp的内容输出"`
	EnableSignServer    bool   `json:"是否使用SignServer"`
	SignServerSoVersion string `json:"SignServer使用的so文件版本"`
}

// 保存配置文件
func SaveConfig(config *Config) {
	utils.WriteJsonData(filepath.Join(utils.GetCurrentDataDir(), "服务器登录配置.json"), config)
}
