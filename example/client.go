package main

import (
	"demo/wsconn"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	wc, err := wsconn.Connect("ws://localhost:8000/ws")
	if err != nil {
		panic(err)
	}
	defer wc.Close()

	for i := 0; i < 10; i++ {
		time.Sleep(time.Millisecond)
		go func() {
			m, err := wc.Request("dc", []byte("Hello from client, DC!"), 2*time.Second)
			if err != nil {
				log.Println(err)
			}
			if m != nil {
				fmt.Println("CLIENT:recv", "id:", m.ID, "reply:", m.Reply, string(m.Data))
			}
		}()
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	select {
	case x := <-interrupt:
		log.Println("received a signal", x.String())
	}
}
