package storage

import (
	"fmt"
	"os"
)

func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func Mkdir(path string) {
	if Exists(path) {
		return
	}
	err := os.MkdirAll(path, 0666)
	if err != nil {
		panic(err)
	}
}

func Create(path string) *os.File {
	file, err := os.Create(path)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return file
}

func Delete(path string) {
	err := os.Remove(path)
	if err != nil {
		panic(err)
	}
}
