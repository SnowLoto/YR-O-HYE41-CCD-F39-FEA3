package utils

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
	"github.com/pterm/pterm"
	"golang.org/x/term"
)

// CONF, 与 INFO 相似的样式
var ConfPrinter = pterm.PrefixPrinter{
	MessageStyle: &pterm.ThemeDefault.InfoMessageStyle,
	Prefix: pterm.Prefix{
		Style: &pterm.ThemeDefault.InfoPrefixStyle,
		Text:  "CONF",
	},
}

func GetInput() string {
	buf := bufio.NewReader(os.Stdin)
	l, _, _ := buf.ReadLine()
	return string(strings.TrimSpace(string(l)))
}

func GetValidInput(text string) string {
	for {
		ConfPrinter.Print(text, ": ")
		result := GetInput()
		if result == "" {
			pterm.Error.Println("无效输入, 输入不能为空")
			continue
		}
		return result
	}
}

func GetPswInput(text string) string {
	ConfPrinter.Printf(text + " (不会回显): ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	pterm.Println()
	if err != nil {
		panic(err)
	}
	return string(bytePassword)
}

func GetInputYN(text string) bool {
	confirm := pterm.DefaultInteractiveConfirm
	// 设置默认值为 Y
	confirm.DefaultValue = true
	// 修改待选项为黄色
	confirm.SuffixStyle = &pterm.Style{pterm.FgYellow}
	// 显示并返回用户输入
	result, _ := confirm.Show(ConfPrinter.Sprint(text))
	return result
}

func GetInputYNInTime(text string, sec int32) bool {
	isInput := false
	// 自动确认
	go func() {
		time.Sleep(time.Second * time.Duration(sec))
		if !isInput {
			keyboard.SimulateKeyPress(keys.Enter)
		}
	}()
	confirm := pterm.DefaultInteractiveConfirm
	// 设置默认值为 Y
	confirm.DefaultValue = true
	// 修改待选项为黄色
	confirm.SuffixStyle = &pterm.Style{pterm.FgYellow}
	// 显示并返回用户输入
	result, _ := confirm.Show(ConfPrinter.Sprint(text) + pterm.Yellow(pterm.Sprintf(" [%d秒后自动确认]", sec)))
	// 取消自动确认
	isInput = true
	return result
}

func GetIntInputInScope(text string, a, b int) int {
	for {
		s := GetValidInput(text)
		num, err := strconv.Atoi(s)
		if err != nil {
			pterm.Error.Println("只能输入数字, 请重新输入")
			continue
		}
		if num < a || num > b {
			pterm.Error.Printfln("只能输入%d到%d之间的整数, 请重新输入", a, b)
			continue
		}
		return num
	}
}
