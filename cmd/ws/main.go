package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			log.Println(err)
			return
		}

		defer conn.Close()

		for {
			_, message, err := conn.ReadMessage()

			if err != nil {
				log.Println("Error reading message from WebSocket", err)
				break
			}

			buf := bytes.NewReader(message)

			var timestamp int64

			err = binary.Read(buf, binary.BigEndian, &timestamp)

			if err != nil {
				log.Println("Error reading timestamp", err)
				break
			}

			var latitude float32

			err = binary.Read(buf, binary.BigEndian, &latitude)

			if err != nil {
				log.Println("Error reading latitude", err)
				break
			}

			var longitude float32

			err = binary.Read(buf, binary.BigEndian, &longitude)

			if err != nil {
				log.Println("Error reading longitude", err)
				break
			}

			t := time.UnixMilli(timestamp)

			log.Printf("time %v, latitude %v, longitude %v", t, latitude, longitude)
		}
	})

	server := http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 3 * time.Second,
	}

	log.Printf("Listening on %v", server.Addr)

	err := server.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}
