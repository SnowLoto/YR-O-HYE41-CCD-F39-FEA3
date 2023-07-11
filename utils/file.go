package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func IsDir(path string) bool {
	stat, err := os.Stat(path)
	return (err == nil && stat.IsDir())
}

func IsFile(path string) bool {
	stat, err := os.Stat(path)
	return (err == nil && !stat.IsDir())
}

func MkDir(path string) bool {
	if !IsDir(path) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return false
		}
	}
	return true
}

func GetFileSize(filePath string) (int64, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}
	fileSize := fileInfo.Size()
	return fileSize, nil
}

func GetFileData(filePath string) ([]byte, error) {
	fp, err := os.OpenFile(filePath, os.O_CREATE|os.O_RDONLY, 0755)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	buf, err := io.ReadAll(fp)
	if err != nil {
		return nil, err
	}
	return buf, err
}

func GetJsonData(filePath string, ptr interface{}) error {
	data, err := GetFileData(filePath)
	if err != nil {
		return err
	}
	if len(data) == 0 || data == nil {
		return nil
	}
	err = json.Unmarshal(data, ptr)
	if err != nil {
		return err
	}
	return nil
}

func WriteFileData(filePath string, data []byte) error {
	MkDir(filepath.Dir(filePath))
	fp, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer fp.Close()
	if _, err := fp.Write(data); err != nil {
		return err
	}
	return nil
}

func WriteJsonData(filePath string, data interface{}) error {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	enc := json.NewEncoder(file)
	enc.SetIndent("", "\t")
	enc.SetEscapeHTML(false)
	err = enc.Encode(data)
	if err != nil {
		return err
	}
	return nil
}

func CopyFile(src, dst string) (nBytes int64, err error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}

	defer destination.Close()
	nBytes, err = io.Copy(destination, source)
	return nBytes, err
}

func RemoveFile(src string) bool {
	err := os.Remove(src)
	return err == nil
}
