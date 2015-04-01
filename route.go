package gtfs

import (
	"strings"
)

//Transit routes. A route is a group of trips that are displayed to riders as a single service.
type Route struct {
	Id        string
	Agency Agency
	ShortName string
	LongName  string
	Desc string

	// The route_type field describes the type of transportation used on a route. Valid values for this field are:
	// 0 - Tram, Streetcar, Light rail. Any light rail or street level system within a metropolitan area.
	// 1 - Subway, Metro. Any underground rail system within a metropolitan area.
	// 2 - Rail. Used for intercity or long-distance travel.
	// 3 - Bus. Used for short- and long-distance bus routes.
	// 4 - Ferry. Used for short- and long-distance boat service.
	// 5 - Cable car. Used for street-level cable cars where the cable runs beneath the car.
	// 6 - Gondola, Suspended cable car. Typically used for aerial cable cars where the car is suspended from the cable.
	// 7 - Funicular. Any rail system designed for steep inclines.
	Type      int
	Url       string
	Color     string
	TextColor string
	Trips     []*Trip
}

// get All shapes for a route
func (route Route) Shapes() []*Shape {
	// collect the unique list of shape pointers
	hsh := make(map[*Shape]bool)

	for _, v := range route.Trips {
		hsh[v.Shape] = true
	}

	retval := []*Shape{}
	for k, _ := range hsh {
		retval = append(retval, k)
	}
	return retval
}

func (route Route) LongestShape() *Shape {
	max := 0
	var shape *Shape
	for _, s := range route.Shapes() {
		if len(s.Coords) > max {
			shape = s
			max = len(s.Coords)
		}
	}
	return shape
}

func (route Route) Stops() []*Stop {
	stops := make(map[*Stop]bool)
	// can't assume the longest shape includes all stops

	for _, t := range route.Trips {
		for _, st := range t.StopTimes {
			stops[st.Stop] = true
		}
	}

	retval := []*Stop{}
	for k, _ := range stops {
		retval = append(retval, k)
	}
	return retval
}

func (route Route) Headsigns() []string {
	max0 := 0
	maxHeadsign0 := ""
	max1 := 1
	maxHeadsign1 := ""

	for _, t := range route.Trips {
		if t.Direction == "0" {
			if len(t.Shape.Coords) > max0 {
				max0 = len(t.Shape.Coords)
				maxHeadsign0 = strings.TrimSpace(t.Headsign)
			}
		} else { // direction == 1. only bidirectional
			if len(t.Shape.Coords) > max1 {
				max1 = len(t.Shape.Coords)
				maxHeadsign1 = strings.TrimSpace(t.Headsign)
			}
		}
	}

	return []string{maxHeadsign0, maxHeadsign1}
}

func (f FeedInfo) loadRoute(s []string) {
	rsn := strings.TrimSpace(s[2])
	rln := strings.TrimSpace(s[3])
	id := strings.TrimSpace(s[0])
	f.Routes[id] = &Route{Id: id, ShortName: rsn, LongName: rln}
}
