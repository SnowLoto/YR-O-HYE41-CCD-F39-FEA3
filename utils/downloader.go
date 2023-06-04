package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/pterm/pterm"
)

func DownloadSmallContent(sourceUrl string) []byte {
	// Get the data
	resp, err := http.Get(sourceUrl)
	if err != nil {
		pterm.Fatal.WithFatal(false).Println("从指定仓库下载资源时出现错误, 请重试或更换仓库")
		panic(err)
	}
	defer resp.Body.Close()
	// Size
	size, _ := strconv.Atoi(resp.Header.Get("Content-Length"))
	downloadSize := int64(size)

	// Progress Bar
	bar := pb.Full.Start64(downloadSize)
	bar.SetWidth(-1)
	bar.SetMaxWidth(100)
	bar.SetRefreshRate(time.Nanosecond)
	defer bar.Finish()

	// Reader
	barReader := bar.NewProxyReader(resp.Body)

	// Buffer
	contents := bytes.NewBuffer([]byte{})
	if _, err := io.Copy(contents, barReader); err == nil {
		return contents.Bytes()
	} else {
		panic(err)
	}
}

func GetLauncherUpdateInfo() string {
	// 请求数据
	resp, err := http.Get("https://api.github.com/repos/Liliya233/omega_launcher/releases/latest")
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	// 解析JSON
	var release struct {
		TagName string `json:"tag_name"`
	}
	err = json.Unmarshal(body, &release)
	if err != nil {
		return ""
	}

	// 返回版本号
	return strings.TrimPrefix(release.TagName, "v")
}
