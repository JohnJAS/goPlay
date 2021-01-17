package main

import (
	"fmt"
)

type upgradeCheckError struct {
	warning bool
}

func New() error {
	return &upgradeCheckError{
		warning: false,
	}
}

func (err *upgradeCheckError) Error() string {
	return fmt.Sprintf("warning : %v", err.warning)
}

func main() {
	//注意,此处err为error接口变量
	err := New()
	if err != nil {
		//只有接口变量才能进行类型断言
		if err, ok := err.(*upgradeCheckError); ok {
			fmt.Printf("%v\n", err.warning)
		}

		switch err.(type) {
		case *upgradeCheckError:
			fmt.Printf("%T\n", err)
		default:
			fmt.Println("I don't know.")
		}

	}
}
