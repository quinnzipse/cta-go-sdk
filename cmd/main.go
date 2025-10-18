package main

import (
	"fmt"

	"git.quinnzipse.dev/cta-tracker/internal/cta"
)

func main() {
	fmt.Println("Grabbing info for Division (sid: 40320)")

	_, err := cta.GetTrainArrivalsByStation(cta.GetTrainArrivalsByStationProps{
		Station: cta.Station{Id: 40320, Name: "Division"},
	})
	if err != nil {
		fmt.Printf("Hit an error\n\n%v", err)
	}
}
