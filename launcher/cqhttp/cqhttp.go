package cqhttp

import (
	"bufio"
	"fmt"
	"io"
	"omega_launcher/defines"
	"omega_launcher/fastbuilder"
	"omega_launcher/utils"
	"os"
	"os/exec"
	"path"

	"github.com/pterm/pterm"
)

func CQHttpEnablerHelper() {
	// 无法创建目录则退出
	if !utils.MkDir(GetCQHttpDir()) {
		panic("无法创建 cqhttp_storage 目录")
	}
	// 导入配置
	UnPackCQHttpRunAuth()
	// 如果go-cqhttp配置文件不存在, 则执行初始化操作
	if cqCfg := getCQConfig(); cqCfg != nil {
		if utils.GetInputYN("已读取到 go-cqhttp 配置文件, 要使用吗?") {
			return
		}
	}
	initCQConfig()
}

func Run(botCfg *defines.LauncherConfig) {
	// 不存在cqhttp目录则退出
	if !utils.IsDir(GetCQHttpDir()) {
		panic("cqhttp_storage 目录不存在, 请使用启动器配置一次 go-cqhttp")
	}
	// 如果不存在cqhttp程序则解压
	if !utils.IsFile(path.Join(GetCqHttpExec())) {
		if err := utils.WriteFileData(GetCqHttpExec(), GetCqHttpBinary()); err != nil {
			pterm.Fatal.WithFatal(false).Println("解压 go-cqhttp 时遇到问题")
			panic(err)
		}
	}
	// 读取Omega配置
	utils.MkDir(path.Join(fastbuilder.GetOmegaStorageDir(), "配置"))
	// 配置地址同步
	port, err := utils.GetAvailablePort()
	if err != nil {
		pterm.Fatal.WithFatal(false).Println("无法为 go-cqhttp 获取可用端口")
		panic(err)
	}
	availableAddress := fmt.Sprintf("127.0.0.1:%d", port)
	qGroupCfgFp, qGuildCfgFp := updateOmegaConfigAddress(availableAddress)
	updateCQConfigAddress(availableAddress)
	// 配置启动参数
	args := []string{"-faststart"}
	// 配置执行目录
	cmd := exec.Command(GetCqHttpExec(), args...)
	cmd.Dir = path.Join(GetCQHttpDir())
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cqhttp_out, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}
	// 从管道中获取并打印cqhttp输出内容
	stopOutput := false
	go func() {
		reader := bufio.NewReader(cqhttp_out)
		for {
			readString, err := reader.ReadString('\n')
			if err != nil || err == io.EOF {
				return
			}
			if stopOutput {
				return
			}
			fmt.Print(readString + "\033[0m")
		}
	}()
	// 启动并持续运行cqhttp
	go func() {
		pterm.Success.Println("正在启动 go-cqhttp")
		err := cmd.Start()
		if err != nil {
			pterm.Fatal.WithFatal(false).Println("go-cqhttp 启动时出现错误")
			panic(err)
		}
		err = cmd.Wait()
		if err != nil {
			pterm.Fatal.WithFatal(false).Println("go-cqhttp 运行时出现错误")
			panic(err)
		}
	}()
	// 等待cqhttp启动完成
	WaitConnect(availableAddress)
	// 配置完成后, 根据设置决定是否关闭go-cqhttp输出
	if botCfg.BlockCQHttpOutput {
		pterm.Warning.Println("将屏蔽 go-cqhttp 的输出内容")
		stopOutput = true
	}
	// 打包data文件
	PackCQHttpRunAuth(qGroupCfgFp, qGuildCfgFp)
	// 打印提示消息
	pterm.Info.Println(
		" 若要为服务器配置 go-cqhttp, 请执行以下的操作：\n",
		"1. 在服务器成功启动一次 Omega, 以生成 omega_storage 目录\n",
		"2. 将 上传这个文件到云服务器以使用云服务器的群服互通.data 上传至服务器 omega_storage 目录下\n",
		"3. 重启启动器并选择启动 go-cqhttp, 此时应该能够识别到 data 文件了\n",
		"如果遇到 go-cqhttp 相关的问题, 可前往 https://docs.go-cqhttp.org/ 寻找可用信息",
	)
}
