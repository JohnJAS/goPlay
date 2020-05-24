package os

import (
	"bufio"
	"os"
	"path/filepath"
)

func check(err error) error {
	if err != nil {
		return err
	}
	return nil
}

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
		err := os.MkdirAll(folder, 0600)
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

	return os.OpenFile(path, os.O_TRUNC, 0600)
}

//OpenFile
func OpenFile(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0600)
}

//WriteFile
func WriteFile(path string, i interface{}) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0600)
	return check(err)

	w := bufio.NewWriter(file)
	_, err = w.WriteString(string(i.(int)))
	return check(err)

	w.Flush()

	return nil
}
