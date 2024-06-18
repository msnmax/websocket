package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	var send string
	c, _, err := websocket.DefaultDialer.Dial("ws://127.0.0.1:8899/echo", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()
	fmt.Println("Connected")
	go func() {
		for {

			_, msg, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			fmt.Printf("%s\n", msg)
		}
	}()
	for {

		fmt.Scan(&send)
		if send == "exit" {
			_ = c.WriteControl(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now())
			break
		}
		fmt.Println("<You> ", send)
		err = c.WriteMessage(websocket.TextMessage, []byte(send))
		if err != nil {
			log.Println(err)
			return
		}

	}

	/*
		c, _, err := websocket.DefaultDialer.Dial("ws://127.0.0.1:8899/echo", nil)
		if err != nil {
			log.Fatal("dial:", err)
		}
		defer c.Close()

		err = c.WriteMessage(websocket.TextMessage, []byte("hello "))
		if err != nil {
			log.Println(err)
			return
		}
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		log.Printf("receive: %s\n", msg)
	*/

}
