package traintracker

import "time"

const (
	CtaBaseUrl = "https://lapi.transitchicago.com"
)

type Station struct {
	Name string
	Id   uint16
}

type LatLon struct {
	Latitude  float64
	Longitude float64
}

type Train struct {
	RunNumber   uint16
	Line        string
	Destination Station
	Direction   string
	Position    LatLon
}

type TrainArrival struct {
	Station     Station
	Train       Train
	Time        time.Time
	IsDue       bool
	IsScheduled bool
}
