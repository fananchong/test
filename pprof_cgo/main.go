package main

/*
#include "./gperftools.h"
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
	DumpHeap()
}
