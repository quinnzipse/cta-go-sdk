//go:build live
// +build live

package traintracker_test

import (
	"os"
	"testing"

	"git.zipse.cloud/zippy/cta-go-sdk/traintracker"
)

func TestArrivals_Live(t *testing.T) {
	key := os.Getenv("CTA_API_KEY")
	if key == "" {
		t.Skip("CTA_API_KEY not set")
	}

	tracker := traintracker.NewTrainTracker(traintracker.TrainTrackerProps{Key: key})

	resp, err := tracker.Arrivals(traintracker.ArrivalsProps{
		MapId: "40380",
	})
	if err != nil {
		t.Fatalf("API call failed: %v", err)
	}
	if resp == nil {
		t.Fatal("response is nil")
	}
	if resp.Ctatt.ErrCd != "0" {
		t.Fatalf("API error %s: %s", resp.Ctatt.ErrCd, resp.Ctatt.ErrNm)
	}
	if resp.Ctatt.Eta == nil {
		t.Fatal("eta is nil")
	}
	if len(resp.Ctatt.Eta) == 0 {
		t.Log("no arrivals returned (may be valid at off hours)")
	}
}
