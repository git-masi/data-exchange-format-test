package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"math/rand"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	serverURL := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/bin"}

	conn, _, err := websocket.DefaultDialer.Dial(serverURL.String(), nil)

	if err != nil {
		log.Fatal("WebSocket connection error:", err)
	}

	defer conn.Close()

	go func() {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				log.Println("WebSocket read error:", err)
				return
			}
		}
	}()

	timer := time.NewTimer(1 * time.Minute)

	// Create a buffer big enough to store a timestamp and 2x 32 bit floats
	buf := bytes.NewBuffer(make([]byte, binary.Size(int64(0))+binary.Size(float32(0))*2))

	log.Println("Starting test")

	for {
		select {
		case <-timer.C:
			log.Println("Finished test")
			return
		default:
			{
				if err := binary.Write(buf, binary.BigEndian, time.Now().UnixMilli()); err != nil {
					log.Println(err)
					return
				}

				for i := 0; i < 2; i++ {
					if err := binary.Write(buf, binary.BigEndian, rand.Float32()); err != nil {
						log.Println(err)
						return
					}
				}

				err := conn.WriteMessage(websocket.BinaryMessage, buf.Bytes())

				if err != nil {
					log.Println("WebSocket write error:", err)
					return
				}

				buf.Reset()
			}
		}
	}
}
