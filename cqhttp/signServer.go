package cqhttp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"omega_launcher/deploy"
	"omega_launcher/plantform"
	"omega_launcher/remote"
	"omega_launcher/utils"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/pterm/pterm"
)

// config.json of sign server
type Config struct {
	Server struct {
		Host string `json:"host"`
		Port int    `json:"port"`
	} `json:"server"`
	Key            string `json:"key"`
	AutoRegister   bool   `json:"auto_register"`
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

// 获取下载文件路径
func getDownloadFilePath() string {
	// 获取远程数据
	remote := remote.GetSignServerRemoteData()
	// 下载文件夹
	return filepath.Join(utils.GetCacheDir(), "downloads", remote.ZipName)
}

// 检查下载文件是否已下载
func isCache() bool {
	// 获取远程数据
	remote := remote.GetSignServerRemoteData()
	// 检查大小是否一致
	return utils.GetFileSize(getDownloadFilePath()) == remote.ZipSize
}

func signServerDeploy() {
	// 获取远程数据
	remote := remote.GetSignServerRemoteData()
	// 缓存与下载
	if !isCache() {
		pterm.Warning.Printfln("正在下载 Sign Server 可执行文件..")
		utils.DownloadFileWithMirror(remote.DownloadUrl, getDownloadFilePath())
	}
	// 解压
	if !utils.IsDir(filepath.Join(utils.GetCacheDir(), "SignServer", remote.UnzipDirName)) {
		fp, err := os.OpenFile(getDownloadFilePath(), os.O_RDONLY, 0755)
		if err != nil {
			panic(err)
		}
		if err := utils.UnZip(fp, filepath.Join(utils.GetCacheDir(), "SignServer")); err != nil {
			utils.RemoveFile(getDownloadFilePath())
			panic(err)
		}
	}
}

func SignServerStart() string {
	// 检查device.json是否存在
	deviceJsonPath := filepath.Join(GetCQHttpDir(), "device.json")
	if !utils.IsFile(deviceJsonPath) {
		pterm.Warning.Println("尚未生成 device.json, 将跳过 Sign Server 的配置")
		return "-"
	}
	// 检查是否为Android协议
	deviceJsonContent, _ := utils.GetFileData(deviceJsonPath)
	var device struct {
		Protocol  int    `json:"protocol"`
		AndroidId string `json:"android_id"`
	}
	err := json.Unmarshal(deviceJsonContent, &device)
	if err != nil {
		pterm.Error.Println("读取 device.json 时出现错误, 将跳过 Sign Server 的配置: " + err.Error())
		return "-"
	}
	if device.Protocol != 1 && device.Protocol != 6 {
		pterm.Error.Println("未使用 Android 登录协议 (1, 6), 将跳过 Sign Server 的配置")
		return "-"
	}
	// 获取远程数据
	remote := remote.GetSignServerRemoteData()
	// 部署SignServer
	signServerDeploy()
	pterm.Success.Printfln("已成功部署最新版本的 Sign Server, 正在尝试启动程序..")
	// 获取可用端口
	availablePort, err := utils.GetAvailablePort()
	if err != nil {
		pterm.Error.Println("无法获取可用的端口来启动 Sign Server: " + err.Error())
		return "-"
	}
	// 获取uin
	if cqCfg := getCQConfig(); cqCfg != nil {
		configPath := filepath.Join(utils.GetCacheDir(), "SignServer", remote.UnzipDirName, "txlib", remote.SoVersion, "config.json")
		setupConfig(configPath, "0.0.0.0", availablePort, cqCfg.Account.Uin, device.AndroidId)
	} else {
		pterm.Error.Println("Sign Server 启动失败, 未能够从 go-cqhttp 配置文件中获取QQ账号")
		return "-"
	}
	// 如果不是Windows则去掉.bat
	cmdStr := filepath.Join(utils.GetCacheDir(), "SignServer", remote.UnzipDirName, "bin", "unidbg-fetch-qsign.bat")
	if currentPlanform := plantform.GetPlantform(); currentPlanform != plantform.WINDOWS_x86_64 && currentPlanform != plantform.WINDOWS_arm64 {
		cmdStr = strings.TrimSuffix(cmdStr, ".bat")
		os.Chmod(cmdStr, 0755)
	}
	// 启动命令
	cmd := exec.Command(cmdStr, []string{fmt.Sprintf("--basePath=%s", filepath.Join(utils.GetCacheDir(), "SignServer", remote.UnzipDirName, "txlib", remote.SoVersion))}...)
	// 配置Java运行环境
	if !deploy.CheckJava() {
		pterm.Warning.Println("系统未安装Java, 将尝试安装对应的 Java 环境")
		deploy.JavaDeploy()
		cmd.Env = append(os.Environ(), fmt.Sprintf("JAVA_HOME=%s", filepath.Join(utils.GetCacheDir(), "Java", "jdk-20.0.1")))
	}
	// 运行SignServer
	go func() {
		if err := cmd.Start(); err != nil {
			pterm.Fatal.WithFatal(false).Println("启动 Sign Server 时出现错误")
			panic(err)
		}
		if err := cmd.Wait(); err != nil {
			pterm.Fatal.WithFatal(false).Println("Sign Server 在运行过程中出现错误")
			panic(err)
		}
	}()
	// 等待连接
	serverUrl := fmt.Sprintf("http://127.0.0.1:%d", availablePort)
	checkHTTPConnection := func() bool {
		client := http.Client{
			Timeout: time.Second,
		}
		resp, err := client.Get(serverUrl)
		if err != nil {
			return false
		}
		defer resp.Body.Close()
		return resp.StatusCode == http.StatusOK
	}
	for !checkHTTPConnection() {
		time.Sleep(time.Second)
	}
	pterm.Success.Println("已成功启动 Sign Server")
	return serverUrl
}
