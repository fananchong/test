// +build !plan9,!windows

package main

/*
#include "./gperftools.h"
*/
import "C"

// SetupGperftools SetupGperftools
func SetupGperftools() {
	C.setup_gperftools()
}

// DumpHeap DumpHeap
func DumpHeap() {
	C.dump_heap()
}

func TestMalloc() {
	C.test_malloc()
}
