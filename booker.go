package main

import (
	"demo/packet"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

var upgrader = websocket.Upgrader{}

var keepers map[string]*websocket.Conn

func main() {
	keepers = make(map[string]*websocket.Conn)
	handleConnectKeeper()
	handleWebsocket()
	log.Fatalln(http.ListenAndServe(":8000", nil))
}

func handleWebsocket() {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		core, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			println(err.Error())
		}
		defer core.Close()
		for {
			_, r, err := core.NextReader()
			if err != nil {
				log.Println(err)
				break
			}
			var p packet.Transport
			if err := json.NewDecoder(r).Decode(&p); err != nil {
				log.Println(err)
				continue
			}

			if p.Handler == "new_address" {
				//--------------->
				keeper, ok := keepers[p.KeeperID]
				if !ok {
					log.Println(err)
					continue
				}
				err = keeper.WriteMessage(websocket.TextMessage, []byte(p.Handler))
				if err != nil {
					log.Println(err)
					break
				}
				mt, m, err := keeper.ReadMessage()
				if err != nil {
					log.Println(err)
					break
				}
				address := string(m)
				//<---------------

				err = core.WriteMessage(mt, []byte(address))
				if err != nil {
					log.Println(err)
					break
				}
			}
		}
	})
}

func handleConnectKeeper() {
	http.HandleFunc("/connect", func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			println(err.Error())
		}
		_, reader, err := c.NextReader()
		if err != nil {
			log.Println(err)
			return
		}
		var p packet.Connection
		if err := json.NewDecoder(reader).Decode(&p); err != nil {
			log.Println(err)
			return
		}
		keepers[p.KeeperID] = c
		log.Println("connected_at " + time.Now().Format("15:04:05"))
	})
}
