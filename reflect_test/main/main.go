//go:generate echo "hello" > 1.txt
package main

import (
	"fmt"
	"reflect"
)

var methods = map[string][]string{
	"F": {"p1", "p2"},
}

type A struct {
}

func (a *A) F(p1 int, p2 string) {
	fmt.Println("A")
}

func main() {

	t1 := reflect.TypeOf(&A{}).Elem()
	a := reflect.New(t1)

	// call
	f1 := a.MethodByName("F")
	f1.Call([]reflect.Value{
		reflect.ValueOf(1),
		reflect.ValueOf("1"),
	})
	fmt.Println(f1.Type)

	f2, _ := reflect.TypeOf(&A{}).MethodByName("F")
	fmt.Println(f2.Type.NumIn())
	// reflect.New(f2.Type.In(i)).Elem()

}
