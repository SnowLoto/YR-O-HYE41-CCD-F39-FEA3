package main

import (
	_ "embed"
	"fmt"
	"omega_launcher/cqhttp"
	"omega_launcher/defines"
	"omega_launcher/fastbuilder"
	"omega_launcher/launcher"
	"omega_launcher/plantform"
	"omega_launcher/utils"
	"os"
	"path"
	"time"

	"github.com/pterm/pterm"
)

//go:embed VERSION
var version []byte

func beforeClose() {
	// 打印错误
	err := recover()
	if err != nil {
		pterm.Fatal.WithFatal(false).Println(err)
	}
	// Make Windows users happy
	if p := plantform.GetPlantform(); p == plantform.WINDOWS_x86_64 || p == plantform.WINDOWS_arm64 {
		time.Sleep(time.Second * 5)
	}
}

func main() {
	defer beforeClose()
	// 确保目录可用
	if err := os.Chdir(utils.GetCurrentDir()); err != nil {
		panic(err)
	}
	// 启动
	// 读取配置
	launcherConfig := &defines.LauncherConfig{}
	utils.GetJsonData(path.Join(utils.GetCurrentDataDir(), "服务器登录配置.json"), launcherConfig)
	// 获取启动器版本信息 (异步)
	go func() {
		latestVer := launcher.GetLauncherUpdateInfo()
		if latestVer != "" {
			launcherConfig.LatestVer = latestVer
			launcher.SaveConfig(launcherConfig)
		}
	}()
	// 版本对比, 不一致时提示更新可用
	verInfo := "Omega Launcher" + pterm.Yellow(" (", string(version), ")")
	if utils.HasGreaterVer(string(version), launcherConfig.LatestVer) {
		verInfo += pterm.Yellow(" (更新可用)")
	}
	// 添加启动信息
	pterm.DefaultBox.Println("https://github.com/Liliya233/omega_launcher")
	pterm.Info.Println(verInfo)
	pterm.Info.Println("Author: CMA2401PT, Modified by Liliya233")
	// 询问是否使用上一次的配置
	if fastbuilder.CheckExecFile() && launcherConfig.FBToken != "" && launcherConfig.RentalCode != "" {
		if utils.GetInputYNInTime("要使用和上次完全相同的配置启动吗?", 10) {
			// 更新FB
			if launcherConfig.UpdateFB {
				fastbuilder.Update(launcherConfig, false)
			}
			// go-cqhttp
			if launcherConfig.EnableCQHttp && launcherConfig.StartOmega {
				cqhttp.Run(launcherConfig)
			}
			// 启动Omega或者FB
			fastbuilder.Run(launcherConfig)
			return
		}
	}
	// 配置FB更新
	if launcherConfig.UpdateFB = utils.GetInputYN("需要启动器帮忙下载或更新 Fastbuilder 吗?"); launcherConfig.UpdateFB {
		fastbuilder.Update(launcherConfig, true)
	}
	// 检查是否已下载FB
	if !fastbuilder.CheckExecFile() {
		pterm.Warning.Printfln("当前目录不存在文件名为 " + fastbuilder.GetFBExecName() + " 的 Fastbuilder")
		fastbuilder.Update(launcherConfig, true)
	}
	// 配置FB
	fastbuilder.FBTokenSetup(launcherConfig)
	// 配置租赁服登录 (如果不为空且选择使用上次配置, 则跳过setup)
	if !(launcherConfig.RentalCode != "" && utils.GetInputYN(fmt.Sprintf("要使用上次 %s 的租赁服配置吗?", launcherConfig.RentalCode))) {
		launcher.RentalServerSetup(launcherConfig)
	}
	// 询问是否使用Omega
	if launcherConfig.StartOmega = utils.GetInputYN("需要启动 Omega 吗?"); launcherConfig.StartOmega {
		// 配置群服互通
		if launcherConfig.EnableCQHttp = utils.GetInputYN("需要启动 go-cqhttp 吗?"); launcherConfig.EnableCQHttp {
			if !utils.IsDir(path.Join(fastbuilder.GetOmegaStorageDir(), "配置")) {
				if launcherConfig.EnableCQHttp = utils.GetInputYN("此时配置 go-cqhttp 会导致新生成的组件均为非启用状态, 要继续吗?"); !launcherConfig.EnableCQHttp {
					// 直接启动Omega或者FB
					fastbuilder.Run(launcherConfig)
				}
			}
			launcherConfig.BlockCQHttpOutput = utils.GetInputYN("需要在配置完成后屏蔽 go-cqhttp 的输出吗?")
			cqhttp.CQHttpEnablerHelper()
			if launcherConfig.EnableSignServer = utils.GetInputYN("需要启动器启动 Sign Server 吗?"); launcherConfig.EnableSignServer {
				launcherConfig.SignServerSoVersion = utils.GetInput("请输入 Sign Server 依赖的文件版本 (不输入则为8.9.63)")
			}
			cqhttp.Run(launcherConfig)
		}
	}
	// 启动Omega或者FB
	fastbuilder.Run(launcherConfig)
}
