package fastbuilder

import (
	"bufio"
	"fmt"
	"io"
	"omega_launcher/defines"
	"omega_launcher/utils"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"

	"github.com/pterm/pterm"
)

// 保存配置文件
func SaveConfig(cfg *defines.LauncherConfig) {
	if err := utils.WriteJsonData(path.Join(utils.GetCurrentDataDir(), "服务器登录配置.json"), cfg); err != nil {
		pterm.Error.Println("无法记录配置, 不过可能不是什么大问题")
	}
}

// 配置Token
func FBTokenSetup(cfg *defines.LauncherConfig) {
	if cfg.FBToken != "" {
		if utils.GetInputYN("要使用上次的 Fastbuilder 账号登录吗?") {
			return
		}
	}
	cfg.FBToken = requestToken()
}

// 配置租赁服信息
func RentalServerSetup(cfg *defines.LauncherConfig) {
	cfg.RentalCode = utils.GetValidInput("请输入租赁服号")
	cfg.RentalPasswd = utils.GetPswInput("请输入租赁服密码")
}

func setupCmdArgs(cfg *defines.LauncherConfig) []string {
	// 配置启动参数
	args := []string{"-M", "--plain-token", cfg.FBToken, "--no-update-check", "-c", cfg.RentalCode}
	// 是否需要租赁服密码
	if cfg.RentalPasswd != "" {
		args = append(args, "-p")
		args = append(args, cfg.RentalPasswd)
	}
	// 是否启动Omega
	if cfg.StartOmega {
		args = append(args, "-O")
	}
	return args
}

func Run(cfg *defines.LauncherConfig) {
	// 获取启动器版本信息
	go func() {
		latestVer := utils.GetLauncherUpdateInfo()
		if latestVer != "" {
			cfg.LatestVer = latestVer
			SaveConfig(cfg)
		}
	}()
	// 启动成功提示语
	successTip := "已成功进入租赁服"
	// 获取命令args
	args := setupCmdArgs(cfg)
	// 建立频道
	readC := make(chan string)
	stop := make(chan string)
	// 持续将输入信息输入到频道中
	go func() {
		for {
			s := utils.GetInput()
			readC <- s
		}
	}()
	// 重启间隔
	restartTime := 0
	// 上次输入
	lastInput := ""
	// 监听程序退出信号
	exitSignal := make(chan os.Signal, 1)
	signal.Notify(exitSignal, os.Interrupt)
	signal.Notify(exitSignal, syscall.SIGTERM)
	signal.Notify(exitSignal, syscall.SIGQUIT)
	for {
		// 记录启动时间
		startTime := time.Now()
		// 是否已正常启动
		isStarted := false
		// 是否停止
		isStopped := false
		// 启动时提示信息
		pterm.Success.Println("正在启动 Omega/Fastbuilder")
		// 启动命令
		cmd := exec.Command(getFBExecPath(), args...)
		cmd.Dir = path.Join(utils.GetCurrentDataDir())
		// 由于需要对内容进行处理, 所以不能直接进行io复制
		// 建立从Fastbuilder到控制台的错误管道
		omega_err, err := cmd.StderrPipe()
		if err != nil {
			panic(err)
		}
		// 建立从Fastbuilder到控制台的输出管道
		omega_out, err := cmd.StdoutPipe()
		if err != nil {
			panic(err)
		}
		// 建立从控制台到Fastbuilder的输入管道
		omega_in, err := cmd.StdinPipe()
		if err != nil {
			panic(err)
		}
		// 仅启动 FB 时需要额外提示
		if !cfg.StartOmega {
			omega_in.Write([]byte("say " + successTip + "\n"))
		}
		// 从管道中获取并打印Fastbuilder错误内容
		go func() {
			reader := bufio.NewReader(omega_err)
			for {
				readString, err := reader.ReadString('\n')
				if readString == "\n" {
					continue
				}
				if err != nil || err == io.EOF {
					//pterm.Error.Println("读取 Omega/Fastbuilder 错误内容时出现错误")
					return
				}
				// 如果程序正在退出则不再发送错误信息
				if isStopped {
					return
				}
				fmt.Print(readString + "\033[0m")
			}
		}()
		// 从管道中获取并打印Fastbuilder输出内容
		go func() {
			reader := bufio.NewReader(omega_out)
			for {
				readString, err := reader.ReadString('\n')
				readString = strings.TrimPrefix(readString, "> ")
				if readString == "\n" || readString == lastInput+"\n" || readString == "say "+successTip+"\n" {
					lastInput = ""
					continue
				}
				if err != nil || err == io.EOF {
					//pterm.Error.Println("读取 Omega/Fastbuilder 输出内容时出现错误")
					return
				}
				fmt.Print(readString + "\033[0m")
				// 成功启动后处理
				if !isStarted && (readString == successTip+"\n" || readString == "Starting Omega in a second\n") {
					isStarted = true
					// 读取验证服务器返回的Token并保存
					cfg.FBToken = loadCurrentFBToken()
					if cfg.FBToken != "" {
						SaveConfig(cfg)
					}
				}
			}
		}()
		// 在未收到停止信号前, 启动器会一直将控制台输入的内容通过管道发送给Fastbuilder
		go func() {
			for {
				select {
				case <-stop:
					return
				case <-exitSignal:
					// 关闭重启
					isStopped = true
					// 正常运行时使用关闭命令退出程序
					if isStarted {
						if cfg.StartOmega {
							omega_in.Write([]byte("stop\n"))
						} else {
							omega_in.Write([]byte("exit\n"))
						}
						return
					}
					// 强制退出
					os.Exit(1)
				case s := <-readC:
					lastInput = s
					// 接收到停止命令时处理
					if (cfg.StartOmega && s == "stop") || s == "exit" || s == "fbexit" {
						// 关闭重启
						isStopped = true
						// 发出停止命令
						omega_in.Write([]byte(s + "\n"))
						// 输出信息
						pterm.Success.Println("正在等待 Omega/Fastbuilder 处理退出命令")
						// 停止接收输入
						return
					} else {
						omega_in.Write([]byte(s + "\n"))
					}
				}
			}
		}()
		// 启动并持续运行Fastbuilder
		err = cmd.Start()
		if !isStopped && err != nil {
			pterm.Error.Println("Omega/Fastbuilder 启动时出现错误")
			panic(err)
		}
		err = cmd.Wait()
		// 如果运行到这里, 说明Fastbuilder出现错误或退出运行了
		cmd.Process.Kill()
		// 判断是否正常退出
		if isStopped {
			pterm.Success.Println("Omega/Fastbuilder 已正常退出, 启动器将结束运行")
			time.Sleep(time.Second * 3)
			break
		} else {
			stop <- "stop!!"
			pterm.Error.Printfln("Oh no! Fastbuilder crashed! (%s)", err.Error()) // ?
		}
		// 为了避免频繁请求, 崩溃后将等待一段时间后重启, 可手动跳过等待
		if time.Since(startTime) < time.Minute {
			if restartTime < 3600 {
				restartTime = restartTime + 60
			}
		} else {
			restartTime = 0
		}
		pterm.Warning.Printf("Omega/Fastbuilder 将在%d秒后重新启动 (按回车立即重启)", restartTime)
		// 等待输入或计时结束
		select {
		case <-readC:
			restartTime = 0
		case <-time.After(time.Second * time.Duration(restartTime)):
			// 换行
			fmt.Println()
		}
	}
}
