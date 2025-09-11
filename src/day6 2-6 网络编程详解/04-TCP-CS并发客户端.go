package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 处理优雅退出
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// 主动发起连接请求
	conn, err := net.Dial("tcp", "127.0.0.1:8001")
	if err != nil {
		fmt.Println("net.Dial err:", err)
		return
	}
	defer conn.Close()
	var wg sync.WaitGroup
	wg.Add(1)
	// 获取用户键盘输入( stdin )，将输入数据发送给服务器
	go func() {
		defer wg.Done()
		defer cancel()

		buf := make([]byte, 4096)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				n, err := os.Stdin.Read(buf)
				if err != nil {
					fmt.Println("os.Stdin.Read err:", err)
					return
				}
				//写给服务器, 读多少，写多少！
				if _, err := conn.Write(buf[:n]); err != nil {
					fmt.Println("conn.Write err:", err)
					return
				}
			}
		}
	}()

	// 回显服务器回发的大写数据
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer cancel()
		buf := make([]byte, 4096)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				n, err := conn.Read(buf)
				if n == 0 {
					fmt.Println("检查到服务器关闭，客户端也关闭")
					return
				}
				if err != nil {
					fmt.Println("conn.Read err:", err)
					return
				}
				fmt.Println("客户端读到服务器回发：", string(buf[:n]))
			}

		}

	}()
	// 等待信号或上下文取消
	select {
	case <-sigCh:
		fmt.Println("\n接收到退出信号")
		cancel()
	case <-ctx.Done():
	}

	wg.Wait()
	fmt.Println("客户端退出")
}
