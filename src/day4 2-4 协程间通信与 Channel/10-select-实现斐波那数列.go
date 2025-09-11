package main

import (
	"fmt"
	"runtime"
)

func fibonacci01(ch <-chan int, quit <-chan bool) int {
	for {
		select {
		case num := <-ch:
			fmt.Print(num, " ")
		case <-quit:
			//return
			runtime.Goexit() //等效于 return
		}
	}
}

func main01() {
	ch := make(chan int)
	quit := make(chan bool)

	go fibonacci01(ch, quit) // 子go 程 打印fibonacci数列

	x, y := 1, 1
	for i := 0; i < 40; i++ {
		ch <- x
		x, y = y, y+x
	}
	quit <- true
}

func fibonacci(n int) int {
	if n <= 1 {
		return n
	}
	a, b := 0, 1
	for i := 2; i <= n; i++ {
		var temp = a
		a = b
		b = temp + b
		//a, b = b, a+b
		fmt.Printf("a:%d b:%d\n", a, b)
	}
	return b
}

func main() {
	n := 10
	result := fibonacci(n)
	fmt.Printf("第 %d 个斐波那契数是: %d\n", n, result)
}
