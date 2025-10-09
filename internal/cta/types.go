package cta

import "time"

type Station struct {
	Name string
	Id   string
}

type LatLon struct {
	Latitude  float64
	Longitude float64
}

type Train struct {
	RunNumber   int16
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
