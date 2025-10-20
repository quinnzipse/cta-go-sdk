// Package traintracker is a good package
package traintracker

import (
	"fmt"
	"io"
	"log"
	"net/http"
	nurl "net/url"
	"strconv"
)

type ArrivalsProps struct {
	MapId string
	StpId string
	Max   uint16
	Rt    string
}

const path = "/api/1.0/ttarrivals.aspx"

// Arrivals will fetch the train arrivals from the cta api by
// route & (station or stop).
// 30000-39999 = Train stops
// 40000-49999 = Train stations (parent stops)
// routes must be Red | Blue | Brn | G | Org | P | Pink | Y

func (tt TrainTracker) Arrivals(props ArrivalsProps) error {

	if err := sanityCheck(props); err != nil {
		return err
	}

	url, err := generateArrivalsUrl(props, tt.key)
	if err != nil {
		return err
	}

	resp, err := http.Get(url.String())
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(body))

	return nil
}

func sanityCheck(props ArrivalsProps) error {
	// Do sanity checks here.
	if props.MapId == "" && props.StpId == "" {
		return fmt.Errorf("must specify either MapId or StpId")
	}

	if props.MapId != "" {
		mapId, err := strconv.ParseUint(props.MapId, 10, 16)
		if err != nil {
			return fmt.Errorf("error converting mapId to uint16\n%w", err)
		}

		if mapId < 40000 && 50000 <= mapId {
			return fmt.Errorf("invalid mapId, must satisfy 40000 <= mapId < 50000 -- got %d", mapId)
		}
	} else {
		stpId, err := strconv.ParseUint(props.StpId, 10, 16)
		if err != nil {
			return fmt.Errorf("error converting StpId to uint16\n%w", err)
		}

		if stpId < 30000 && 40000 <= stpId {
			return fmt.Errorf("invalid StpId, must satisfy 30000 <= StpId < 40000 -- got %d", stpId)
		}
	}

	switch props.Rt {
	case "":
	case "Red":
	case "Blue":
	case "Brn":
	case "G":
	case "Org":
	case "P":
	case "Pink":
	case "Y":
	default:
		return fmt.Errorf("invalid Rt, must be one of: \"\", Red, Blue, Brn, G, Org, P, Pink, Y -- got %s", props.Rt)
	}

	return nil
}

func generateArrivalsUrl(props ArrivalsProps, key string) (*nurl.URL, error) {
	fullUrl := CtaBaseUrl + path

	// Parse the base URL first
	u, err := nurl.Parse(fullUrl)
	if err != nil {
		return nil, err
	}

	params := nurl.Values{}

	params.Add("outputType", "JSON")

	if props.MapId != "" {
		params.Add("mapid", string(props.MapId))
	}

	if props.StpId != "" {
		params.Add("stpid", props.StpId)
	}

	if props.Max != 0 {
		max := strconv.FormatUint(uint64(props.Max), 10)

		params.Add("max", max)
	}

	if props.Rt != "" {
		params.Add("rt", props.Rt)
	}

	params.Add("key", key)

	u.RawQuery = params.Encode()

	return u, nil
}
