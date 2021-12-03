package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"time"
)

func main() {
	err := AskForConfirm("")
	fmt.Println(err)

	err = AskForConfirm("")
	fmt.Println(err)
}


func AskForConfirm(msg string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	inputChan := make(chan string, 1)
	defer close(inputChan)
	received := make(chan struct{})
	defer close(received)
	reput := make(chan struct{})
	defer close(reput)
	go readUserInput(ctx, msg, inputChan, received, reput)
Loop:
	for {
		select {
		case <-ctx.Done():
			err = ctx.Err()
			break Loop
		case i := <-inputChan:
			fmt.Printf("User input: %s\n", i)
			if i == "y" || i == "yes" || i == "Y" || i == "Yes" || i == "YES" {
				received <- struct{}{}
				break Loop
			} else if i == "n" || i == "no" || i == "N" || i == "No" || i == "NO" {
				received <- struct{}{}
				cancel()
				err = ctx.Err()
			} else {
				fmt.Println("Please Input Y/N :")
				reput <- struct{}{}
				continue
			}
		default:
			time.Sleep(time.Second)
		}
	}
	return
}

func readUserInput(ctx context.Context, message string, inputChan chan string, received, reput chan struct{}) {
	cc, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	f := bufio.NewReader(os.Stdin)
	var input string
	if message == "" {
		message = "Please check the warning messages and confirm to continue: (Y/N)"
	}
	fmt.Println(message)
	for {
		str, _ := f.ReadString('\n')
		fmt.Sscan(str, &input)
		inputChan <- input
		select {
		case <-cc.Done():
			return
		case <-received:
			return
		case <-reput:
			continue
		}
	}
}
