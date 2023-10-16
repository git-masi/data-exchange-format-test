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

	flag.Parse()

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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
			resp, err := http.Get("http://127.0.0.1:8080/metrics")

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

		if strings.Contains(line, "binary_message_latency_microseconds_bucket") {
			before, after, _ := strings.Cut(line, " ")

			time := latencyRegex.FindString(before)
			numRequests, err := strconv.Atoi(after)

			if err != nil {
				return nil, nil, err
			}

			bm[time] = numRequests
		}

		if strings.Contains(line, "json_message_latency_milliseconds_bucket") {
			before, after, _ := strings.Cut(line, " ")

			time := latencyRegex.FindString(before)
			numRequests, err := strconv.Atoi(after)

			if err != nil {
				return nil, nil, err
			}

			jm[time] = numRequests
		}
	}

	if err = scanner.Err(); err != nil {
		return nil, nil, err
	}

	return bm, jm, nil
}
