// Package cta is a good package
package cta

import (
	"errors"
)

type GetTrainArrivalsByStationProps struct {
	Station   Station
	Direction string
}

// GetTrainArrivalsByStation will fetch the train arrivals from the cta api by
// station and optionally direction of travel.
func GetTrainArrivalsByStation(props GetTrainArrivalsByStationProps) ([]TrainArrival, error) {
	return nil, errors.ErrUnsupported
}
