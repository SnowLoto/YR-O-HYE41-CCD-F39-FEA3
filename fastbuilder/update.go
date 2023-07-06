package fastbuilder

import (
	"bytes"
	"encoding/json"
	"io"
	"omega_launcher/defines"
	"omega_launcher/utils"
	"time"

	"github.com/pterm/pterm"
)

var (
	// Fastbuilder远程仓库地址
	STORAGE_REPO = ""

	// Github 镜像地址
	GITHUB_MIRROR = "https://www.omega-download.top/"

	// Github 官方仓库
	OFFICIAL_REPO = "https://github.com/LNSSPsd/PhoenixBuilder/releases/latest/download/"

	// 海的官方镜像仓库
	SEA_REPO = "https://likemc.xyz/build/"

	// Github rnhws-Team 预览版仓库
	RNHWS_TEAM_REPO = "https://github.com/rnhws-Team/PhoenixBuilder/releases/latest/download/"

	// Github 官方 Libre 分支仓库
	OFFICIAL_LIBRE_REPO = "https://github.com/LNSSPsd/PhoenixBuilder/releases/download/"

	// 本地 (自用)
	LOCAL_REPO = "http://fileserver:12333/res/"
)

type release struct {
	TagName         string `json:"tag_name"`
	TargetCommitish string `json:"target_commitish"`
	PublishedAt     string `json:"published_at"`
}

func getLatestLibreReleaseVersion() string {
	pterm.Warning.Printfln("正在获取最新的 Libre-release 版本号..")
	var releases []release
	err := json.Unmarshal(utils.DownloadSmallContent("https://api.github.com/repos/LNSSPsd/PhoenixBuilder/releases"), &releases)
	if err != nil {
		pterm.Error.Println("正在获取最新的 Libre-release 版本号时出现问题")
		panic(err)
	}
	var latestPreRelease release
	for _, release := range releases {
		// 不是 libre 分支则直接跳过
		if release.TargetCommitish != "libre" {
			continue
		}
		// 解析 latest 的时间, 出错时代表 latestPreRelease 为空, 可直接替换
		latestTime, err := time.Parse(time.RFC3339, latestPreRelease.PublishedAt)
		if err != nil {
			latestPreRelease = release
			continue
		}
		time, err := time.Parse(time.RFC3339, release.PublishedAt)
		if err != nil {
			panic("解析 release 信息时出现错误")
		}
		// 对比两个 release 的发布时间, 且替换为最新
		if time.After(latestTime) {
			latestPreRelease = release
		}
	}
	if latestPreRelease.TagName == "" {
		panic("未获取到任何 Libre-release 信息")
	}
	return latestPreRelease.TagName
}

// 仓库选择
func selectRepo(cfg *defines.LauncherConfig, reselect bool) {
	if reselect || cfg.Repo < 1 || cfg.Repo > 6 {
		// 不再于列表提示自用仓库
		utils.ConfPrinter.Println(
			"当前可选择的仓库有：\n",
			"1. 官方仓库"+pterm.Yellow(" (推荐)\n"),
			"2. 官方镜像仓库\n",
			"3. 海的官方镜像仓库\n",
			"4. rnhws-Team 预览版仓库\n",
			"5. rnhws-Team 预览版镜像仓库\n",
			"6. 官方 Libre 分支仓库\n",
			"7. 官方 Libre 分支镜像仓库",
		)
		cfg.Repo = utils.GetIntInputInScope("请输入序号来选择一个仓库", 1, 8)
	}
	switch cfg.Repo {
	case 1:
		STORAGE_REPO = OFFICIAL_REPO
		pterm.Info.Printfln("将使用 官方仓库 (%s) 进行更新", STORAGE_REPO)
	case 2:
		STORAGE_REPO = GITHUB_MIRROR + OFFICIAL_REPO
		pterm.Info.Printfln("将使用 官方镜像仓库 (%s) 进行更新", STORAGE_REPO)
	case 3:
		STORAGE_REPO = SEA_REPO
		pterm.Info.Printfln("将使用 海的官方镜像仓库 (%s) 进行更新", STORAGE_REPO)
	case 4:
		STORAGE_REPO = RNHWS_TEAM_REPO
		pterm.Info.Printfln("将使用 rnhws-Team 预览版仓库 (%s) 进行更新", STORAGE_REPO)
	case 5:
		STORAGE_REPO = GITHUB_MIRROR + RNHWS_TEAM_REPO
		pterm.Info.Printfln("将使用 rnhws-Team 预览版镜像仓库 (%s) 进行更新", STORAGE_REPO)
	case 6:
		STORAGE_REPO = OFFICIAL_LIBRE_REPO + getLatestLibreReleaseVersion() + "/"
		pterm.Info.Printfln("将使用 官方 Libre 分支仓库 (%s) 进行更新", STORAGE_REPO)
	case 7:
		STORAGE_REPO = GITHUB_MIRROR + OFFICIAL_LIBRE_REPO + getLatestLibreReleaseVersion() + "/"
		pterm.Info.Printfln("将使用 官方 Libre 分支镜像仓库 (%s) 进行更新", STORAGE_REPO)
	case 8:
		STORAGE_REPO = LOCAL_REPO
		pterm.Info.Printfln("将使用 本地仓库 (%s) 进行更新", STORAGE_REPO)
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
