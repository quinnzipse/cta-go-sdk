package main

import (
	"errors"
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"git.quinnzipse.dev/cta-go-sdk/internal/traintracker"
)

var getStopName = traintracker.GetStopName

func main() {
	fmt.Println("Initializing Train Tracker")

	key := os.Getenv("CTA_API_KEY")
	if strings.TrimSpace(key) == "" {
		panic(errors.New("CTA_API_KEY is required"))
	}

	tracker := traintracker.NewTrainTracker(traintracker.TrainTrackerProps{
		Key: key,
	})

	displayBlueLineLocations(tracker)

	fmt.Println("\n" + strings.Repeat("=", 60) + "\n")

	displayBlueLineArrivals(tracker)
}

func displayBlueLineLocations(tracker traintracker.TrainTracker) {
	fmt.Println("=== Blue Line Train Locations ===")

	locations, err := tracker.Locations(traintracker.LocationsProps{
		Rt: "Blue",
	})

	if err != nil {
		fmt.Printf("Error fetching locations: %v\n", err)
		return
	}

	if locations.Ctatt.ErrCd != "0" {
		fmt.Printf("API Error (%s): %s\n", locations.Ctatt.ErrCd, locations.Ctatt.ErrNm)
		return
	}

	fmt.Printf("Response: ErrCd=%s, Route count: %d\n", locations.Ctatt.ErrCd, len(locations.Ctatt.Route))

	if len(locations.Ctatt.Route) > 0 {
		for _, r := range locations.Ctatt.Route {
			fmt.Printf("DEBUG: Route.Rt='%s', Position len=%d\n", r.Rt, len(r.Position))
			if len(r.Position) == 0 {
				fmt.Println("No trains in this route")
				continue
			}

			for _, train := range r.Position {
				flags := ""
				if train.IsDly {
					flags += " DELAYED"
				}
				if train.IsSch {
					flags += " SCHEDULED"
				}

				fmt.Printf("  Run %s | (% .5f, % .5f) | hdg %d | next: %s (%s)%s\n",
					train.RunNumber,
					train.Lat, train.Lon,
					train.Heading,
					train.NextStpNm, train.NextStpId,
					flags)
			}
		}
	} else {
		fmt.Println("No train data returned")
	}

	if len(locations.Ctatt.Route) == 0 {
		fmt.Println("No train data returned - trying without route filter...")

		allLocs, err := tracker.Locations(traintracker.LocationsProps{})
		if err != nil {
			fmt.Printf("Error fetching all locations: %v\n", err)
			return
		}
		fmt.Printf("All routes count: %d\n", len(allLocs.Ctatt.Route))
		for _, r := range allLocs.Ctatt.Route {
			fmt.Printf("  Route: %s has %d trains\n", r.Rt, len(r.Position))
		}
		return
	}

	fmt.Printf("As of %s\n\n", locations.Ctatt.Tmst.Format(time.Kitchen))

	for _, route := range locations.Ctatt.Route {
		fmt.Printf("Route: %s (%d trains)\n", route.Rt, len(route.Position))

		for _, train := range route.Position {
			flags := ""
			if train.IsDly {
				flags += " DELAYED"
			}
			if train.IsSch {
				flags += " SCHEDULED"
			}

			fmt.Printf("  Run %s | (% .5f, % .5f) | hdg %d | next: %s (%s)%s\n",
				train.RunNumber,
				train.Lat, train.Lon,
				train.Heading,
				train.NextStpNm, getStopName(train.NextStpId),
				flags)
		}
	}
}

func displayBlueLineArrivals(tracker traintracker.TrainTracker) {
	fmt.Println("=== Blue Line Arrivals ===")

	blueLineStations := map[string]string{
		"30062": "Division",
		"30100": "O'Hare",
		"40380": "Clark/Lake",
		"40410": "Jefferson Park",
		"40510": "Rosemont",
		"40580": "Logan Square",
	}

	stationId := "30062"
	stationName := blueLineStations[stationId]

	fmt.Printf("Fetching Blue Line arrivals for %s (stop: %s)\n", stationName, stationId)

	arrivals, err := tracker.Arrivals(traintracker.ArrivalsProps{
		StpId: stationId,
		Rt:    "Blue",
		Max:   5,
	})

	if err != nil {
		fmt.Printf("Error fetching arrivals: %v\n", err)
		return
	}

	if arrivals.Ctatt.ErrCd != "0" {
		fmt.Printf("API Error (%s): %s\n", arrivals.Ctatt.ErrCd, arrivals.Ctatt.ErrNm)
		return
	}

	if arrivals.Ctatt.Eta == nil {
		fmt.Println("No ETA data returned")
		return
	}

	fmt.Printf("As of %s\n\n", arrivals.Ctatt.Tmst.Format(time.Kitchen))

	if len(arrivals.Ctatt.Eta) == 0 {
		fmt.Println("No arrivals")
		return
	}

	for _, eta := range arrivals.Ctatt.Eta {
		flags := ""

		if eta.IsSch {
			flags += " [SCHEDULED]"
		}

		if eta.IsDly {
			flags += " [DELAYED]"
		}

		if eta.IsFlt {
			flags += " [FAULT]"
		}

		if eta.IsApp {
			fmt.Printf("  Run %s -> %s | DUE NOW | (% .5f, % .5f)\n",
				eta.Rn, eta.DestNm, *eta.Lat, *eta.Lon)
		} else {
			waitMins := math.Round(time.Until(eta.ArrT).Minutes())
			fmt.Printf("  Run %s -> %s | in %.0f mins%s | (% .5f, % .5f)\n",
				eta.Rn, eta.DestNm, waitMins, flags, *eta.Lat, *eta.Lon)
		}
	}
}
