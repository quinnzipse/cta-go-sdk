//go:build live
// +build live

package traintracker_test

import (
	"os"
	"testing"

	"git.zipse.cloud/zippy/cta-go-sdk/traintracker"
)

func TestLocationsToFollow_Live(t *testing.T) {
	key := os.Getenv("CTA_API_KEY")
	if key == "" {
		t.Skip("CTA_API_KEY not set")
	}

	tracker := traintracker.NewTrainTracker(traintracker.TrainTrackerProps{Key: key})

	locs, err := tracker.Locations(traintracker.LocationsProps{
		Rt: "Blue",
	})
	if err != nil {
		t.Fatalf("Locations API failed: %v", err)
	}
	if locs.Ctatt.ErrCd != "0" {
		t.Fatalf("Locations API error %s: %s", locs.Ctatt.ErrCd, locs.Ctatt.ErrNm)
	}

	var runNum string
	for _, route := range locs.Ctatt.Route {
		for _, train := range route.Position {
			runNum = train.RunNumber
			break
		}
		if runNum != "" {
			break
		}
	}

	if runNum == "" {
		t.Skip("no Blue line trains available at this time")
	}

	t.Logf("Found Blue Line run: %s", runNum)

	follow, err := tracker.Follow(traintracker.FollowProps{
		RunNumber: runNum,
	})
	if err != nil {
		t.Fatalf("Follow API failed: %v", err)
	}
	if follow == nil {
		t.Fatal("Follow response is nil")
	}
	if follow.Ctatt.ErrCd != "0" {
		t.Fatalf("Follow API error %s: %s", follow.Ctatt.ErrCd, follow.Ctatt.ErrNm)
	}
	if len(follow.Ctatt.Eta) == 0 {
		t.Log("no stop predictions returned (valid for some runs)")
	} else {
		t.Logf("Follow returned %d stops", len(follow.Ctatt.Eta))
		for i, eta := range follow.Ctatt.Eta {
			t.Logf("  stop %d: %s -> %s", i+1, eta.StaNm, eta.DestNm)
		}
	}

	if follow.Ctatt.Position != nil {
		t.Logf("Train position: (%.5f, %.5f) heading %d",
			*follow.Ctatt.Position.Lat,
			*follow.Ctatt.Position.Lon,
			*follow.Ctatt.Position.Heading)
	}
}
