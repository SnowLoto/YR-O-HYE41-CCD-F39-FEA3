package remote

import (
	"encoding/json"
	"omega_launcher/utils"

	"github.com/pterm/pterm"
)

var (
	signServerRemoteData    *SignServerRemoteData
	signServerRemoteDataUrl = "https://github.com/Liliya233/omega_launcher/blob/main/remote_data/config/sign_server.json"
)

type SignServerRemoteData struct {
	DownloadUrl  string `json:"DownloadUrl"`
	ZipName      string `json:"ZipName"`
	ZipSize      int64  `json:"ZipSize"`
	UnzipDirName string `json:"UnzipDirName"`
	SoVersion    string `json:"SoVersion"`
}

func GetSignServerRemoteData() *SignServerRemoteData {
	if signServerRemoteData == nil {
		// 下载远程数据
		pterm.Warning.Println("正在获取 Sign Server 配置所需的数据..")
		bytes := utils.DownloadBytesWithMirror(signServerRemoteDataUrl)
		// 解析远程数据
		newSignServerRemoteData := SignServerRemoteData{}
		if err := json.Unmarshal(bytes, &newSignServerRemoteData); err != nil {
			panic(err)
		}
		signServerRemoteData = &newSignServerRemoteData
	}
	// 返回结构体指针
	return signServerRemoteData
}
