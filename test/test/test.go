package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os/exec"
	"syscall"
	"time"
)

type CMD struct {
	name    string
	args    []string
	stdout  io.ReadCloser
	stderr  io.ReadCloser
	print   bool
	timeout int
	cancel  context.CancelFunc

	buff []byte
}

func Shell(script string) *CMD {
	return Command("/bin/bash", "-c", script)
}

func Command(name string, args ...string) *CMD {
	return &CMD{name: name, args: args}
}

func (cmd *CMD) readPrint() ([]byte, error) {
	extend := func(n int) {
		tmp := make([]byte, cap(cmd.buff)+n)
		copy(tmp, cmd.buff)
		cmd.buff = tmp[:len(cmd.buff)]
	}
	read := func() (b []byte, err error) {
		n := 1024
		rear := len(cmd.buff)
		if rear+n > cap(cmd.buff) {
			extend(1024 + n)
		}
		n, err = cmd.stdout.Read(cmd.buff[rear : rear+n])
		cmd.buff = cmd.buff[:rear+n]
		return cmd.buff[rear : rear+n], err
	}

	for {
		b, err := read()
		fmt.Printf("%s", b)
		if err != nil {
			if err.Error() == "EOF" {
				return cmd.buff, nil
			} else {
				return nil, err
			}
		}
	}
}

func (cmd *CMD) Print2Console() *CMD {
	cmd.print = true
	return cmd
}

func (cmd *CMD) SetTimeout(t int) *CMD {
	cmd.timeout = t
	return cmd
}

func (cmd *CMD) Cancel() {
	if cmd.cancel != nil {
		cmd.cancel()
	}
}

func (cmd *CMD) buildCtxt() (ctx context.Context, command *exec.Cmd) {
	if cmd.timeout > 0 {
		ctx, cmd.cancel = context.WithTimeout(context.Background(), time.Duration(cmd.timeout)*time.Second)
	} else {
		ctx, cmd.cancel = context.WithCancel(context.Background())
	}
	command = exec.CommandContext(ctx, cmd.name, cmd.args...)
	return
}

func (cmd *CMD) Run() (stdout, stderr []byte, err error, retCode int) {
	ctx, command := cmd.buildCtxt()

	defer func() {
		if err != nil {
			retCode = -1000
			log.Printf("Run CMD err: %s(%+v)", err, ctx.Err())
		}
		err = command.Wait()
		if err != nil {
			// unknown error code
			retCode = -1
			if e, ok := err.(*exec.ExitError); ok {
				if status, ok := e.Sys().(syscall.WaitStatus); ok {
					retCode = status.ExitStatus()
				}
				if ctxErr := ctx.Err(); ctxErr != nil {
					err = ctxErr
					switch ctxErr {
					case context.Canceled:
						log.Printf("cmd [%s] canceled", cmd.name)
					case context.DeadlineExceeded:
						log.Printf("run cmd [%s] timeout", cmd.name)
					}
				}
			}
		}
		cmd.cancel()
	}()

	cmd.stdout, err = command.StdoutPipe()
	if err != nil {
		return
	}
	cmd.stderr, err = command.StderrPipe()
	if err != nil {
		return
	}
	log.Printf("Run CMD: %s %v", cmd.name, cmd.args)
	err = command.Start()
	if err != nil {
		return
	}
	//pid := command.Process.Pid
	if cmd.print {
		stdout, err = cmd.readPrint()
	} else {
		stdout, err = ioutil.ReadAll(cmd.stdout)
	}
	if err != nil {
		return
	}
	stderr, err = ioutil.ReadAll(cmd.stderr)
	if err != nil {
		return
	}
	if cmd.print {
		fmt.Printf("%s", stderr)
	}
	return
}

func main(){
	_,_,_,_ = Command("hostname").Run()
	//commit3
	//commit2
}