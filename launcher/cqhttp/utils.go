package cqhttp

import (
	"bytes"
	"io"
	"net/url"
	"omega_launcher/embed_binary"
	"omega_launcher/fastbuilder"
	"omega_launcher/utils"
	"os"
	"path"
	"path/filepath"

	"github.com/andybalholm/brotli"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/pterm/pterm"
)

func GetCqHttpExec() string {
	cqhttp := "cqhttp"
	if embed_binary.GetPlantform() == embed_binary.WINDOWS_x86_64 {
		cqhttp = "cqhttp.exe"
	}
	cqhttp = path.Join(utils.GetCurrentDir(), cqhttp)
	p, err := filepath.Abs(cqhttp)
	if err != nil {
		panic(err)
	}
	return p
}

func GetCqHttpBinary() []byte {
	compressedData := embed_binary.GetCqHttpBinary()
	var execBytes []byte
	var err error
	if execBytes, err = io.ReadAll(brotli.NewReader(bytes.NewReader(compressedData))); err != nil {
		panic(err)
	}
	return execBytes
}

func GetCQHttpHash() string {
	exec := GetCqHttpExec()
	return utils.GetFileHash(exec)
}

func GetEmbeddedCQHttpHash() string {
	return utils.GetBinaryHash(GetCqHttpBinary())
}

func WaitConnect(addr string) {
	for {
		u := url.URL{Scheme: "ws", Host: addr}
		if _, _, err := websocket.DefaultDialer.Dial(u.String(), nil); err != nil {
			// time.Sleep(1)
			continue
		} else {
			return
		}
	}
}

func GetCQHttpDir() string {
	return path.Join(utils.GetCurrentDir(), "cqhttp_storage")
}

func PackCQHttpRunAuth(qGroupLinkFp, qGuildLinkFp string) {
	_uuid, _ := uuid.NewUUID()
	uuid := _uuid.String()
	uuidFile := path.Join(GetCQHttpDir(), "uuid")
	if err := utils.WriteFileData(uuidFile, []byte(uuid)); err != nil {
		panic(err)
	}
	if _, err := utils.CopyFile(qGroupLinkFp, path.Join(GetCQHttpDir(), "组件-群服互通.json")); err != nil {
		panic(err)
	}
	if _, err := utils.CopyFile(qGuildLinkFp, path.Join(GetCQHttpDir(), "组件-第三方__Liliya233__频服互通.json")); err != nil {
		panic(err)
	}
	fileName := path.Join(fastbuilder.GetOmegaStorageDir(), "上传这个文件到云服务器以使用云服务器的群服互通.data")
	fp, err := os.OpenFile(fileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
	fp.Write([]byte(uuid))
	if err != nil {
		panic(err)
	}
	if err := utils.Zip(GetCQHttpDir(), fp, []string{"data", "logs"}); err != nil {
		panic(err)
	}
	fp.Close()
	os.Remove(path.Join(GetCQHttpDir(), "组件-群服互通.json"))
	os.Remove(path.Join(GetCQHttpDir(), "组件-第三方__Liliya233__频服互通.json"))
}

func UnPackCQHttpRunAuth() {
	fileName := path.Join(fastbuilder.GetOmegaStorageDir(), "上传这个文件到云服务器以使用云服务器的群服互通.data")
	if utils.IsFile(fileName) {
		var fp *os.File
		defer func() {
			if fp != nil {
				fp.Close()
			}
		}()
		unzipSize, err := utils.GetUnZipSize(fileName)
		if err != nil {
			panic(err)
		}
		fp, err = os.OpenFile(fileName, os.O_RDONLY, 0755)
		if err != nil {
			panic(err)
		}
		uuidBytes := make([]byte, 36)
		if _, err := fp.Read(uuidBytes); err != nil {
			panic(err)
		}
		uuidFile := path.Join(GetCQHttpDir(), "uuid")
		if utils.IsFile(uuidFile) {
			if thisUUidBytes, err := utils.GetFileData(uuidFile); err == nil {
				if string(thisUUidBytes) == string(uuidBytes) {
					return
				}
			}
		}
		pterm.Info.Print("已读取到 .data 文件，要导入吗? 要请输入 y, 不要请输入 n: ")
		accept := utils.GetInputYN()
		if accept {
			os.RemoveAll(GetCQHttpDir())
			zipData, err := io.ReadAll(fp)
			if err != nil {
				panic(err)
			}
			if err := utils.UnZip(bytes.NewReader(zipData), unzipSize, GetCQHttpDir()); err != nil {
				panic(err)
			}
			if _, err := utils.CopyFile(path.Join(GetCQHttpDir(), "组件-群服互通.json"), path.Join(fastbuilder.GetOmegaStorageDir(), "配置", "群服互通", "组件-群服互通.json")); err != nil {
				panic(err)
			}
			if _, err := utils.CopyFile(path.Join(GetCQHttpDir(), "组件-第三方__Liliya233__频服互通.json"), path.Join(fastbuilder.GetOmegaStorageDir(), "配置", "第三方", "Liliya233", "频服互通", "组件-第三方__Liliya233__频服互通.json")); err != nil {
				panic(err)
			}
			pterm.Success.Println("导入应该成功了")
		}
	}
}
