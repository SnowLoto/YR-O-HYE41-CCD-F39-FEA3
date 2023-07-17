package cqhttp

import (
	"bytes"
	"net/url"
	"omega_launcher/fastbuilder"
	"omega_launcher/plantform"
	"omega_launcher/utils"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func GetCQHttpDir() string {
	return filepath.Join(utils.GetCurrentDir(), "cqhttp_storage")
}

func GetCqHttpExec() string {
	cqhttp := filepath.Join(utils.GetCurrentDir(), plantform.GetCQHttpName())
	p, err := filepath.Abs(cqhttp)
	if err != nil {
		panic(err)
	}
	return p
}

func GetCQHttpHash() string {
	return utils.GetFileHash(GetCqHttpExec())
}

func WaitConnect(addr string) {
	u := url.URL{Scheme: "ws", Host: addr}
	for {
		if _, _, err := websocket.DefaultDialer.Dial(u.String(), nil); err == nil {
			return
		}
	}
}

func PackCQHttpRunAuth(qGroupLinkFp, qGuildLinkFp string) {
	_uuid, _ := uuid.NewUUID()
	uuid := _uuid.String()
	uuidFile := filepath.Join(GetCQHttpDir(), "uuid")
	if err := utils.WriteFileData(uuidFile, []byte(uuid)); err != nil {
		panic(err)
	}
	if _, err := utils.CopyFile(qGroupLinkFp, filepath.Join(GetCQHttpDir(), "组件-群服互通.json")); err != nil {
		panic(err)
	}
	if _, err := utils.CopyFile(qGuildLinkFp, filepath.Join(GetCQHttpDir(), "组件-频服互通.json")); err != nil {
		panic(err)
	}
	fileName := filepath.Join(fastbuilder.GetOmegaStorageDir(), "上传这个文件到云服务器以使用云服务器的群服互通.data")
	fp, err := os.OpenFile(fileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0755)
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	fp.Write([]byte(uuid))
	if err := utils.Zip(GetCQHttpDir(), fp, []string{"data", "logs"}); err != nil {
		panic(err)
	}
	os.Remove(filepath.Join(GetCQHttpDir(), "组件-群服互通.json"))
	os.Remove(filepath.Join(GetCQHttpDir(), "组件-频服互通.json"))
}

func UnPackCQHttpRunAuth() bool {
	fileName := filepath.Join(fastbuilder.GetOmegaStorageDir(), "上传这个文件到云服务器以使用云服务器的群服互通.data")
	fp, err := os.OpenFile(fileName, os.O_RDONLY, 0)
	if err != nil {
		return false
	}
	uuidBytes := make([]byte, 36)
	if _, err := fp.Read(uuidBytes); err != nil {
		panic(err)
	}
	uuidFile := filepath.Join(GetCQHttpDir(), "uuid")
	if thisUUidBytes, err := utils.GetFileData(uuidFile); err == nil {
		if bytes.Equal(thisUUidBytes, uuidBytes) {
			return false
		}
	}
	if utils.GetInputYN("已读取到 data 文件，要导入吗?") {
		os.RemoveAll(GetCQHttpDir())
		if err := utils.UnZip(fp, GetCQHttpDir()); err != nil {
			panic(err)
		}
		utils.CopyFile(filepath.Join(GetCQHttpDir(), "组件-群服互通.json"), filepath.Join(fastbuilder.GetOmegaStorageDir(), "配置", "群服互通", "组件-群服互通.json"))
		utils.CopyFile(filepath.Join(GetCQHttpDir(), "组件-频服互通.json"), filepath.Join(fastbuilder.GetOmegaStorageDir(), "配置", "第三方_by_Liliya233", "频服互通", "组件-第三方__Liliya233__频服互通.json"))
		return true
	}
	return false
}
