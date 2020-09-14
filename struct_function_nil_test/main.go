package main

import "fmt"

type ii interface {
	F()
}

type aa struct{}

func (a *aa) F() {
	if a == nil {
		fmt.Println("aa is nil")
	} else {
		fmt.Println("aa.F() is called")
	}
}

func test1() {
	var a *aa
	a.F()
}

func test2() {
	var a interface{} = nil
	b, _ := a.(*aa) // 类型转化失败，则 b 为 nil , 类型为 aa
	b.F()
}

func test3() {
	var a ii
	b, _ := a.(*aa) // 类型转化失败，则 b 为 nil , 类型为 aa
	b.F()
}

func test4() {
	var a ii = (*aa)(nil)
	a.(*aa).F()
}

func test5() {
	var a interface{} = (*aa)(nil)
	a.(*aa).F()
}

func main() {
	test1() // aa is nil
	test2() // aa is nil
	test3() // aa is nil
	test4() // aa is nil
	test5() // aa is nil
}
