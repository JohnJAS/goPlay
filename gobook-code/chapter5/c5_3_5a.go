/*
代码所在章节：5.3.5节
*/

package main

import (
	"context"
	"fmt"
	"time"
)

//define a new type include a Context Field
type otherContext struct {
	context.Context
}

func main() {

	//Construct a *cancelCtx type object
	ctxa, cancel := context.WithCancel(context.Background())
	/*
			ctxa内部状态 --> ctxa=&cancelCtx{
								Context: new(emptyCtx),
		                   }
	*/
	go work(ctxa, "work1")

	//Construct a *timerCtx type object wrapped by *cancelCtx
	tm := time.Now().Add(3 * time.Second)
	//ctxa is child thread of ctxa
	ctxb, _ := context.WithDeadline(ctxa, tm)
	/*
			ctxb内部状态 --> ctxb=&timeCtx{
								cancelCtx: ctxa,
		                        dataline:tm,
							}
			同时触发ctxa ，在children 中维护ctxb 作为子节点
	*/
	go work(ctxb, "work2")

	oc := otherContext{ctxb}
	//Construct a *cancelCtx type object wrapped by oc
	ctxc := context.WithValue(oc, "key", "god andes,pass from main ")
	/*
		ctxc -->  ctxc=&cancelCtx {
						Context: oc,
					}
		同时通过oc.C ontext 找到ctxb ，通过ctxb . cancelCtx 找到ctxa ，在ctxa 的children
		字段中维护ctxc 作为其子节点
	*/
	go workWithValue(ctxc, "work3")

	time.Sleep(10 * time.Second)
	cancel()
	time.Sleep(5 * time.Second)
	fmt.Println("main stop")
}

//do something
func work(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("%s get msg to cancel\n", name)
			return
		default:
			fmt.Printf("%s is running \n", name)
			time.Sleep(1 * time.Second)
		}
	}
}

//do something and pass values by context
func workWithValue(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("%s get msg to cancel\n", name)
			return
		default:
			value := ctx.Value("key").(string)
			fmt.Printf("%s is running value=%s \n", name, value)
			time.Sleep(1 * time.Second)
		}
	}
}
