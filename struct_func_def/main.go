package main

import (
	"fmt"
)

type sss struct {
	A int
}

func (s sss) F1() {
	s.A++
	fmt.Printf("F1, s addr: %p. s.A=%d\n", &s, s.A)
}

func (s *sss) F2() {
	s.A++
	fmt.Printf("F2, s addr: %p. s.A=%d\n", s, s.A)
}

func main() {
	var s sss
	fmt.Printf("s addr: %p. s.A=%d\n", &s, s.A)
	s.F1()
	s.F1()
	s.F1()
	s.F2()
	s.F2()
	s.F2()
}
