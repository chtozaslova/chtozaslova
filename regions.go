package chtozaslova

type Region struct {
	x1, y  int
	length int
	q      uint64
}

type RegionSlice []Region

func (r RegionSlice) Len() int      { return len(r) }
func (r RegionSlice) Swap(i, j int) { r[i], r[j] = r[j], r[i] }
func (r RegionSlice) Less(i, j int) bool {
	return r[i].x1 < r[j].x1
}


var regions []Region
var regionsY [][]Region