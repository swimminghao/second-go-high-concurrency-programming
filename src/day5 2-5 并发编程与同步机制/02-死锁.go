package main

import (
	"fmt"
	"time"
)

// 死锁1
/*func main()  {
	ch :=make(chan int)
	ch <- 789
	num := <-ch
	fmt.Println("num = ", num)
}*/

// 死锁2
/*func main()  {
	ch := make(chan int)
	go func() {
		ch <- 789
	}()
	num := <- ch
	fmt.Println("num = ", num)
}*/

// 死锁 3
func main01() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	go func() { // 子
		for {
			select {
			case num := <-ch1:
				ch2 <- num
			}
		}
	}()
	for {
		select {
		case num := <-ch2:
			ch1 <- num
		}
	}
}

func main02() {
	ch := make(chan int)

	// 启动发送方的 Goroutine
	go func() {
		ch <- 1
		fmt.Println("send")
	}()

	// 主 Goroutine 稍作等待，确保发送方 Goroutine 已经运行并阻塞在发送操作上
	time.Sleep(time.Millisecond * 100)

	// 主 Goroutine 充当接收者
	num := <-ch
	fmt.Println("received:", num)

	fmt.Println("over")
}

func main() {
	ch := make(chan int)
	ch <- 1 // I'm blocked because there is no channel read yet.
	fmt.Println("send")
	go func() {
		<-ch // I will never be called for the main routine is blocked!
		fmt.Println("received")
	}()
	fmt.Println("over")
}
