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
	err := New()
	if err != nil {
		if err, ok := err.(*upgradeCheckError); ok {
			fmt.Printf("%v", err.warning)
		}
	}
}
