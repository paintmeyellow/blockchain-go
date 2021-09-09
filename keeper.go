package main

import (
	"demo/packet"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	booker, _, err := websocket.DefaultDialer.Dial("ws://localhost:8000/connect", nil)
	if err != nil {
		panic(err)
	}
	defer booker.Close()

	err = booker.WriteJSON(&packet.Connection{KeeperID: "k1"})
	if err != nil {
		panic(err)
	}
	println("connected")

	go func() {
		for {
			_, m, err := booker.ReadMessage()
			if err != nil {
				log.Println(err)
				break
			}
			handler := string(m)
			if handler == "new_address" {
				err = booker.WriteMessage(websocket.TextMessage, []byte("bc1qxy2kgdygjrsqtzq2n0yrf2493p83kkfjhx0wlh"))
				if err != nil {
					fmt.Println(err)
					break
				}
			}
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	select {
	case x := <-interrupt:
		log.Println("received a signal", x.String())
	}
}
