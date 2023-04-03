package utils

import (
	"omega_launcher/plantform"
	"os"
	"path"
	"path/filepath"
)

func GetCurrentDir() string {
	// 兼容配套的Docker
	if IsFile(path.Join("/ome", "launcher_liliya")) {
		return path.Join("/workspace")
	}
	pathExecutable, err := os.Executable()
	if err != nil {
		panic(err)
	}
	return filepath.Dir(pathExecutable)
}

func GetCurrentDataDir() string {
	// Android环境下, 尝试将数据文件放在 /sdcard/Download
	currentPlantform := plantform.GetPlantform()
	if currentPlantform == plantform.Android_arm64 || currentPlantform == plantform.Android_x86_64 {
		if IsDir("/sdcard/Download/omega_storage") {
			return path.Join("/sdcard/Download")
		} else {
			if IsDir("/sdcard") {
				if MkDir("/sdcard/Download/omega_storage") {
					return path.Join("/sdcard/Download")
				}
			}
		}
	}
	return GetCurrentDir()
}
