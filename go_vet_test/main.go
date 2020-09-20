package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func test1() {
	str := "hello world!"
	fmt.Printf("%d\n", str)
}

func test2() {
	str := "hello world!"
	fmt.Printf("%s\n", &str)
}

func customLogf(str string, args ...interface{}) {
	fmt.Printf(str, args...)
}

func test3() {
	i := 42
	customLogf("the answer is %s\n", i)
}

func test4() {
	var i int
	fmt.Println(i != 0 || i != 1) // always true
	fmt.Println(i == 0 && i == 1) // always false
	fmt.Println(i == 0 && i == 0) // redundant check
}

func test5() {
	words := []string{"foo", "bar", "baz"}
	for _, word := range words {
		fmt.Println("in for", word)
		go func() {
			fmt.Println("in goroutine", word) // 3 prints are all `in goroutine baz`
		}()
	}
	time.Sleep(1 * time.Second)
}

func test6() {

	defer func() {
		r := recover()
		log.Println(r)
	}()

	res, err := http.Get("https://www.xxxxx.io/")
	defer res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func test7() {
	i := 42
	i = i
}

func test8() {
	fmt.Println(test7 == nil)
}

func test9() {
	i := 42
	fmt.Println("test9:", i>>64)
}

func main() {
	test1()
	test2()
	test3()
	test4()
	test5()
	test6()
	test7()
	test8()
	test9()
}
