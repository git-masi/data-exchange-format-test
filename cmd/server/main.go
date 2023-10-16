package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type JsonData struct {
	Time      int64   `json:"time"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var binaryMessageLatency = prometheus.NewHistogram(
	prometheus.HistogramOpts{
		Name:    "binary_message_latency_microseconds",
		Help:    "Latency of binary WebSocket messages in microseconds",
		Buckets: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 15, 20, 25, 30, 40, 50},
	},
)

var jsonMessageLatency = prometheus.NewHistogram(
	prometheus.HistogramOpts{
		Name:    "json_message_latency_microseconds",
		Help:    "Latency of JSON WebSocket messages in microseconds",
		Buckets: []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 15, 20, 25, 30, 40, 50},
	},
)

func init() {
	prometheus.MustRegister(binaryMessageLatency)
	prometheus.MustRegister(jsonMessageLatency)
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	mux.Handle("/metrics", promhttp.Handler())

	mux.HandleFunc("/bin", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			log.Println(err)
			return
		}

		defer conn.Close()

		for {
			_, message, err := conn.ReadMessage()

			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					log.Println("Received close code 1000. Closing connection.")
					return
				}

				log.Println("Error reading message from WebSocket", err)
				return
			}

			startTime := time.Now()

			buf := bytes.NewReader(message)

			var timestamp int64

			err = binary.Read(buf, binary.BigEndian, &timestamp)

			if err != nil {
				log.Println("Error reading timestamp", err)
				return
			}

			var latitude float32

			err = binary.Read(buf, binary.BigEndian, &latitude)

			if err != nil {
				log.Println("Error reading latitude", err)
				return
			}

			var longitude float32

			err = binary.Read(buf, binary.BigEndian, &longitude)

			if err != nil {
				log.Println("Error reading longitude", err)
				return
			}

			// Calculate the size of one record: int64 + 2 * float32
			// Pre-allocate a buffer
			resp := bytes.NewBuffer(make([]byte, 0, binary.Size(int64(0))+2*binary.Size(float32(0))))

			binary.Write(resp, binary.BigEndian, timestamp)

			binary.Write(resp, binary.BigEndian, latitude)

			binary.Write(resp, binary.BigEndian, longitude)

			err = conn.WriteMessage(websocket.BinaryMessage, resp.Bytes())

			if err != nil {
				log.Println("Error sending binary message", err)
				return
			}

			endTime := time.Now()

			latency := endTime.Sub(startTime).Microseconds()

			binaryMessageLatency.Observe(float64(latency))
		}
	})

	mux.HandleFunc("/json", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			log.Println(err)
			return
		}

		defer conn.Close()

		for {
			_, _, err := conn.ReadMessage()

			if err != nil {
				log.Println("Error reading message from WebSocket", err)
				return
			}

			startTime := time.Now()

			var data JsonData

			err = conn.ReadJSON(&data)

			if err != nil {
				log.Println("Error decoding JSON", err)
				return
			}

			err = conn.WriteJSON(data)

			if err != nil {
				log.Println("Error encoding JSON", err)
				return
			}

			endTime := time.Now()

			latency := endTime.Sub(startTime).Microseconds()

			jsonMessageLatency.Observe(float64(latency))
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
