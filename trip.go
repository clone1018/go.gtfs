package gtfs


type Trip struct {
	Id        string
	Shape     *Shape
	Route     *Route
	Service   string
	Direction string
	Headsign  string

	// may not be loaded
	StopTimes []StopTime
}