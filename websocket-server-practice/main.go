package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type connPool struct {
	connections []*websocket.Conn
}

func main() {
	var connPool connPool
	upgrader := &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		connPool.connections = append(connPool.connections, c)
		if err != nil {
			log.Println("upgrade:", err)
			return
		}
		fmt.Println("New connection", c.RemoteAddr())
		defer func() {
			log.Println("disconnect !!")
			c.Close()
		}()
		for {
			mtype, msg, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				break
			}
			if mtype == websocket.CloseNormalClosure {
				log.Println("closeNormal")
				break
			}
			log.Printf("receive: %s\n", msg)
			rev := string(msg)
			send := "broadcast:" + rev
			for _, s_c := range connPool.connections {
				err = s_c.WriteMessage(mtype, []byte(send))
				if err != nil {
					log.Println("write:", err)
					break
				}
			}

		}
	})
	log.Println("server start at :8899")
	log.Fatal(http.ListenAndServe(":8899", nil))
}
