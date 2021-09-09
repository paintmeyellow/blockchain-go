package main

import (
	"demo/packet"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	booker, _, err := websocket.DefaultDialer.Dial("ws://localhost:8000/ws", nil)
	if err != nil {
		panic(err)
	}
	defer booker.Close()

	go func() {
		for {
			_, m, err := booker.ReadMessage()
			if err != nil {
				fmt.Println(err)
				break
			}
			println(string(m))
		}
	}()

	go func() {
		for {
			<-time.Tick(time.Second)
			err = booker.WriteJSON(&packet.Packet{
				Type:    packet.Request,
				Payload: map[string]interface{}{"handler": "new_address", "keeper_id": "k1"},
			})
			if err != nil {
				log.Println(err)
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
