package main

import (
	"bufio"
	"flag"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func main() {
	testFlag := flag.Bool("test", false, "Enable test mode")
	metricsUrl := flag.String("url", "http://127.0.0.1:8080/metrics", "The URL to use to fetch Prometheus metrics")

	flag.Parse()

	mux := http.NewServeMux()

	mux.HandleFunc("/charts", func(w http.ResponseWriter, r *http.Request) {
		var scanner *bufio.Scanner

		if *testFlag == true {
			file, err := os.Open("example_metrics.txt")

			if err != nil {
				log.Println(err)
				http.Error(w, "Internal error", http.StatusInternalServerError)
				return
			}

			defer file.Close()

			scanner = bufio.NewScanner(file)
		} else {
			resp, err := http.Get(*metricsUrl)
			if err != nil {
				log.Println(err)
				http.Error(w, "Internal error", http.StatusInternalServerError)
				return
			}

			defer resp.Body.Close()

			scanner = bufio.NewScanner(resp.Body)
		}

		bm, jm, err := readMetrics(scanner)

		if err != nil {
			log.Println(err)
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}

		chart(bm, jm).Render(r.Context(), w)
	})

	server := http.Server{
		Addr:              ":8081",
		Handler:           mux,
		ReadHeaderTimeout: 3 * time.Second,
	}

	log.Printf("Listening on %v", server.Addr)

	err := server.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}

func readMetrics(scanner *bufio.Scanner) (map[string]int, map[string]int, error) {
	latencyRegex, err := regexp.Compile("\\d+|\\+Inf")

	if err != nil {
		return nil, nil, err
	}

	bm := map[string]int{}
	jm := map[string]int{}

	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "binary_message_latency_bucket") {
			before, after, _ := strings.Cut(line, " ")

			time := latencyRegex.FindString(before)
			numRequests, err := strconv.ParseFloat(after, 32)

			if err != nil {
				return nil, nil, err
			}

			bm[time] = int(numRequests)
		}

		if strings.Contains(line, "json_message_latency_bucket") {
			before, after, _ := strings.Cut(line, " ")

			time := latencyRegex.FindString(before)
			numRequests, err := strconv.ParseFloat(after, 32)

			if err != nil {
				return nil, nil, err
			}

			jm[time] = int(numRequests)
		}
	}

	if err = scanner.Err(); err != nil {
		return nil, nil, err
	}

	return bm, jm, nil
}
