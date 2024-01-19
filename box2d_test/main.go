package main

import (
	"flag"
	"fmt"
	"runtime"

	"github.com/ByteArena/box2d"
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
	world := box2d.MakeB2World(box2d.B2Vec2{X: 0, Y: -10})

	// 创建一个静态地面
	groundDef := box2d.MakeB2BodyDef()
	groundDef.Position.Set(0, -10)
	ground := world.CreateBody(&groundDef)

	// 创建一个边框作为静态地面
	groundBox := box2d.MakeB2EdgeShape()
	groundBox.Set(box2d.B2Vec2{X: -5000000, Y: 0}, box2d.B2Vec2{X: 5000000, Y: 0})
	ground.CreateFixture(&groundBox, 0)

	// 创建 n 个静态物体
	for i := 0; i < n; i++ {
		// 创建静态物体的定义
		bodyDef := box2d.MakeB2BodyDef()
		bodyDef.Type = box2d.B2BodyType.B2_staticBody
		bodyDef.Position.Set(float64(i), 5)

		// 在世界中创建静态物体
		body := world.CreateBody(&bodyDef)

		// 创建静态物体的形状（这里使用简单的盒子形状）
		shape := box2d.MakeB2PolygonShape()
		shape.SetAsBox(0.5, 0.5)

		// 将形状添加到物体中
		body.CreateFixture(&shape, 0)
	}
}

func getMem() (alloc, totalAlloc, sys uint64) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.Alloc, m.TotalAlloc, m.Sys
}
