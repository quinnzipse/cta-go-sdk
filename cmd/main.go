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

func main() {
	fmt.Println("Initializing Train Tracker")

	key := os.Getenv("CTA_API_KEY")
	if strings.TrimSpace(key) == "" {
		panic(errors.New("CTA_API_KEY is required"))
	}

	tracker := traintracker.NewTrainTracker(traintracker.TrainTrackerProps{
		Key: key,
	})

	fmt.Println("Grabbing info for Division (sid: 40320)")

	arrivals, err := tracker.Arrivals(traintracker.ArrivalsProps{
		MapId: "40590",
	})

	if err != nil {
		fmt.Printf("Hit an error\n\n%v", err)
	}

	fmt.Println("Retrieved info")
	if arrivals.Ctatt.ErrCd != "0" {
		fmt.Printf("Error (@%s) %s %s", arrivals.Ctatt.Tmst.Format(time.Kitchen), arrivals.Ctatt.ErrCd, arrivals.Ctatt.ErrNm)
	}

	fmt.Printf("\n----------------\nTrain Schedule (as of %s)\n----------------\n", arrivals.Ctatt.Tmst.Format(time.Kitchen))

	etas := arrivals.Ctatt.Eta

	if len(etas) == 0 {
		fmt.Println("No Arrivals")
		return
	}

	fmt.Printf("Station: %s\n\n", arrivals.Ctatt.Eta[0].StaNm)

	for _, eta := range etas {

		flags := ""

		if eta.IsSch {
			flags += "⏰"
		}

		if eta.IsFlt {
			flags += "!"
		}

		if eta.IsDly {
			flags += "⏳"
		}

		if eta.IsApp {
			fmt.Printf("Run %s\t\t%s    \tdue %s         \t\t(%.5f, %.5f) hdg %d\n", eta.Rn, eta.DestNm, flags, eta.Lat, eta.Lon, eta.Heading)
		} else {
			fmt.Printf("Run %s\t\t%s    \tin %.0f mins %s\t\t(%.5f, %.5f) hdg %d\n", eta.Rn, eta.DestNm, math.Round(eta.ArrT.Sub(time.Now()).Minutes()), flags, eta.Lat, eta.Lon, eta.Heading)
		}
	}

}
