package main

import (
	"demo/packet"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

var upgrader = websocket.Upgrader{}

var requests map[uint64]*websocket.Conn

var keepers map[string]*websocket.Conn

var requestN uint64

func main() {
	requests = make(map[uint64]*websocket.Conn)
	keepers = make(map[string]*websocket.Conn)
	go checkConnections()
	http.HandleFunc("/ws", serveWS)
	log.Fatalln(http.ListenAndServe(":8000", nil))
}

func checkConnections() {
	for {
		<-time.Tick(4 * time.Second)
		fmt.Println("requests", requests)
		fmt.Println("keepers", keepers)
		fmt.Println()
	}
}

func serveWS(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		println(err.Error())
	}
	var keeperID string
	defer func() {
		c.Close()
		if keeperID != "" {
			delete(keepers, keeperID)
		}
	}()
	for {
		_, reader, err := c.NextReader()
		if err != nil {
			log.Println(err)
			break
		}
		var p packet.Packet
		if err := json.NewDecoder(reader).Decode(&p); err != nil {
			log.Println(err)
			continue
		}
		switch p.Type {
		case packet.ConnectionInit:
			keeperID, err = p.Payload.StringValue("keeper_id")
			if err != nil {
				log.Println(err)
			}
			keepers[keeperID] = c
		case packet.Request:
			atomic.AddUint64(&requestN, 1)
			requests[requestN] = c
			if err := handleRequest(&p); err != nil {
				log.Println(err)
			}
		case packet.Next:
			c, ok := requests[p.ID]
			if !ok {
				log.Println(err)
				continue
			}
			delete(requests, p.ID)
			if err := c.WriteJSON(&p); err != nil {
				log.Println(err)
			}
		}
	}
}

func handleRequest(p *packet.Packet) error {
	keeperID, err := p.Payload.StringValue("keeper_id")
	if err != nil {
		return err
	}
	handler, err := p.Payload.StringValue("handler")
	if err != nil {
		return err
	}
	keeper, ok := keepers[keeperID]
	if !ok {
		return errors.New("keeper is not found")
	}
	switch handler {
	case "new_address":
		kp := packet.Packet{
			ID:      requestN,
			Type:    packet.Request,
			Payload: map[string]interface{}{"handler": "new_address"},
		}
		if err := keeper.WriteJSON(&kp); err != nil {
			return err
		}
		return nil
	default:
		return errors.New("handler is not found")
	}
}
