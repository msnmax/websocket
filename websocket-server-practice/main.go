package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var RecordLogger *log.Logger

func manageUserName(pool map[string]*websocket.Conn, conn *websocket.Conn) (string, error) {
	var username string
	for {
		mtype, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("name input err:", err)
			return "", err

		}
		fmt.Println("name:", string(msg))
		if mtype == websocket.CloseNormalClosure {
			log.Println("closeNormal")
			return "", err
		}
		_, ok := pool[string(msg)]
		if ok {
			continue
		}
		pool[string(msg)] = conn
		username = string(msg)
		break
	}

	return username, nil
}

func initRecordLog() (*os.File, error) {
	file, err := os.OpenFile("record.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	//defer file.Close()
	RecordLogger = log.New(file, "", log.Ldate|log.Ltime)
	return file, err
}
func record(msg string) {
	RecordLogger.Println(msg)
}
func broadcast(pool map[string]*websocket.Conn, user string, msg string) {

	for name, c := range pool {
		if name == user {
			continue
		}

		err := c.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func main() {
	connPool := make(map[string]*websocket.Conn)
	file, err := initRecordLog()
	upgrader := &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	if err != nil {
		log.Fatalf("Failed to initialize log file: %v", err)
	}
	defer file.Close()

	http.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		username, erruser := manageUserName(connPool, c)

		if erruser != nil {
			return
		}
		fmt.Println(connPool)
		if err != nil {
			log.Println("upgrade:", err)
			return
		}
		fmt.Println("New connection", c.RemoteAddr())
		defer func() {
			log.Println("disconnect !!")
			delete(connPool, username)
			c.Close()
		}()

		broadcast(connPool, "", fmt.Sprintf("[%s]-into chat", username))
		record(username + " into chat")
		if err != nil {
			log.Println("write:", err)
		}
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
			send := username + ":" + rev
			recordmsg := fmt.Sprintf("[%s]:%s", username, rev)
			record(recordmsg)
			broadcast(connPool, username, send)

		}
	})

	http.HandleFunc("/dump", func(w http.ResponseWriter, r *http.Request) {

		if err != nil {
			http.Error(w, "Failed to encode data to JSON", http.StatusInternalServerError)
			return
		}

		// 设置 HTTP 响应头部
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		data, _ := os.ReadFile("record.log")
		fmt.Fprintf(w, "%s\n", string(data))
	})

	log.Println("server starting")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
