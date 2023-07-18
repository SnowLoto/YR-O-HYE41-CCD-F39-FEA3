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

var (
	// 使用此项目搭建: https://github.com/hunshcn/gh-proxy
	MIRROR_URL = "https://www.omega-download.top/"
)

func DownloadBytes(sourceUrl string) ([]byte, error) {
	// Get the data
	resp, err := http.Get(sourceUrl)
	if err != nil {
		return nil, err
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
		return nil, err
	}
	return contents.Bytes(), nil
}

func DownloadFile(sourceUrl string, destinationPath string) error {
	// 获取目录路径
	destinationDir := filepath.Dir(destinationPath)
	// 创建目录
	if !MkDir(destinationDir) {
		return errors.New("创建目录时出现错误")
	}
	// 发起HTTP GET请求
	resp, err := http.Get(sourceUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("下载失败，状态码：%d", resp.StatusCode)
	}
	// 获取文件大小
	size, err := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		return err
	}
	// 创建文件
	file, err := os.Create(destinationPath)
	if err != nil {
		return err
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
		return err
	}
	return nil
}

func DownloadBytesWithMirror(sourceUrl string) ([]byte, error) {
	if bytes, err := DownloadBytes(MIRROR_URL + sourceUrl); err == nil {
		return bytes, nil
	}
	pterm.Warning.Println("下载时出现错误, 将再次尝试下载..")
	if bytes, err := DownloadBytes(sourceUrl); err == nil {
		return bytes, nil
	} else {
		panic(err)
	}
}

func DownloadFileWithMirror(sourceUrl string, destinationPath string) error {
	if DownloadFile(MIRROR_URL+sourceUrl, destinationPath) == nil {
		return nil
	}
	pterm.Warning.Println("下载时出现错误, 将再次尝试下载..")
	if err := DownloadFile(sourceUrl, destinationPath); err != nil {
		return err
	}
	return nil
}
