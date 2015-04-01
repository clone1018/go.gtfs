package gtfs

type FeedInfo struct {
    Dir             string
    Routes          map[string]*Route
    Shapes          map[string]*Shape
    Stops           map[string]*Stop
    Trips           map[string]*Trip
    CalendarEntries map[string]CalendarEntry
}


func (feed *FeedInfo) RouteByShortName(shortName string) *Route {
    for _, v := range feed.Routes {
        if v.ShortName == shortName {
            return v
        }
    }
    //TODO error here
    return &Route{}
}

func (feed FeedInfo) Calendar() []string {
    retval := []string{}
    for i := 0; i <= 6; i++ {
        for k, v := range feed.CalendarEntries {
            if v.Days[i] == "1" {
                retval = append(retval, k)
            }
        }
    }
    return retval
}
