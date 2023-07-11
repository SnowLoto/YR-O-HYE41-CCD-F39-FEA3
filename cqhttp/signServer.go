package cqhttp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"omega_launcher/deploy"
	"omega_launcher/plantform"
	"omega_launcher/utils"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/pterm/pterm"
)

var (
	signServerDownloadUrl = "https://www.omega-download.top/https://github.com/fuqiuluo/unidbg-fetch-qsign/releases/latest/download/"
)

// config.json of sign server
type Config struct {
	Server struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	} `json:"server"`
	Key            string `json:"key"`
	ReloadInterval int    `json:"reload_interval"`
	Protocol       struct {
		Qua     string `json:"qua"`
		Version string `json:"version"`
		Code    string `json:"code"`
	} `json:"protocol"`
	Unidbg struct {
		Dynarmic bool `json:"dynarmic"`
		Unicorn  bool `json:"unicorn"`
		Debug    bool `json:"debug"`
	} `json:"unidbg"`
}

func checkHTTPConnection(url string) bool {
	client := http.Client{
		Timeout: time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

func isSignServerCache() (string, string, bool) {
	pterm.Warning.Printfln("正在检查 Sign Server 更新..")
	var release struct {
		TagName string `json:"tag_name"`
		Assets  []struct {
			Name string `json:"name"`
			Size int64  `json:"size"`
		} `json:"assets"`
	}
	err := json.Unmarshal(utils.DownloadBytes("https://api.github.com/repos/fuqiuluo/unidbg-fetch-qsign/releases/latest"), &release)
	if err != nil {
		pterm.Error.Println("检查 Sign Server 更新时出现问题")
		panic(err)
	}
	downloadFileName := fmt.Sprintf("unidbg-fetch-qsign-%s.zip", release.TagName)
	downloadFilePath := filepath.Join(utils.GetCacheDir(), "downloads", downloadFileName)
	// 找到asset的Name与downloadFileName相等的一项, 然后对比它的size
	for _, asset := range release.Assets {
		if asset.Name == downloadFileName {
			currentSize, err := utils.GetFileSize(downloadFilePath)
			if err != nil {
				return downloadFileName, downloadFilePath, false
			}
			if currentSize == asset.Size {
				return downloadFileName, downloadFilePath, true
			}
			break
		}
	}
	return downloadFileName, downloadFilePath, false
}

func signServerDeploy() string {
	downloadFileName, downloadFilePath, isCache := isSignServerCache()
	if !isCache {
		pterm.Warning.Printfln("正在下载 Sign Server 可执行文件..")
		utils.DownloadFile(signServerDownloadUrl+downloadFileName, downloadFilePath)
	}
	execDirName := strings.TrimSuffix(downloadFileName, ".zip")
	if !utils.IsDir(filepath.Join(utils.GetCacheDir(), "SignServer", execDirName)) {
		pterm.Warning.Printfln("正在解压 Sign Server 可执行文件..")
		fp, err := os.OpenFile(downloadFilePath, os.O_RDONLY, 0755)
		if err != nil {
			panic(err)
		}
		if err := utils.UnZip(fp, filepath.Join(utils.GetCacheDir(), "SignServer")); err != nil {
			utils.RemoveFile(downloadFilePath)
			panic(err)
		}
		pterm.Success.Printfln("Sign Server 已成功部署")
	}
	return execDirName
}

func setupConfig(configPath string, host string, port int, uin int64, androidID string) {
	// 读取配置文件
	configBytes, err := utils.GetFileData(configPath)
	if err != nil {
		panic("读取 Sign Server 配置时出现错误: " + err.Error())
	}
	// 解析配置文件
	var config Config
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		panic("解析 Sign Server 配置时出现错误: " + err.Error())
	}
	// 修改配置
	config.Key = "233233"
	config.Server.Host = host
	config.Server.Port = port
	modifiedConfigBytes, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		panic("序列化配置文件时出现错误: " + err.Error())
	}
	err = utils.WriteFileData(configPath, modifiedConfigBytes)
	if err != nil {
		panic("写入配置文件时出现错误: " + err.Error())
	}
}

func SignServerStart(soVersion string) string {
	deviceJsonPath := filepath.Join(GetCQHttpDir(), "device.json")
	if !utils.IsFile(deviceJsonPath) {
		pterm.Warning.Println("尚未生成 device.json, 将跳过 Sign Server 的配置")
		return ""
	}
	deviceJsonContent, _ := utils.GetFileData(deviceJsonPath)
	var device struct {
		Protocol  int    `json:"protocol"`
		AndroidId string `json:"android_id"`
	}
	err := json.Unmarshal(deviceJsonContent, &device)
	if err != nil {
		pterm.Error.Println("读取 device.json 时出现错误, 将跳过 Sign Server 的配置: " + err.Error())
		return ""
	}
	if device.Protocol != 6 {
		pterm.Error.Println("未使用 Android Pad 登录协议 (6), 将跳过 Sign Server 的配置")
		return ""
	}
	signServerDirName := signServerDeploy()
	availablePort, err := utils.GetAvailablePort()
	if err != nil {
		pterm.Error.Println("无法获取可用的端口来启动 Sign Server: " + err.Error())
		return ""
	}
	if cqCfg := getCQConfig(); cqCfg != nil {
		configPath := filepath.Join(utils.GetCacheDir(), "SignServer", signServerDirName, "txlib", soVersion, "config.json")
		setupConfig(configPath, "0.0.0.0", availablePort, cqCfg.Account.Uin, device.AndroidId)
	} else {
		pterm.Error.Println("Sign Server 启动失败, 未能够从 go-cqhttp 配置文件中获取QQ账号")
		return ""
	}
	args := []string{
		fmt.Sprintf("--basePath=%s", filepath.Join(utils.GetCacheDir(), "SignServer", signServerDirName, "txlib", soVersion)),
	}
	// 如果不是Windows则去掉.bat
	cmdStr := filepath.Join(utils.GetCacheDir(), "SignServer", signServerDirName, "bin", "unidbg-fetch-qsign.bat")
	if currentPlanform := plantform.GetPlantform(); currentPlanform != plantform.WINDOWS_x86_64 && currentPlanform != plantform.WINDOWS_arm64 {
		cmdStr = strings.TrimSuffix(cmdStr, ".bat")
		os.Chmod(cmdStr, 0755)
	}
	// 运行Sign Server
	cmd := exec.Command(cmdStr, args...)
	if !deploy.CheckJava() {
		pterm.Warning.Println("系统未安装Java, 将尝试安装对应的 Java 环境")
		deploy.JavaDeploy()
		cmd.Env = append(os.Environ(), fmt.Sprintf("JAVA_HOME=%s", filepath.Join(utils.GetCacheDir(), "Java", "jdk-20.0.1")))
	}
	go func() {
		pterm.Success.Println("正在启动 Sign Server..")
		for {
			if err := cmd.Start(); err != nil {
				pterm.Fatal.WithFatal(false).Println("启动 Sign Server 时出现错误")
				panic(err)
			}
			err := cmd.Wait()
			if err != nil {
				pterm.Error.Println(err)
			}
			cmd.Process.Kill()
			pterm.Warning.Println("重新启动 Sign Server..")
		}
	}()
	// 等待连接
	serverUrl := fmt.Sprintf("http://127.0.0.1:%d", availablePort)
	for !checkHTTPConnection(serverUrl) {
		time.Sleep(time.Second)
	}
	pterm.Success.Println("Sign Server 已成功启动")
	return serverUrl
}
