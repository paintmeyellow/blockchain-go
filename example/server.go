package main

import (
	"demo/wsconn"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wc, err := wsconn.Upgrade(w, r, nil)
		if err != nil {
			log.Println("SERVER", err)
			return
		}
		wc.Subscribe("dc", func(m *wsconn.Msg) {
			fmt.Println("SERVER:recv", "id:", m.ID, "reply:", m.Reply, string(m.Data))
			if err = wc.RespondOnMsg(m, []byte("Hello from server, DC!")); err != nil {
				panic(err)
			}
		})
	})
	log.Fatalln(http.ListenAndServe(":8000", nil))
}
