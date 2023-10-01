package main

import (
	_ "embed"
	"fmt"
	"omega_launcher/cqhttp"
	"omega_launcher/fastbuilder"
	"omega_launcher/launcher"
	"omega_launcher/plantform"
	"omega_launcher/utils"
	"os"
	"path/filepath"
	"runtime/debug"
	"time"

	"github.com/pterm/pterm"
	"golang.org/x/term"
)

//go:embed VERSION
var version []byte

func beforeClose() {
	// 打印错误
	err := recover()
	if err != nil {
		pterm.Fatal.WithFatal(false).Println(err)
		// Make Contributors happy
		debug.PrintStack()
	}
	if p := plantform.GetPlantform(); p == plantform.WINDOWS_amd64 || p == plantform.WINDOWS_arm64 {
		// Make Windows users happy
		time.Sleep(time.Second * 5)
	} else {
		// Make Unix users happy
		term.MakeRaw(0)
	}
}

func main() {
	defer beforeClose()
	// 确保目录可用
	if err := os.Chdir(utils.GetCurrentDir()); err != nil {
		panic(err)
	}
	// 启动器自更新 (异步)
	go launcher.CheckUpdate(string(version))
	// 启动
	// 读取配置
	launcherConfig := &launcher.Config{}
	utils.GetJsonData(filepath.Join(utils.GetCurrentDataDir(), "SnowConfig.json"), launcherConfig)
	// 添加启动信息
	pterm.DefaultBox.Println("https://Snow.fastbuilder.icu/SnowLotus/")
	pterm.Info.Println("Omega Launcher" + pterm.Yellow(" (Legacy Omega Only)") + pterm.Yellow(" (", string(version), ")"))
	pterm.Info.Println("Author: CMA2401PT, Modified by SnowLotus")
	// 询问是否使用上一次的配置
	if fastbuilder.CheckExecFile() && launcherConfig.RentalCode != "" {
		if result, _ := utils.GetInputYNInTime("是否使用上一次的登录配置？", 10); result {
			// 更新FB
			if launcherConfig.UpdateFB {
				fastbuilder.Update(launcherConfig)
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
	if launcherConfig.UpdateFB = utils.GetInputYN("是否需要下载或更新PhoenixBuilder？"); launcherConfig.UpdateFB {
		fastbuilder.UpdateRepo(launcherConfig)
		fastbuilder.Update(launcherConfig)
	}
	// 检查是否已下载FB
	if !fastbuilder.CheckExecFile() {
		pterm.Warning.Printfln("Error " + plantform.GetFastBuilderName() + " no FastBuilder")
		fastbuilder.UpdateRepo(launcherConfig)
		fastbuilder.Update(launcherConfig)
	}
	// 配置FB
	fastbuilder.FBTokenSetup(launcherConfig)
	// 配置租赁服登录 (如果不为空且选择使用上次配置, 则跳过setup)
	if !(launcherConfig.RentalCode != "" && utils.GetInputYN(fmt.Sprintf("是否使用上一次的 %s 的租赁服登陆配置?", launcherConfig.RentalCode))) {
		launcherConfig.RentalCode = utils.GetValidInput("输入服务器号")
		launcherConfig.RentalPasswd = utils.GetPswInput("输入服务器密码")
	}
	// 询问是否使用Omega
	if launcherConfig.StartOmega = utils.GetInputYN("要启动 Omega 吗?"); launcherConfig.StartOmega {
		// 配置群服互通
		if launcherConfig.EnableCQHttp = utils.GetInputYN("要启动 go-cqhttp/群服互通 吗?"); launcherConfig.EnableCQHttp {
			if !utils.IsDir(filepath.Join(fastbuilder.GetOmegaStorageDir(), "配置")) {
				if launcherConfig.EnableCQHttp = utils.GetInputYN("此时配置 go-cqhttp 会导致新生成的组件均为非启用状态, 要继续吗?"); !launcherConfig.EnableCQHttp {
					// 直接启动Omega或者FB
					fastbuilder.Run(launcherConfig)
					return
				}
			}
			launcherConfig.BlockCQHttpOutput = utils.GetInputYN("要在配置完成后屏蔽 go-cqhttp 的输出吗?")
			cqhttp.CQHttpEnablerHelper()
			launcherConfig.EnableSignServer = utils.GetInputYN("要启动器启动 Sign Server 吗?")
			cqhttp.Run(launcherConfig)
		}
	}
	// 启动Omega或者FB
	fastbuilder.Run(launcherConfig)
}
