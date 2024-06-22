package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	//var send string
	c, _, err := websocket.DefaultDialer.Dial("ws://127.0.0.1:8080/echo", nil)
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
	fmt.Println("Please, input your user name:")
	///msg_in, _ := reader.ReadString('\n')
	//msg_in = msg_in[:len(msg_in)-1]

	send, _ := reader.ReadString('\n')
	send = strings.TrimSpace(send)
	err = c.WriteMessage(websocket.TextMessage, []byte(send))
	if err != nil {
		log.Println(err)
		return
	}
	for {
		msg_in, _ := reader.ReadString('\n')
		msg_in = strings.TrimSpace(msg_in)

		if msg_in == "exit" {
			_ = c.WriteControl(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""), time.Now())
			break
		}
		fmt.Println("<You> ", msg_in)
		err = c.WriteMessage(websocket.TextMessage, []byte(msg_in))
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
