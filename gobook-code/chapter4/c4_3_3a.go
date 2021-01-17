/*
代码所在章节：4.3.3节
*/

package main

import "fmt"

type Inter interface {
	Ping()
	Pang()
}

type St struct{}

func (St) Ping() {
	fmt.Println("ping")
}
func (*St) Pang() {
	fmt.Println("pang")
}

func main() {
	var st *St = nil
	var it Inter = st

	//结果0x0表示空指针
	fmt.Printf("%T %p %v\n", st, st, st)
	fmt.Printf("%T %p %v\n", it, it, it)

	if st == it {
		fmt.Println(true)
	}

	if st == nil {
		//此处st确实为nil
		fmt.Println("st is nil")
	}

	//但这里it不为nil,是因为空接口有2个字段，一个是实例类型，一个是指向绑定实例的指针
	if it == nil {
		fmt.Println("it is nil")
	} else {
		it.Pang()
		//下面的语句会导致panic
		//panic: value method main.St.Ping called using nil *St pointer
		//方法转换为值接受者函数调用，第一个参数是St类型，由于*St是nil，无法获取指针所指的对象值，所以panic.
		//it.Ping()
	}
}

//*main.St 0x0 <nil>
//*main.St 0x0 <nil>
//true
//st is nil
//pang
