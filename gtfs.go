package gtfs

import (
	"encoding/csv"
	"io"
	logger "log"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
)

var (
	log = logger.New(os.Stderr, "", logger.LstdFlags)
)


type Headsign struct {
	Direction string
	Text      string
}

type Shape struct {
	Id     string
	Coords []Coord
}

type CalendarEntry struct {
	ServiceId string
	Days      []string
}

type Coord struct {
	Lat float64
	Lon float64
	Seq int
}

type CoordBySeq []Coord

func (a CoordBySeq) Len() int           { return len(a) }
func (a CoordBySeq) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a CoordBySeq) Less(i, j int) bool { return a[i].Seq < a[j].Seq }

// main utility function for reading GTFS files
func (feed *FeedInfo) readCsv(filename string, f func([]string)) {
	file, err := os.Open(path.Join(feed.Dir, filename))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	reader.TrailingComma = true
	firstLineSeen := false
	for {
		record, err := reader.Read()
		if !firstLineSeen {
			firstLineSeen = true
			continue
		}
		if err == io.EOF {
			break
		} else if err != nil {
			log.Panicln(err)
		} else {
			f(record)
		}
	}
}

func Load(feed_path string, loadStopTimes bool) FeedInfo {
	f := FeedInfo{Dir: feed_path}
	f.Routes = make(map[string]*Route)
	f.Shapes = make(map[string]*Shape)
	f.Stops = make(map[string]*Stop)
	f.Trips = make(map[string]*Trip)
	f.CalendarEntries = make(map[string]CalendarEntry)

	f.readCsv("calendar.txt", func(s []string) {
		c := CalendarEntry{ServiceId: s[0], Days: s[1:8]}
		f.CalendarEntries[s[0]] = c
	})

	// we assume that this CSV is grouped by shape_id
	// but this is not guaranteed in spec?
	var curShape *Shape
	var found = false
	f.readCsv("shapes.txt", func(s []string) {
		shape_id := s[0]
		if !found || shape_id != curShape.Id {
			if found {
				f.Shapes[curShape.Id] = curShape
			}
			found = true
			curShape = &Shape{Id: shape_id}
		}
		lon, _ := strconv.ParseFloat(s[1], 64)
		lat, _ := strconv.ParseFloat(s[2], 64)
		seq, _ := strconv.Atoi(s[3])
		curShape.Coords = append(curShape.Coords, Coord{Lat: lat, Lon: lon, Seq: seq})
	})
	if found {
		f.Shapes[curShape.Id] = curShape
	}

	// sort coords by their sequence
	for _, v := range f.Shapes {
		sort.Sort(CoordBySeq(v.Coords))
	}

	f.readCsv("routes.txt", f.loadRoute)

	f.readCsv("trips.txt", func(s []string) {
		route_id := s[0]
		service := s[1]
		trip_id := s[2]
		direction := s[4]
		shape_id := s[6]
		headsign := s[3]

		var shape *Shape
		shape = f.Shapes[shape_id]
		var trip Trip
		trip.StopTimes = []StopTime{}
		f.Trips[trip_id] = &trip

		route := f.Routes[route_id]
		trip = Trip{Shape: shape, Route: route, Id: trip_id, Direction: direction, Service: service, Headsign: headsign}
		route.Trips = append(route.Trips, &trip)
		f.Routes[route_id] = route
	})

	f.readCsv("stops.txt", func(s []string) {
		stop_id := s[0]
		stop_name := s[1]
		stop_lat, _ := strconv.ParseFloat(s[3], 64)
		stop_lon, _ := strconv.ParseFloat(s[4], 64)
		coord := Coord{Lat: stop_lat, Lon: stop_lon}
		f.Stops[stop_id] = &Stop{Coord: coord, Name: stop_name, Id: stop_id}
	})

	if !loadStopTimes {
		return f
	}
	f.readCsv("stop_times.txt", func(s []string) {
		trip_id := s[0]
		stop_id := s[3]
		seq, _ := strconv.Atoi(s[4])
		time := Hmstoi(s[1])
		stop := f.Stops[stop_id]
		trip := f.Trips[trip_id]
		newStopTime := StopTime{Trip: trip, Stop: stop, Seq: seq, Time: time}
		trip.StopTimes = append(trip.StopTimes, newStopTime)
	})

	// sort stops by seq

	for _, v := range f.Trips {
		sort.Sort(StopTimeBySeq(v.StopTimes))
	}

	return f
}

func Hmstoi(str string) int {
	components := strings.Split(str, ":")
	hour, _ := strconv.Atoi(components[0])
	min, _ := strconv.Atoi(components[1])
	sec, _ := strconv.Atoi(components[2])
	retval := hour*60*60 + min*60 + sec
	return retval
}
