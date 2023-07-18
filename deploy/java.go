package deploy

import (
	"omega_launcher/plantform"
	"omega_launcher/utils"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pterm/pterm"
)

var javaDownloadUrl = "https://download.oracle.com/java/20/latest/"

func isJavaCache(fileName, downloadDst string) bool {
	pterm.Warning.Printfln("正在检查 Java 运行环境文件..")
	remoteHash, err := utils.DownloadBytes(javaDownloadUrl + fileName + ".sha256")
	if err != nil {
		panic(err)
	}
	return string(remoteHash) == utils.GetFileHash(downloadDst)
}

func CheckJava() bool {
	cmd := exec.Command("java", "-version")
	return cmd.Run() == nil
}

func JavaDeploy() {
	javaExecFile := filepath.Join(utils.GetCacheDir(), "Java", "jdk-20.0.1", "bin", "java")
	if utils.IsFile(javaExecFile) {
		return
	}
	fileName := plantform.GetJDKDownloadName()
	downloadDst := filepath.Join(utils.GetCacheDir(), "downloads", fileName)
	if !isJavaCache(fileName, downloadDst) {
		pterm.Warning.Printfln("正在下载 Java 运行环境文件..")
		err := utils.DownloadFile(javaDownloadUrl+fileName, downloadDst)
		if err != nil {
			panic(err)
		}
	}
	pterm.Warning.Printfln("正在解压 Java 运行环境文件..")
	fp, err := os.OpenFile(downloadDst, os.O_RDONLY, 0755)
	if err != nil {
		panic(err)
	}
	uncompressDir := filepath.Join(utils.GetCacheDir(), "Java")
	if strings.HasSuffix(fileName, ".tar.gz") {
		if err := utils.UnTarGz(fp, uncompressDir); err != nil {
			utils.RemoveFile(downloadDst)
			panic(err)
		}
	} else {
		if err := utils.UnZip(fp, uncompressDir); err != nil {
			utils.RemoveFile(downloadDst)
			panic(err)
		}
	}
	pterm.Success.Printfln("Java 运行环境已成功部署")
}
