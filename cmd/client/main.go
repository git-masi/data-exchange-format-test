package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"log"
	"math/rand"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

type JsonData struct {
	Time      int64   `json:"time"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

func main() {
	messageType := flag.String("type", "bin", "Which type of message to use")

	flag.Parse()

	switch *messageType {
	case "bin":
		handleBin()
	case "json":
		handleJson()
	default:
		log.Fatalln("No valid message type supplied")
	}
}

func handleBin() {
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

	log.Println("Test started")

	for {
		select {
		case <-timer.C:
			err = conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(1000, ""), time.Now().Add(time.Millisecond))

			if err != nil {
				log.Println(err)
			}

			log.Println("Test ended")
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

func handleJson() {
	serverURL := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/json"}

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

	log.Println("Test started")

	for {
		select {
		case <-timer.C:
			err = conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(1000, ""), time.Now().Add(time.Millisecond))

			if err != nil {
				log.Println(err)
			}

			log.Println("Test ended")
			return
		default:
			{
				message := JsonData{
					Time:      time.Now().UnixMilli(),
					Latitude:  rand.Float32(),
					Longitude: rand.Float32(),
				}

				err = conn.WriteJSON(message)

				if err != nil {
					log.Println(err)
					return
				}
			}
		}
	}
}
