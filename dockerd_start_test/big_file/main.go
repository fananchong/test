package main

import (
	"fmt"
	"io/ioutil"
	"time"
	"unsafe"
)

func main() {
	m1 := make([]byte, 4*1024*1024*1024)
	m1[0] = 'a'
	memset(unsafe.Pointer(&m1[0]), 'c', 1*1024*1024*1024)
	m1[len(m1)-1] = 'b'
	for i := 0; i < 5; i++ {
		err := ioutil.WriteFile(fmt.Sprintf("./output%d.txt", i), m1, 0666)
		if err != nil {
			panic(err)
		}
	}
	for {
		time.Sleep(1 * time.Minute)
	}
}

func memset(s unsafe.Pointer, c byte, n uintptr) {
	ptr := uintptr(s)
	var i uintptr
	for i = 0; i < n; i++ {
		pByte := (*byte)(unsafe.Pointer(ptr + i))
		*pByte = c
	}
}
