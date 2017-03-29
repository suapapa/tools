package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/tarm/serial"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: TODO")
	}

	pn := os.Args[1]
	s, err := serial.OpenPort(&serial.Config{
		Name: pn,
		Baud: 9600,
	})
	if err != nil {
		panic(err)
	}

	quitC := make(chan struct{})
	ctx, cancle := context.WithCancel(context.Background())

	go func(ctx context.Context) {
		rBuff := make([]byte, 1)

	loopR:
		for {
			select {
			case <-ctx.Done():
				log.Println("quit loopR")
				break loopR
			default:
				s.Read(rBuff)
				os.Stdout.Write(rBuff)
			}
		}
		quitC <- struct{}{}
	}(ctx)

	go func(ctx context.Context) {
		wBuff := make([]byte, 1)

	loopW:
		for {
			select {
			case <-ctx.Done():
				log.Println("quit loopW")
				break loopW
			default:
				os.Stdin.Read(wBuff)
				s.Write(wBuff)
			}
		}
		quitC <- struct{}{}
	}(ctx)

	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigC
		cancle()
	}()

	<-quitC
	<-quitC
}
