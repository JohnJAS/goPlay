package os

import (
	"os"
	"path/filepath"
)

//PathExists return whether file or folder exists or not
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//CreateFolder
func CreateFolder(path string) error {
	folder := filepath.Dir(path)
	exist, _ := PathExists(folder)
	if exist == false {
		err := os.MkdirAll(folder, 0666)
		if err != nil {
			return err
		}
	}
	return nil
}

//CreateFile
func CreateFile(path string) (*os.File, error) {
	folder := filepath.Dir(path)
	exist, _ := PathExists(folder)
	if exist == false {
		err := os.MkdirAll(folder, 0666)
		if err != nil {
			return nil, err
		}
	}

	return os.Create(path)
}

//OpenFile
func OpenFile(path string) (*os.File, error) {
	return os.OpenFile(path,os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
}