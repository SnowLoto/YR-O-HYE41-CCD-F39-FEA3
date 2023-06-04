package launcher

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

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
