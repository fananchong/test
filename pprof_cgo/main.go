package main

/*
#include "./gperftools.h"
#include <stdlib.h>
#include <stdio.h>
#cgo LDFLAGS: -L. -L./lib/lib/ -lgperftools -ltcmalloc -lstdc++ -lpthread -lm
*/
import "C"
import (
	_ "fmt"
)

func main() {
	SetupGperftools()
	TestMalloc()
	TestMalloc()
	TestMalloc()
	ptr1 := C.malloc(C.size_t(8 * 1024 * 1024))
	ptr2 := C.malloc(C.size_t(8 * 1024 * 1024))
	C.malloc(C.size_t(8 * 1024 * 1024))
	C.free(ptr1)
	C.free(ptr2)
	DumpHeap()
}
