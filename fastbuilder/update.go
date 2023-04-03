package fastbuilder

import (
	"bytes"
	"io"
	"omega_launcher/defines"
	"omega_launcher/utils"

	"github.com/pterm/pterm"
)

// Fastbuilder远程仓库地址
var STORAGE_REPO = ""

// 仓库选择
func selectRepo(cfg *defines.LauncherConfig, reselect bool) {
	if reselect || cfg.Repo < 1 || cfg.Repo > 6 {
		// 不再于列表提示自用仓库
		utils.ConfPrinter.Println(
			"当前可选择的仓库有：\n",
			"1. 官方仓库\n",
			"2. 官方镜像仓库\n",
			"3. 云裳公益镜像仓库\n",
			"4. rnhws-Team 仓库\n",
			"5. rnhws-Team 镜像仓库",
		)
		cfg.Repo = utils.GetIntInputInScope("请输入序号来选择一个仓库", 1, 6)
	}
	switch cfg.Repo {
	case 1:
		pterm.Info.Printfln("将使用官方仓库 (%s) 进行更新", defines.OFFICIAL_REPO)
		STORAGE_REPO = defines.OFFICIAL_REPO
	case 2:
		pterm.Info.Printfln("将使用官方镜像仓库 (%s) 进行更新", defines.OFFICIAL_REPO_MIRROR)
		STORAGE_REPO = defines.OFFICIAL_REPO_MIRROR
	case 3:
		pterm.Info.Printfln("将使用云裳公益镜像仓库 (%s) 进行更新", defines.YSCLOUD_REPO)
		STORAGE_REPO = defines.YSCLOUD_REPO
	case 4:
		pterm.Info.Printfln("将使用 rnhws-Team 仓库 (%s) 进行更新", defines.RNHWS_TEAM_REPO)
		STORAGE_REPO = defines.RNHWS_TEAM_REPO
	case 5:
		pterm.Info.Printfln("将使用 rnhws-Team 镜像仓库 (%s) 进行更新", defines.RNHWS_TEAM_REPO_MIRROR)
		STORAGE_REPO = defines.RNHWS_TEAM_REPO_MIRROR
	case 6:
		pterm.Info.Printfln("将使用 本地仓库 (%s) 进行更新", defines.LOCAL_REPO)
		STORAGE_REPO = defines.LOCAL_REPO
	default:
		panic("无效的仓库, 请重新配置")
	}
}

// 下载FB
func download() {
	var execBytes []byte
	var err error
	// 获取写入路径与远程仓库url
	path := getFBExecPath()
	url := STORAGE_REPO + GetFBExecName()
	// 下载
	compressedData := utils.DownloadSmallContent(url)
	// 官网并没有提供brotli, 所以对读取操作进行修改
	if execBytes, err = io.ReadAll(bytes.NewReader(compressedData)); err != nil {
		panic(err)
	}
	// 写入文件
	if err := utils.WriteFileData(path, execBytes); err != nil {
		panic(err)
	}
}

// 升级FB
func Update(cfg *defines.LauncherConfig, reselect bool) {
	selectRepo(cfg, reselect)
	pterm.Warning.Println("正在从指定仓库获取更新信息..")
	if getRemoteFBHash(STORAGE_REPO) == getCurrentFBHash() {
		pterm.Success.Println("太好了, 你的 Fastbuilder 已经是最新的了!")
	} else {
		pterm.Warning.Println("正在为你下载最新的 Fastbuilder, 请保持耐心..")
		download()
	}
}
