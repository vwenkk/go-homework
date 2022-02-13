package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

//  基于 errgroup 实现一个 http server 的启动和关闭 ，以及 linux signal 信号的注册和处理，要保证能够一个退出，全部注销退出。
//
func main() {
	var (
		quit = make(chan struct{})
	)

	group, ctx := errgroup.WithContext(context.Background())
	server1 := http.Server{
		Addr: ":8081",
		Handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			fmt.Fprintln(writer, "hello, server1")
		}),
	}

	server2 := http.Server{
		Addr: ":8082",
		Handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			fmt.Fprintln(writer, "hello, server2")
		}),
	}
	// 启动第一个server
	group.Go(func() error {
		go func() {
			select {
			case <-quit:
				log.Printf("server1, 开始退出")
				time.Sleep(2 * time.Second)
				server1.Shutdown(ctx)
			}
		}()
		return server1.ListenAndServe()
	})

	// 启动第二个server
	group.Go(func() error {
		go func() {
			select {
			case <-quit:
				log.Printf("server2, 开始退出")
				time.Sleep(1 * time.Second)
				server2.Shutdown(ctx)
			}
		}()
		return server2.ListenAndServe()
	})

	// 监听退出信号
	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

		select {
		case <-ctx.Done():
			log.Printf("err:%v", ctx.Err())
			close(quit)
		case <-sigs:
			close(quit)
			log.Printf("收到信号进行退出...")
		}
		// 三秒未退出，强制退出 防止长时间卡死
		time.AfterFunc(3*time.Second, func() {
			log.Fatalf("长时间未退出, 强制退出")
		})
	}()

	group.Wait()
	fmt.Print("退出成功")
}
