package main

import (
	"encoding/json"
	"fmt"
)

type I interface {
	F()
}

type aaa struct {
	A int
	C float32
}

type bbb struct {
	A int
	B string
}

func (a aaa) F() {
	fmt.Printf("aaa.A=%d aaa.C=%f\n", a.A, a.C)
}

func (b bbb) F() {
	fmt.Printf("bbb.A=%d bbb.B=%s\n", b.A, b.B)
}

func myNew(t string) I {
	switch t {
	case "aaa":
		return &aaa{10, 99.1}
	case "bbb":
		return &bbb{2, "hello b"}
	}
	return nil
}

func test1() {
	a1 := myNew("aaa")
	data := enc(a1)
	a2 := dec("aaa", data)
	a2.F()
}

func test2() {
	a1 := myNew("aaa")
	data := enc(a1)
	a2 := dec("bbb", data)
	a2.F()
}

func test3() {
	b1 := myNew("bbb")
	data := enc(b1)
	b2 := dec("aaa", data)
	b2.F()
}

func main() {
	test1()
	test2()
	test3()
}

func enc(obj interface{}) []byte {
	buf, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	return buf
}

func dec(t string, data []byte) I {
	i := myNew(t)
	err := json.Unmarshal(data, i)
	if err != nil {
		panic(err)
	}
	return i
}
