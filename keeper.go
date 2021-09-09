package main

import (
	"demo/packet"
	"encoding/json"
	"errors"
	"github.com/gorilla/websocket"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	booker, _, err := websocket.DefaultDialer.Dial("ws://localhost:8000/ws", nil)
	if err != nil {
		panic(err)
	}
	defer booker.Close()
	p := packet.Packet{
		Type:    packet.ConnectionInit,
		Payload: map[string]interface{}{"keeper_id": "k1"},
	}
	if err = booker.WriteJSON(&p); err != nil {
		panic(err)
	}
	println("connected")

	go func() {
		for {
			_, r, err := booker.NextReader()
			if err != nil {
				log.Println(err)
				break
			}
			var p packet.Packet
			if err := json.NewDecoder(r).Decode(&p); err != nil {
				log.Println(err)
				continue
			}
			switch p.Type {
			case packet.Request:
				if err := handleRequestKeeper(booker, &p); err != nil {
					log.Println(err)
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

func handleRequestKeeper(c *websocket.Conn, p *packet.Packet) error {
	handler, err := p.Payload.StringValue("handler")
	if err != nil {
		return err
	}
	switch handler {
	case "new_address":
		kp := packet.Packet{
			ID:   p.ID,
			Type: packet.Next,
			Payload: map[string]interface{}{
				"address":     "bc1qxy2kgdygjrsqtzq2n0yrf2493p83kkfjhx0wlh",
				"private_key": "5Kb8kLf9zgWQnogidDA76MzPL6TsZZY36hWXMssSzNydYXYB9KF",
			},
		}
		return c.WriteJSON(&kp)
	default:
		return errors.New("handler is not found")
	}
}
