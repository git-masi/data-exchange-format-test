package main

import (
	"bufio"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadMetrics(t *testing.T) {
	file, err := os.Open("example.txt")

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	// Create a new scanner
	scanner := bufio.NewScanner(file)

	bm, jm, err := readMetrics(scanner)

	assert.Nil(t, err, "there should be no error when reading metrics")

	want := map[string]int{
		"5":    0,
		"10":   1,
		"15":   8,
		"20":   35,
		"25":   99,
		"30":   290,
		"35":   536,
		"40":   698,
		"45":   812,
		"50":   900,
		"55":   971,
		"60":   1015,
		"65":   1030,
		"70":   1038,
		"75":   1042,
		"80":   1047,
		"+Inf": 1059,
	}

	assert.Equal(t, want, bm, "the returned binary metrics should match the expected value")

	want = map[string]int{
		"5":    2,
		"6":    2,
		"7":    8,
		"8":    93,
		"9":    151,
		"10":   807,
		"11":   1001,
		"12":   1028,
		"13":   1029,
		"14":   1030,
		"15":   1030,
		"16":   1032,
		"17":   1032,
		"18":   1032,
		"19":   1032,
		"20":   1032,
		"+Inf": 1032,
	}

	assert.Equal(t, want, jm, "the returned json metrics should match the expected value")
}
