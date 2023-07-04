package utils

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func Zip(src_dir string, zipfile *os.File, ignores []string) error {
	archive := zip.NewWriter(zipfile)
	defer archive.Close()
	return filepath.Walk(src_dir, func(filePath string, info os.FileInfo, _ error) error {
		if filePath == src_dir {
			return nil
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = filePath[len(src_dir)+1:]
		header.Name = strings.ReplaceAll(header.Name, "\\", "/")
		for _, ignore := range ignores {
			if strings.HasPrefix(header.Name, ignore) {
				return nil
			}
		}
		if info.IsDir() {
			return nil
		} else {
			// 设置：zip的文件压缩算法
			header.Method = zip.Deflate
		}
		// 创建：压缩包头部信息
		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			file, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer file.Close()
			if _, err := io.Copy(writer, file); err != nil {
				return err
			}
		}
		return nil
	})
}

func UnZip(zipfile []byte, dstDir string) error {
	// 创建zip文件的Reader
	zipReader, err := zip.NewReader(bytes.NewReader(zipfile), int64(len(zipfile)))
	if err != nil {
		return err
	}
	// 遍历ZIP文件中的每个文件/目录
	for _, file := range zipReader.File {
		// 获取解压后的文件路径
		extractedFilePath := filepath.Join(dstDir, file.Name)
		// 如果是目录，创建目录
		if file.FileInfo().IsDir() {
			err = os.MkdirAll(extractedFilePath, file.Mode())
			if err != nil {
				return err
			}
			continue
		}
		// 创建解压后的文件所在的目录
		err = os.MkdirAll(filepath.Dir(extractedFilePath), file.Mode())
		if err != nil {
			return err
		}
		// 创建解压后的文件
		extractedFile, err := os.Create(extractedFilePath)
		if err != nil {
			return err
		}
		defer extractedFile.Close()
		// 打开ZIP文件中的文件
		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()
		// 将ZIP文件中的内容复制到解压后的文件
		_, err = io.Copy(extractedFile, fileReader)
		if err != nil {
			return err
		}
	}
	return nil
}

func UnTarGz(tarGzData []byte, dstDir string) error {
	// 创建gzip.Reader来读取压缩数据
	gzipReader, err := gzip.NewReader(bytes.NewReader(tarGzData))
	if err != nil {
		return err
	}
	defer gzipReader.Close()
	// 创建tar.Reader来读取.tar数据
	tarReader := tar.NewReader(gzipReader)
	// 遍历tar文件中的每个文件/目录
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break // 读取结束
		}
		if err != nil {
			return err
		}
		// 获取解压后的文件路径
		extractedFilePath := filepath.Join(dstDir, header.Name)
		// 根据文件类型进行相应的操作
		switch header.Typeflag {
		case tar.TypeDir: // 目录
			err = os.MkdirAll(extractedFilePath, os.ModePerm)
		case tar.TypeReg: // 文件
			// 创建解压后的文件所在的目录
			err = os.MkdirAll(filepath.Dir(extractedFilePath), os.ModePerm)
			if err != nil {
				return err
			}
			// 创建解压后的文件
			extractedFile, err := os.Create(extractedFilePath)
			if err != nil {
				return err
			}
			defer extractedFile.Close()
			// 将tar文件中的内容复制到解压后的文件
			_, err = io.Copy(extractedFile, tarReader)
			if err != nil {
				return err
			}
		}
		if err != nil {
			return err
		}
	}
	return nil
}
