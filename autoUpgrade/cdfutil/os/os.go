package os

import (
	"bufio"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
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

	return os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0600)
}

//OpenFile
func OpenFile(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0600)
}

//WriteFile
func WriteFile(path string, i interface{}) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0600)
	defer f.Close()
	if err != nil {
		return err
	}

	var s string
	switch i.(type) {
	case int:
		s = strconv.Itoa(i.(int))
	case string:
		s = i.(string)
	}

	w := bufio.NewWriter(f)
	_, err = w.WriteString(s)
	if err != nil {
		return err
	}

	w.Flush()

	return nil
}

//ReadFile
func ReadFile(path string, n ...int) (string, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0600)
	defer f.Close()
	if err != nil {
		return "", err
	}

	switch {
	case len(n) == 0:
		dat, err := ioutil.ReadFile(path)
		return string(dat), err
	case len(n) == 1:
		i := n[0]
		r := bufio.NewReader(f)
		dat, err := r.Peek(i)
		return string(dat), err
	case len(n) < 0 || len(n) > 1:
		return "", errors.New("Internal error in " + reflect.TypeOf(func() {}).PkgPath())
	}

	return "", nil
}
