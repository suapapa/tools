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

	defer s.Close()

	quitC := make(chan struct{})
	stdinC := make(chan []byte, 1)
	serialC := make(chan []byte, 1)

	ctx, cancle := context.WithCancel(context.Background())

	// serial <- stdin
	go func(ctx context.Context) {
	loop:
		for {
			select {
			case <-ctx.Done():
				log.Println("quit loopR")
				break loop
			case b := <-stdinC:
				s.Write(b)
			}
		}
		quitC <- struct{}{}
	}(ctx)

	// TODO: how cat I put context here?
	go func() {
		b := make([]byte, 1)
		for {
			os.Stdin.Read(b)
			stdinC <- b
		}
	}()

	// stdout <- serial
	go func(ctx context.Context) {
	loop:
		for {
			select {
			case <-ctx.Done():
				log.Println("quit loop")
				break loop
			case b := <-serialC:
				os.Stdout.Write(b)
			}
		}
		quitC <- struct{}{}
	}(ctx)

	// TODO: how cat I put context here?
	go func() {
		b := make([]byte, 1)
		for {
			s.Read(b)
			serialC <- b
		}
	}()

	// ----
	sigC := make(chan os.Signal, 1)
	signal.Notify(sigC, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigC
		cancle()
	}()

	<-quitC
	<-quitC
}
