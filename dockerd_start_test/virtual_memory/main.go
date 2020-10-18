package main

import (
	"time"
	"unsafe"
)

var m1 []byte
var m2 []byte
var m3 []byte
var m4 []byte

func main() {
	m1 = make([]byte, 4*1024*1024*1024)
	m2 = make([]byte, 4*1024*1024*1024)
	m3 = make([]byte, 4*1024*1024*1024)
	m4 = make([]byte, 4*1024*1024*1024)
	m1[0] = 'a'
	memset(unsafe.Pointer(&m1[0]), 'c', 1*1024*1024*1024)
	m1[len(m1)-1] = 'b'
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
