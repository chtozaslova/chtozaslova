package chtozaslova

import (
	"math"
)

func W(Y int) int {
	floor := int(math.Floor(1546 * math.Sin(deg2rad((float64(Y)+0.5)/24))))
	if floor > 1 {
		return floor
	} else {
		return 1
	}
}

func deg2rad(deg float64) float64 {
	return deg * math.Pi / 180
}

func cuberoot(n float64) float64 {
	return math.Pow(n, 1.0/3.0)
}

func frac(n float64) float64 {
	return n - math.Floor(n)
}
