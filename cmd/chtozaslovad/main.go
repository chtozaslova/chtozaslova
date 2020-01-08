package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"chtozaslova/chtozaslova"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Chtozaslova API server. Try /api/convert-to-3wa and /api/convert-to-coordinates.\n")
	})

	http.HandleFunc("/api/convert-to-3wa", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		params := r.URL.Query()

		coordinates, ok := params["coordinates"]
		if !ok {
			w.WriteHeader(400)
			fmt.Fprintf(w, "{\"error\":\"No coordinates parameter passed to /api/convert-to-3wa.\"}")
			return
		}

		numArr := strings.Split(coordinates[0], ",")
		if len(numArr) != 2 {
			w.WriteHeader(400)
			fmt.Fprintf(w, "{\"error\":\"Expected exactly 2 coordinates separated by comma.\"}")
			return
		}

		lat, err := strconv.ParseFloat(numArr[0], 64)
		if err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "{\"error\":\"Can't parse lat\"}")
			return
		}

		lon, err := strconv.ParseFloat(numArr[1], 64)
		if err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "{\"error\":\"Can't parse lon\"}")
			return
		}

		words, err := chtozaslova.LatLon2Words(lat, lon)
		if err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "{\"error\":\"Error encoding lat,lon\"}")
			return
		}

		fmt.Fprintf(w, "{\"words\":\"%s\",\"language\":\"en\",\"coordinates\":{\"lat\":%.6f,\"lon\":%.6f}}", words, lat, lon)
	})

	http.HandleFunc("/api/convert-to-coordinates", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		params := r.URL.Query()
		words, ok := params["words"]
		if !ok {
			w.WriteHeader(400)
			fmt.Fprintf(w, "{\"error\":\"No words parameter passed to /api/convert-to-coordinates.\"}")
			return
		}

		lat, lon, err := chtozaslova.Words2LatLon(words[0])
		if err != nil {
			w.WriteHeader(400)
			fmt.Fprintf(w, "{\"error\":\"Error decoding words\"}")
			return
		}

		fmt.Fprintf(w, "{\"words\":\"%s\",\"language\":\"en\",\"coordinates\":{\"lat\":%.6f,\"lon\":%.6f}}", words[0], lat, lon)
	})

	log.Fatal(http.ListenAndServe(":8081", nil))
}
