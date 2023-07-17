package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/andybalholm/brotli"
)

func GetFileData(fname string) ([]byte, error) {
	fp, err := os.OpenFile(fname, os.O_CREATE|os.O_RDONLY, 0755)
	if err != nil {
		return nil, err
	}
	defer fp.Close()
	buf, err := ioutil.ReadAll(fp)
	if err != nil {
		return nil, err
	}
	return buf, err
}

func WriteFileData(fname string, data []byte) error {
	fp, err := os.OpenFile(fname, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer fp.Close()
	if _, err := fp.Write(data); err != nil {
		return err
	}
	return nil
}

func main() {
	var cqList []string = []string{
		"./build/go-cqhttp_darwin_amd64",
		"./build/go-cqhttp_darwin_arm64",
		"./build/go-cqhttp_linux_amd64",
		"./build/go-cqhttp_linux_arm64",
		"./build/go-cqhttp_windows_amd64.exe",
		"./build/go-cqhttp_windows_arm64.exe",
		"./build/go-cqhttp_android_amd64",
		"./build/go-cqhttp_android_arm64",
	}

	for _, inFile := range cqList {
		outFile := inFile + ".brotli"
		fmt.Printf("comprerss: %v -> %v \n", inFile, outFile)
		var origData []byte
		var err error
		origData, err = GetFileData(inFile)
		if err != nil || len(origData) == 0 {
			fmt.Printf("read %v fail\n", inFile)
			continue
		}

		buf := bytes.NewBuffer([]byte{})
		compressor := brotli.NewWriterLevel(buf, brotli.DefaultCompression)
		compressor.Write(origData)
		compressor.Close()
		newData := buf.Bytes()

		if err := WriteFileData(outFile, newData); err != nil {
			panic(err)
		}

		fmt.Printf(" compress %.3f\n", float32(len(newData))/float32(len(origData)))
		os.Remove(inFile)
	}

	// var compressedData []byte
	// if compressedData, err = GetFileData(*outFile); err != nil {
	// 	panic(err)
	// }

	// if recoveredData, err := ioutil.ReadAll(brotli.NewReader(bytes.NewReader(compressedData))); err != nil {
	// 	panic(err)
	// } else {
	// 	if bytes.Compare(recoveredData, origData) != 0 {
	// 		panic("not same")
	// 	}
	// 	fmt.Println("Success")
	// }
}
