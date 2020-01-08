package main

import (
	"fmt"

	"chtozaslova/chtozaslova"
)

func main() {
	words, _ := chtozaslova.LatLon2Words(37.234332, -115.806663)
	fmt.Println(words)
	lat, lon, _ := chtozaslova.Words2LatLon("joyful.nail.harmonica")
	fmt.Printf("%.6f %.6f\n", lat, lon)
}
