# cta-go-sdk

A Go SDK for the [Chicago Transit Authority (CTA) Train Tracker API](https://www.transitchicago.com/developers/traintracker/).

It wraps the CTA's train endpoints so you can get real-time train data in your Go programs without building your own HTTP calls and JSON parsing.

## What you can do with it

- Get live **train locations** on a route (or all routes).
- Get **arrival predictions** for a station or stop.
- Follow a single **train** by its run number.

## Before you start

You need a CTA API key. Get one for free from the [CTA Developers page](https://www.transitchicago.com/developers/).

You also need Go installed (1.22 or newer).

## Installation

Add the package to your project:

```bash
go get github.com/quinnzipse/cta-go-sdk/traintracker
```

## Basic example

```go
package main

import (
	"fmt"
	"os"

	"github.com/quinnzipse/cta-go-sdk/traintracker"
)

func main() {
	key := os.Getenv("CTA_API_KEY")

	tracker := traintracker.NewTrainTracker(traintracker.TrainTrackerProps{
		Key: key,
	})

	// Get arrivals for the Division Blue Line stop.
	arrivals, err := tracker.Arrivals(traintracker.ArrivalsProps{
		StpId: "30062",
		Rt:    "Blue",
		Max:   5,
	})
	if err != nil {
		panic(err)
	}

	for _, eta := range arrivals.Ctatt.Eta {
		fmt.Printf("Run %s -> %s arriving at %s\n", eta.Rn, eta.DestNm, eta.ArrT)
	}
}
```

Run it with your key in the environment:

```bash
CTA_API_KEY=your_key_here go run .
```

## More examples

Get the live locations of every Blue Line train:

```go
locations, err := tracker.Locations(traintracker.LocationsProps{
	Rt: "Blue",
})
```

Follow a specific train by its run number:

```go
train, err := tracker.Follow(traintracker.FollowProps{
	RunNumber: "430",
})
```

## Route names

The API accepts these route codes: `Red`, `Blue`, `Brn`, `G`, `Org`, `P`, `Pink`, and `Y`. Pass an empty string to get all routes.

Stop IDs are 5-digit numbers. Station IDs are `40000-49999` and individual stop IDs are `30000-39999`.

## Running the demo

There is a small demo in `cmd/main.go` that prints Blue Line locations and arrivals. It needs your key:

```bash
CTA_API_KEY=your_key_here go run ./cmd
```

## Testing

Run the unit tests with:

```bash
go test ./...
```

Some tests are integration tests that call the live CTA API and are skipped unless the right build tag and key are set. Check the test files for details.

## License

See the [LICENSE](LICENSE) file.
