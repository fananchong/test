package main

import (
	"flag"
	"fmt"
	"runtime"

	cp "github.com/jakecoffman/cp/v2"
)

var num int

func main() {
	flag.IntVar(&num, "num", 0, "num")
	flag.Parse()
	alloc1, totalAlloc1, sys1 := getMem()
	createScene(num)
	alloc2, totalAlloc2, sys2 := getMem()
	fmt.Printf("Allocated memory: %v bytes\n", alloc2-alloc1)
	fmt.Printf("Total memory allocated (including garbage collector overhead): %v bytes\n", totalAlloc2-totalAlloc1)
	fmt.Printf("System memory obtained from the OS: %v bytes\n", sys2-sys1)
}

func createScene(n int) {
	space := cp.NewSpace()

	addBox := func(space *cp.Space, x, y float64) *cp.Body {
		staticBody := cp.NewStaticBody()
		staticBody.SetPosition(cp.Vector{X: x, Y: y})
		body := space.AddBody(staticBody)

		space.AddShape(cp.NewBox(body, 1, 1, 0))
		return body
	}

	for i := 0; i < num; i++ {
		addBox(space, float64(i*10), 0)
	}
}

func getMem() (alloc, totalAlloc, sys uint64) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.Alloc, m.TotalAlloc, m.Sys
}
