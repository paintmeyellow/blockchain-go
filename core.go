package main

import (
	"demo/packet"
	"fmt"
	"github.com/gorilla/websocket"
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
			address := string(m)
			println(address)
		}
	}()

	for {
		<-time.Tick(3 * time.Second)
		err = booker.WriteJSON(&packet.Transport{KeeperID: "k1", Handler:  "new_address"})
		if err != nil {
			fmt.Println(err)
			break
		}
	}
}
