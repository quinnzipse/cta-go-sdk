package main

import (
	"fmt"

	"git.quinnzipse.dev/cta-go-sdk/internal/traintracker"
)

func main() {
	fmt.Println("Grabbing info for Division (sid: 40320)")

	_, err := traintracker.GetTrainArrivalsByStation(traintracker.GetTrainArrivalsByStationProps{
		Station: traintracker.Station{Id: 40320, Name: "Division"},
	})

	if err != nil {
		fmt.Printf("Hit an error\n\n%v", err)
	}
}
