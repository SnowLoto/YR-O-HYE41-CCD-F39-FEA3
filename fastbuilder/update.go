package fastbuilder

import (
	"encoding/json"
	"fmt"
	"omega_launcher/launcher"
	"omega_launcher/plantform"
	"omega_launcher/remote"
	"omega_launcher/utils"
	"strings"
	"time"

	"github.com/pterm/pterm"
)

type release struct {
	TagName         string `json:"tag_name"`
	TargetCommitish string `json:"target_commitish"`
	PublishedAt     string `json:"published_at"`
}

func getLatestLibreReleaseVersion() string {
	pterm.Warning.Printfln("正在获取最新的 Libre-release 版本号..")
	var releases []release
	bytes, err := utils.DownloadBytes("https://api.github.com/repos/LNSSPsd/PhoenixBuilder/releases")
	if err != nil {
		pterm.Error.Println("获取最新的 Libre-release 版本号时出现问题")
		panic(err)
	}
	if err := json.Unmarshal(bytes, &releases); err != nil {
		pterm.Error.Println("解析最新的 Libre-release 版本号时出现问题")
		panic(err)
	}
	var latestPreRelease release
	for _, release := range releases {
		// 不是 libre 分支则直接跳过
		if !strings.Contains(release.TargetCommitish, "libre") {
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

func UpdateRepo(cfg *launcher.Config) {
	// 获取远程数据
	remoteDatas := remote.GetFastbuilderRepoRemoteData()
	// 仓库切片
	repos := []launcher.ConfigRepo{}
	// 额外提示信息样式
	mirrorTip := pterm.Green("[镜像]")
	// 提示信息切片
	tips := []string{"当前可选择的仓库有："}
	for _, data := range remoteDatas {
		// 基本链接
		repos = append(repos,
			launcher.ConfigRepo{
				Name:         data.Name,
				Url:          data.Url,
				IsPreRelease: data.IsPreRelease,
			},
		)
		tips = append(tips, fmt.Sprintf("%v. %v", len(tips), data.Name))
		// 是否 Github
		if data.IsGithub {
			repos = append(repos,
				launcher.ConfigRepo{
					Name:         fmt.Sprintf("%v %v", data.Name, mirrorTip),
					Url:          data.Url,
					IsPreRelease: data.IsPreRelease,
					UseMirror:    true,
				},
			)
			tips = append(tips, fmt.Sprintf("%v. %v %v", len(tips), data.Name, mirrorTip))
		}
	}
	// 打印仓库列表
	utils.ConfPrinter.Println(strings.Join(tips, "\n"))
	cfg.Repo = &repos[utils.GetIntInputInScope("请输入序号来选择一个仓库", 1, len(tips))-1]
}

// 下载FB
func download(useMirror bool, url string) {
	// 镜像下载
	if useMirror && utils.DownloadFileWithMirror(url+plantform.GetFastBuilderName(), getFBExecPath()) == nil {
		return
	}
	// 普通下载
	err := utils.DownloadFile(url+plantform.GetFastBuilderName(), getFBExecPath())
	if err != nil {
		panic(err)
	}
}

// 升级FB
func Update(cfg *launcher.Config) {
	// 检查是否为空的url
	if cfg.Repo == nil || cfg.Repo.Url == "" {
		pterm.Warning.Println("FastBuilder 更新链接为空, 请重新选择仓库来启用更新功能")
		cfg.UpdateFB = false
		return
	}
	pterm.Warning.Println(fmt.Sprintf("正在从 %v 获取更新信息..", cfg.Repo.Name))
	// 更新使用的url
	url := cfg.Repo.Url
	// 如果是预构建, 需要获取tag来拼接为完整的url
	if cfg.Repo.IsPreRelease {
		url = url + getLatestLibreReleaseVersion() + "/"
	}
	if getRemoteFBHash(url) == getCurrentFBHash() {
		pterm.Success.Println("太好了, 你的 FastBuilder 已经是最新的了!")
	} else {
		pterm.Warning.Println("正在为你下载最新的 FastBuilder, 请保持耐心..")
		download(cfg.Repo.UseMirror, url)
	}
}
