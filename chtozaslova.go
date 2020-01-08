package chtozaslova

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io"
	"math"
	"math/big"
	"sort"
	"strings"
)

func LatLon2Words(lat, lon float64) (string, error) {
	X, Y, x, y := LatLon2XYxy(lat, lon)

	n, err := XYxy2n(X, Y, x, y)
	if err != nil {
		return "", err
	}

	m, err := N2m(n)
	if err != nil {
		return "", err
	}

	i, j, k, err := M2ijk(m)
	if err != nil {
		return "", err
	}

	words, err := Ijk2words(i, j, k)
	if err != nil {
		return "", err
	}

	return words, nil
}

func Words2LatLon(words string) (lat, lon float64, err error) {
	i, j, k, err := Words2ijk(words)
	if err != nil {
		return 0, 0, err
	}

	m := Ijk2m(i, j, k)

	n, err := M2n(m)
	if err != nil {
		return 0, 0, err
	}

	X, Y, x, y, err := N2XYxy(n)
	if err != nil {
		return 0, 0, err
	}

	lat, lon = XYxy2LatLon(X, Y, x, y)

	return lat, lon, nil
}

func LatLon2XYxy(lat, lon float64) (X, Y, x, y int) {
	for lon < -180 {
		lon += 360
	}
	for lon >= 180 {
		lon -= 360
	}
	for lat < -90 {
		lat += 180
	}
	for lat >= 90 {
		lat -= 180
	}

	X = int(math.Floor((lon + 180) * 24))
	Y = int(math.Floor((lat + 90) * 24))
	x = int(math.Floor(float64(W(Y)) * frac((lon+180)*24)))
	y = int(math.Floor(1546 * frac((lat+90)*24)))
	return X, Y, x, y
}

func XYxy2LatLon(X, Y, x, y int) (lat, lon float64) {
	lat = (float64(Y)+(float64(y)+0.5)/1546)/24 - 90
	lon = (float64(X)+(float64(x)+0.5)/float64(W(Y)))/24 - 180

	return lat, lon
}

func XYxy2n(X, Y, x, y int) (uint64, error) {
	q, err := XY2q(X, Y)
	if err != nil {
		return 0, err
	}
	return q + uint64(x)*1546 + uint64(y), nil
}

func N2XYxy(n uint64) (X, Y, x, y int, err error) {
	X, Y, q, err := N2XYq(n)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	n -= q
	x = int(math.Floor(float64(n) / 1546))
	y = int(n - uint64(x)*1546)

	return X, Y, x, y, nil
}

func M2ijk(m uint64) (i, j, k int, err error) {
	l := uint64(math.Floor(cuberoot(float64(m))))
	l2 := l * l
	l3 := l * l * l
	if l3 <= m && m < l3+l2+2*l+1 {
		r := m - l3
		i = int(l)
		j = int(r / (l + 1))
		k = int(r % (l + 1))
	} else if l3+l2+2*l+1 <= m && m < l3+2*l2+3*l+1 {
		r := m - (l3 + l2 + 2*l + 1)
		i = int(r / (l + 1))
		j = int(l)
		k = int(r % (l + 1))
	} else {
		r := m - (l3 + 2*l2 + 3*l + 1)
		i = int(r / l)
		j = int(r % l)
		k = int(l)
	}

	if i >= len(num2word) || j >= len(num2word) || k >= len(num2word) {
		return 0, 0, 0, fmt.Errorf("chtozaslova: m out of range: %d", m)
	}

	return i, j, k, nil
}

func Ijk2m(i, j, k int) uint64 {
	l := uint64(i)
	if j > i && j > k {
		l = uint64(j)
	} else if k > i && k > j {
		l = uint64(k)
	}

	l2 := l * l
	l3 := l * l * l

	if uint64(i) == l {
		return l3 + (l+1)*uint64(j) + uint64(k)
	} else if uint64(j) == l {
		return l3 + l2 + 2*l + 1 + (l+1)*uint64(i) + uint64(k)
	} else {
		return l3 + 2*l2 + 3*l + 1 + l*uint64(i) + uint64(j)
	}
}

func Ijk2words(i, j, k int) (string, error) {
	if i >= len(num2word) || j >= len(num2word) || k >= len(num2word) || i < 0 || j < 0 || k < 0 {
		return "", fmt.Errorf("chtozaslova: i,j,k value out of range: %d,%d,%d", i, j, k)
	}
	s := []string{num2word[i], num2word[j], num2word[k]}
	return strings.Join(s, "."), nil
}

func Words2ijk(words string) (i, j, k int, err error) {
	wordsArr := strings.Split(words, ".")
	if len(wordsArr) != 3 {
		return 0, 0, 0, fmt.Errorf("chtozaslova: expected 3 words but got %d: %s", len(wordsArr), words)
	}
	i, ok := word2num[wordsArr[0]]
	if !ok {
		return 0, 0, 0, fmt.Errorf("chtozaslova: don't recognise word: %s", wordsArr[0])
	}
	j, ok = word2num[wordsArr[1]]
	if !ok {
		return 0, 0, 0, fmt.Errorf("chtozaslova: don't recognise word: %s", wordsArr[1])
	}
	k, ok = word2num[wordsArr[2]]
	if !ok {
		return 0, 0, 0, fmt.Errorf("chtozaslova: don't recognise word: %s", wordsArr[2])
	}
	return i, j, k, nil
}

func N2m(n uint64) (uint64, error) {
	for _, blk := range shuffleBlocks {
		if n >= blk.start && n < blk.start+blk.size {
			blockstart := big.NewInt(int64(blk.start))
			blocksize := big.NewInt(int64(blk.size))
			F_i := big.NewInt(int64(blk.f_i))
			n := big.NewInt(int64(n))

			n.Sub(n, blockstart).Mul(n, F_i).Mod(n, blocksize).Add(n, blockstart)
			return n.Uint64(), nil
		}
	}

	return 0, fmt.Errorf("chtozaslova: n out of range: %d", n)
}

func M2n(m uint64) (uint64, error) {
	for _, blk := range shuffleBlocks {
		if m >= blk.start && m < blk.start+blk.size {
			blockstart := big.NewInt(int64(blk.start))
			blocksize := big.NewInt(int64(blk.size))
			R_i := big.NewInt(int64(blk.r_i))
			m := big.NewInt(int64(m))

			m.Sub(m, blockstart).Mul(m, R_i).Mod(m, blocksize).Add(m, blockstart)
			return m.Uint64(), nil
		}
	}

	return 0, fmt.Errorf("chtozaslova: m out of range: %d", m)
}

func XY2q(X, Y int) (uint64, error) {
	buildRegions()

	min := 0
	max := len(regionsY[Y])

	for min < max {
		i := (min + max) / 2
		b := regionsY[Y][i]

		if b.x1 <= X && b.x1+b.length > X && b.y == Y {
			return b.q + uint64(X-b.x1)*uint64(W(Y)*1546), nil
		} else if X < b.x1 {
			max = i
		} else {
			min = i + 1
		}
	}

	return 0, fmt.Errorf("chtozaslova: no such cell (%d,%d)", X, Y)
}

func N2XYq(n uint64) (X, Y int, q uint64, err error) {
	buildRegions()

	min := 0
	max := len(regions)

	for min < max {
		i := (min + max) / 2
		b := regions[i]

		onecellSize := W(b.y) * 1546
		endq := b.q + uint64(b.length*onecellSize)

		if n >= b.q && n < endq {
			i := int((n - b.q) / uint64(onecellSize))
			return b.x1 + i, b.y, b.q + uint64(i*onecellSize), nil
		} else if n < b.q {
			max = i
		} else {
			min = i + 1
		}
	}

	return 0, 0, 0, fmt.Errorf("chtozaslova: n out of range: %d", n)
}

func buildRegions() {
	if len(regions) > 0 {
		return
	}

	compressed, err := base64.StdEncoding.DecodeString(regionData)
	if err != nil {
		panic("chtozaslova: can't un-base64 compressed region data")
	}

	regionReader, err := gzip.NewReader(strings.NewReader(string(compressed)))
	if err != nil {
		panic("chtozaslova: can't read compressed region data out of memory")
	}

	var regionBytes bytes.Buffer
	if _, err := io.Copy(&regionBytes, regionReader); err != nil {
		panic("chtozaslova: can't decompress region data from memory")
	}

	r := regionBytes.String()

	regions = make([]Region, len(r)/6)
	regionsY = make([][]Region, 4320)

	for i := 0; i < 4320; i++ {
		regionsY[i] = make([]Region, 0)
	}

	q := uint64(0)

	for i := 0; i < len(r)/6; i++ {
		x := int(r[i*6]) + 256*int(r[i*6+1])
		y := int(r[i*6+2]) + 256*int(r[i*6+3])
		length := int(r[i*6+4]) + 256*int(r[i*6+5])
		region := Region{x, y, length, q}
		regions[i] = region
		regionsY[y] = append(regionsY[y], region)
		q += uint64(length) * uint64(W(y)) * 1546
	}

	for i := 0; i < 4320; i++ {
		sort.Sort(RegionSlice(regionsY[i]))
	}
}
