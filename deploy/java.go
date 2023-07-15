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

func getJavaCompressFileName() string {
	switch plantform.GetPlantform() {
	case plantform.Android_arm64:
		pterm.Error.Println("对于 Android Termux, 请尝试输入命令 \"pkg install openjdk-17\" 来安装 Java")
		// return "jdk-20_linux-aarch64_bin.tar.gz"
	case plantform.Android_x86_64:
		pterm.Error.Println("对于 Android Termux, 请尝试输入命令 \"pkg install openjdk-17\" 来安装 Java")
		// return "jdk-20_linux-x64_bin.tar.gz"
	case plantform.Linux_arm64:
		return "jdk-20_linux-aarch64_bin.tar.gz"
	case plantform.Linux_x86_64:
		return "jdk-20_linux-x64_bin.tar.gz"
	case plantform.MACOS_arm64:
		return "jdk-20_macos-aarch64_bin.tar.gz"
	case plantform.MACOS_x86_64:
		return "jdk-20_macos-x64_bin.tar.gz"
	case plantform.WINDOWS_arm64:
		// No support
	case plantform.WINDOWS_x86_64:
		return "jdk-20_windows-x64_bin.zip"
	}
	panic("请尝试自行安装 Java")
}

func isJavaCache(fileName, downloadDst string) bool {
	pterm.Warning.Printfln("正在检查 Java 运行环境文件..")
	return string(utils.DownloadBytes(javaDownloadUrl+fileName+".sha256")) == utils.GetFileHash(downloadDst)
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
	fileName := getJavaCompressFileName()
	downloadDst := filepath.Join(utils.GetCacheDir(), "downloads", fileName)
	if !isJavaCache(fileName, downloadDst) {
		pterm.Warning.Printfln("正在下载 Java 运行环境文件..")
		utils.DownloadFile(javaDownloadUrl+fileName, downloadDst)
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
