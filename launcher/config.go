package launcher

import (
	"omega_launcher/defines"
	"omega_launcher/utils"
	"path"
	"strings"
)

// 配置租赁服信息
func RentalServerSetup(cfg *defines.LauncherConfig) {
	cfg.RentalCode = utils.GetValidInput("请输入租赁服号")
	cfg.RentalPasswd = utils.GetPswInput("请输入租赁服密码")
}

// 保存配置文件
func SaveConfig(config *defines.LauncherConfig) {
	copyConfig := *config
	if strings.HasPrefix(copyConfig.FBToken, "{\"encrypt_token\"") {
		copyConfig.FBToken = ""
	}
	utils.WriteJsonData(path.Join(utils.GetCurrentDataDir(), "服务器登录配置.json"), copyConfig)
}
