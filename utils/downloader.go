package utils

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/pterm/pterm"
)

func DownloadBytes(sourceUrl string) []byte {
	// Get the data
	resp, err := http.Get(sourceUrl)
	if err != nil {
		pterm.Fatal.WithFatal(false).Println("下载资源时出现错误")
		panic(err)
	}
	defer resp.Body.Close()
	// Size
	size, _ := strconv.Atoi(resp.Header.Get("Content-Length"))
	// Progress Bar
	bar := pb.Full.Start64(int64(size))
	bar.SetWidth(-1)
	bar.SetMaxWidth(100)
	bar.SetRefreshRate(time.Millisecond)
	defer bar.Finish()
	// Reader
	barReader := bar.NewProxyReader(resp.Body)
	// Buffer
	contents := bytes.NewBuffer([]byte{})
	if _, err := io.Copy(contents, barReader); err != nil {
		panic(err)
	}
	return contents.Bytes()
}

func DownloadFile(sourceURL string, destinationPath string) {
	// 获取目录路径
	destinationDir := filepath.Dir(destinationPath)
	// 创建目录
	if !MkDir(destinationDir) {
		err := errors.New("创建目录时出现错误")
		pterm.Fatal.WithFatal(false).Println(err.Error())
		panic(err)
	}
	// 发起HTTP GET请求
	resp, err := http.Get(sourceURL)
	if err != nil {
		pterm.Fatal.WithFatal(false).Println("下载资源时出现错误")
		panic(err)
	}
	defer resp.Body.Close()
	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		panic(fmt.Errorf("下载失败，状态码：%d", resp.StatusCode))
	}
	// 获取文件大小
	size, err := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		pterm.Fatal.WithFatal(false).Println("获取文件大小时出现错误")
		panic(err)
	}
	// 创建文件
	file, err := os.Create(destinationPath)
	if err != nil {
		pterm.Fatal.WithFatal(false).Println("创建文件时出现错误")
		panic(err)
	}
	defer file.Close()
	// 创建进度条
	bar := pb.Full.Start64(size)
	bar.SetWidth(-1)
	bar.SetMaxWidth(100)
	bar.SetRefreshRate(time.Millisecond)
	defer bar.Finish()
	// 创建进度条写入器
	writer := bar.NewProxyWriter(file)
	// 将响应体写入文件并显示进度
	_, err = io.Copy(writer, resp.Body)
	if err != nil {
		pterm.Fatal.WithFatal(false).Println("写入文件时出现错误")
		panic(err)
	}
}
