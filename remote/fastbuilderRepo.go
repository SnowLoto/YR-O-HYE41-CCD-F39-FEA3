package remote

import (
	"encoding/json"
	"omega_launcher/utils"

	"github.com/pterm/pterm"
)

var (
	fastbuilderRepoRemoteData    []*FastbuilderRepoRemoteData
	fastbuilderRepoRemoteDataUrl = "https://github.com/Liliya233/omega_launcher/blob/main/remote_data/config/fastbuilder_repo.json"
)

type FastbuilderRepoRemoteData struct {
	Name         string `json:"Name"`
	Url          string `json:"Url"`
	IsGithub     bool   `json:"IsGithub"`
	IsPreRelease bool   `json:"IsPreRelease"`
}

func GetFastbuilderRepoRemoteData() []*FastbuilderRepoRemoteData {
	if fastbuilderRepoRemoteData == nil {
		// 下载远程数据
		pterm.Warning.Println("正在获取 FastBuilder 仓库列表..")
		bytes, err := utils.DownloadBytesWithMirror(fastbuilderRepoRemoteDataUrl)
		if err != nil {
			panic(err)
		}
		// 解析远程数据
		newFastbuilderRepoRemoteData := []*FastbuilderRepoRemoteData{}
		if err := json.Unmarshal(bytes, &newFastbuilderRepoRemoteData); err != nil {
			panic(err)
		}
		fastbuilderRepoRemoteData = newFastbuilderRepoRemoteData
	}
	// 返回结构体指针
	return fastbuilderRepoRemoteData
}
