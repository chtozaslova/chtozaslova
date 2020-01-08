package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"chtozaslova/chtozaslova"
)

func latlon2words(input string) {
	numArr := strings.Split(input, ",")
	if len(numArr) != 2 {
		fmt.Printf("%s\texpected lat,lon\n", input)
		return
	}

	lat, err := strconv.ParseFloat(numArr[0], 64)
	if err != nil {
		fmt.Printf("%s\tcouldn't understand lat: %s\n", input, numArr[0])
		return
	}

	lon, err := strconv.ParseFloat(numArr[1], 64)
	if err != nil {
		fmt.Printf("%s\tcouldn't understand lon: %s\n", input, numArr[1])
		return
	}

	words, err := chtozaslova.LatLon2Words(lat, lon)
	if err != nil {
		fmt.Printf("%s\t%v\n", input, err)
		return
	}

	fmt.Printf("%s\t%s\n", input, words)
}

func words2latlon(input string) {
	lat, lon, err := chtozaslova.Words2LatLon(input)
	if err != nil {
		fmt.Printf("%s\t%v\n", input, err)
		return
	}

	fmt.Printf("%s\t%.6f,%.6f\n", input, lat, lon)
}

func isNumeric(c byte) bool {
	return c >= '0' && c <= '9' || c == '-'
}

func handleInput(input string) {
	if isNumeric(input[0]) {
		latlon2words(input)
	} else {
		words2latlon(input)
	}
}

func main() {
	readFromStdinPtr := flag.Bool("i", false, "read from stdin")
	flag.Parse()

	if *readFromStdinPtr {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			handleInput(scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "chtozaslova: read stdin: %v", err)
			os.Exit(1)
		}
	} else {
		if len(os.Args) < 2 {
			fmt.Fprintf(os.Stderr, "usage: chtozaslova [-i] [INPUT...]\n\nOptions:\n  -i    read from stdin\n\nIf not -i, you should specify at least one INPUT in the form lat,lon or words. It is perfectly acceptable to mix lat,lon and words inputs in a single invocation.\n\nExamples:\n    chtozaslova -i < words.txt\n    chtozaslova joyful.nail.harmonica\n    chtozaslova 37.234332,-115.806657\n\nOutput is of the form \"INPUT[tab]OUTPUT\" with one line per INPUT, where OUTPUT will be either the converted INPUT, or an error message.\n\n")
			os.Exit(1)
		}

		for _, input := range os.Args[1:] {
			handleInput(input)
		}
	}
}
