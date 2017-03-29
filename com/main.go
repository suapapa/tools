package main

import (
	"context"
	"fmt"
	"os"

	"github.com/tarm/serial"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: TODO")
	}

	pn := os.Args[1]
	s, err := serial.OpenPort(&serial.Config{
		Name: pn,
		Baud: 115200,
	})
	if err != nil {
		panic(err)
	}

	rBuff := make([]byte, 1)
	wBuff := make([]byte, 1)

	quitC := make(chan struct{})
	ctx, cancle := context.WithCancel(context.Background())
	defer cancle()

	go func(ctx context.Context) {
	loopR:
		for {
			select {
			case <-ctx.Done():
				break loopR
			default:
				s.Read(rBuff)
				os.Stdout.Write(rBuff)
			}
		}
		quitC <- struct{}{}
	}(ctx)

	go func(ctx context.Context) {
	loopW:
		for {
			select {
			case <-ctx.Done():
				break loopW
			default:
				os.Stdin.Read(wBuff)
				s.Write(wBuff)
			}
		}
		quitC <- struct{}{}
	}(ctx)

	<-quitC
	<-quitC
}
