/*
代码所在章节：4.2.1节
*/

package main

import "fmt"

type Inter interface {
	Ping()
	Pang()
}

type Anter interface {
	Inter
	String()
}

type St struct {
	Name string
}

func (St) Ping() {
	fmt.Println("ping")
}
func (*St) Pang() {
	fmt.Println("pang")
}

func main() {
	st := &St{"andes"}
	var i interface{} = st

	//判断i绑定的实例是否实现了接口类型Inter
	o := i.(Inter)
	o.Ping()
	o.Pang()

	//不加类型断言会报错
	if p, ok := i.(Anter); ok {
		p.String()
	}

	//判断i绑定的实例是否就是具体类型St
	s := i.(*St)
	fmt.Printf("%s", s.Name)

}
