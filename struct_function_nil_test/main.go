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
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	var a interface{} = nil
	a.(*aa).F()
}

func test3() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	var a ii
	a.(*aa).F()
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
	test2() // interface conversion: interface {} is nil, not *main.aa
	test3() // interface conversion: main.ii is nil, not *main.aa
	test4() // aa is nil
	test5() // aa is nil
}
