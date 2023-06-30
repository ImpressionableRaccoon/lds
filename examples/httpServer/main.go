package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/ImpressionableRaccoon/lds"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cancel()

	var (
		serverAddress = flag.String("a", ":15848", "http server address")
		serialPort    = flag.String("p", "", "serial port name")
		baudRate      = flag.Int("b", 230400, "baud rate")
	)
	flag.Parse()

	if *serverAddress == "" || *serialPort == "" {
		flag.PrintDefaults()
		return
	}

	l, err := lds.New(lds.Config{
		PortName: *serialPort,
		BaudRate: *baudRate,
	})
	if err != nil {
		panic(err)
	}
	defer func() {
		err := l.Close()
		if err != nil {
			fmt.Printf("lidar close error: %s\n", err)
		}
	}()

	go func() {
		err := l.Worker(ctx)
		if err != nil {
			fmt.Printf("worker stopped: %s\n", err)
		}
	}()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		data, err := json.Marshal(l.Get())
		if err != nil {
			http.Error(w, fmt.Sprintf("parse data error: %s", err), http.StatusInternalServerError)
			return
		}
		_, err = w.Write(data)
		if err != nil {
			fmt.Printf("http write error: %s\n", err)
		}
	})
	go func() {
		srv := http.Server{
			ReadHeaderTimeout: time.Second,
		}

		ln, err := net.Listen("tcp", *serverAddress)
		if err != nil {
			fmt.Printf("listen http server address error: %s\n", err)
			return
		}

		err = srv.Serve(ln)
		if err != nil {
			fmt.Printf("listed and serve error: %s\n", err)
		}
	}()

	<-ctx.Done()
}
