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

func getJavaCache(fileName string) []byte {
	localFilePath := filepath.Join(utils.GetCacheDir(), "downloads", fileName)
	if !utils.IsFile(localFilePath) {
		return nil
	}
	pterm.Warning.Printfln("正在检查 Java 运行环境文件..")
	if string(utils.DownloadSmallContent(javaDownloadUrl+fileName+".sha256")) != utils.GetFileHash(localFilePath) {
		return nil
	}
	result, err := utils.GetFileData(localFilePath)
	if err != nil {
		panic(err)
	}
	return result
}

func CheckJava() bool {
	cmd := exec.Command("java", "-version")
	return cmd.Run() == nil
}

func JavaDeploy() {
	if utils.IsDir(filepath.Join(utils.GetCacheDir(), "Java", "jdk-20.0.1")) {
		return
	}
	fileName := getJavaCompressFileName()
	fileContent := getJavaCache(fileName)
	if fileContent == nil {
		pterm.Warning.Printfln("正在下载 Java 运行环境文件..")
		fileContent = utils.DownloadSmallContent(javaDownloadUrl + fileName)
		utils.WriteFileData(filepath.Join(utils.GetCacheDir(), "downloads", fileName), fileContent)
	}
	pterm.Warning.Printfln("正在解压 Java 运行环境文件..")
	dstDir := filepath.Join(utils.GetCacheDir(), "Java")
	if strings.HasSuffix(fileName, ".tar.gz") {
		if err := utils.UnTarGz(fileContent, dstDir); err != nil {
			utils.RemoveFile(filepath.Join(utils.GetCacheDir(), "downloads", fileName))
			panic(err)
		}
		os.Chmod(filepath.Join(utils.GetCacheDir(), "Java", "jdk-20.0.1", "bin", "java"), 0755)
	} else {
		if err := utils.UnZip(fileContent, dstDir); err != nil {
			utils.RemoveFile(filepath.Join(utils.GetCacheDir(), "downloads", fileName))
			panic(err)
		}
	}
	pterm.Success.Printfln("Java 运行环境已成功部署")
}
