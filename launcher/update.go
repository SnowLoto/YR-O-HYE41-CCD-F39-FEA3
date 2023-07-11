package launcher

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"omega_launcher/plantform"
	"omega_launcher/utils"
	"os"
	"path/filepath"

	selfupdate "github.com/creativeprojects/go-selfupdate"
)

func CheckUpdate(currentVer string) {
	// Docker无需执行
	if utils.IsDocker() {
		return
	}
	// 获取原始文件名
	execName := plantform.GetOringinExecName()
	// 移除旧文件
	utils.RemoveFile(filepath.Join(utils.GetCurrentDir(), fmt.Sprintf(".%s.old", execName)))
	// 请求数据
	resp, err := http.Get("https://api.github.com/repos/Liliya233/omega_launcher/releases/latest")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	// 解析JSON
	var release struct {
		TagName string `json:"tag_name"`
	}
	err = json.Unmarshal(body, &release)
	if err != nil {
		return
	}
	// 自更新
	if utils.HasGreaterVer(currentVer, release.TagName) {
		exe, err := os.Executable()
		if err != nil {
			return
		}
		if err := selfupdate.UpdateTo(context.Background(), "https://github.com/Liliya233/omega_launcher/releases/latest/download/"+execName, execName, exe); err != nil {
			if err := selfupdate.UpdateTo(context.Background(), "https://www.omega-download.top/https://github.com/Liliya233/omega_launcher/releases/latest/download/"+execName, execName, exe); err != nil {
				return
			}
		}
	}
}
