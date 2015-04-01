package gtfs

type Stop struct {
	Id    string
	Code string
	Name  string
	Desc string
	Lat float64
	Lon float64
	FareZone FareZone
	Url string
	LocationType int
	// TODO: ParentStation?
	Timezone string
	WheelchairBoarding int
	Coord Coord
}

type StopTime struct {
	Stop *Stop
	Trip *Trip
	Time int
	Seq  int
}

type StopTimeBySeq []StopTime

func (a StopTimeBySeq) Len() int           { return len(a) }
func (a StopTimeBySeq) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a StopTimeBySeq) Less(i, j int) bool { return a[i].Seq < a[j].Seq }