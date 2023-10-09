package main

import (
	"bufio"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func fetchChartData() error {
	resp, err := http.Get("http://127.0.0.1:8080/metrics")

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)

	readMetrics(scanner)

	return nil
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
