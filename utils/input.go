package utils

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
	"github.com/pterm/pterm"
)

// CONF, 与 INFO 相似的样式
var (
	IsFirstPsw  = true
	ConfPrinter = pterm.PrefixPrinter{
		MessageStyle: &pterm.ThemeDefault.InfoMessageStyle,
		Prefix: pterm.Prefix{
			Style: &pterm.ThemeDefault.InfoPrefixStyle,
			Text:  "CONF",
		},
	}
)

func ReadLine() string {
	buf := bufio.NewReader(os.Stdin)
	l, _, _ := buf.ReadLine()
	return string(strings.TrimSpace(string(l)))
}

func GetInput(text string) string {
	ConfPrinter.Print(text, ": ")
	result := ReadLine()
	return result
}

func GetValidInput(text string) string {
	for {
		ConfPrinter.Print(text, ": ")
		result := ReadLine()
		if result == "" {
			pterm.Error.Println("无效输入, 输入不能为空")
			continue
		}
		return result
	}
}

func GetPswInput(text string) string {
	if !IsFirstPsw {
		// 光标上移一行
		fmt.Print("\033[1A")
		fmt.Print("\r")
	}
	IsFirstPsw = false
	result, err := pterm.DefaultInteractiveTextInput.WithMask("*").Show(ConfPrinter.Sprintf(text + " (不会回显)"))
	if err != nil {
		panic(err)
	}
	return result
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

func GetInputYNInTime(text string, sec int32) (bool, bool) {
	isFinish := false
	isUser := true
	// 自动确认
	go func() {
		time.Sleep(time.Second * time.Duration(sec))
		if !isFinish {
			keyboard.SimulateKeyPress(keys.Enter)
			isUser = false
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
	isFinish = true
	return result, isUser
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

func GetInt64Input(text string) int64 {
	for {
		s := GetValidInput(text)
		num, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			pterm.Error.Println("只能输入数字, 请重新输入")
			continue
		}
		return num
	}
}
