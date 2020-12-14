package main

import (
	"context"
	"fmt"
	"sync"
)

func main() {
	ch := make(chan struct{})

	//fmt.Println(<-ch)
	close(ch)

	fmt.Println(<-ch)
	fmt.Println(<-ch)

	ctx := context.Background()

	valuectx := context.WithValue(ctx, "test", "abc")

	fmt.Println(valuectx.Value("test"))

	valuectx2 := context.WithValue(ctx, "test2", "bcd")

	fmt.Println(valuectx2.Value("test2"))
	fmt.Println(valuectx2.Value("test"))

	cctx, cctxFunc := context.WithCancel(ctx)

	cctx.Done()

}



func parentCancelCtx(parent context.Context) (*context.cancelCtx, bool) {
	for {
		switch c := parent.(type) {
		case *context.cancelCtx:
			return c, true
		case *context.timerCtx:
			return &c.cancelCtx, true
		case *context.valueCtx:
			parent = c.Context
		default:
			return nil, false
		}
	}
}
