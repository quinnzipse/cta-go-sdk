// Package cta is a good package
package cta

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	nurl "net/url"
	"os"
	"strconv"

	"github.com/google/uuid"
)

type GetTrainArrivalsByStationProps struct {
	Station   Station
	Direction string
}

// GetTrainArrivalsByStation will fetch the train arrivals from the cta api by
// station and optionally direction of travel.
// 0-29999 = Bus stops
// 30000-39999 = Train stops
// 40000-49999 = Train stations (parent stops)
func GetTrainArrivalsByStation(props GetTrainArrivalsByStationProps) ([]TrainArrival, error) {

	if props.Station.Id == 0 {
		return nil, errors.New("station id is required")
	}

	if props.Station.Id < 30000 || props.Station.Id > 49999 {
		return nil, errors.New("station id (train) is invalid")
	}

	fmt.Println("Calling arrivals api")

	mapId := strconv.FormatUint(uint64(props.Station.Id), 10)

	err := callArrivalsApi(arrivalsApiProps{mapId: mapId, apiKey: uuid.New().String()})
	if err != nil {
		return nil, err
	}

	return []TrainArrival{}, nil
}

type arrivalsApiProps struct {
	mapId  string
	stpId  string
	apiKey string
}

func generateArrivalsUrl(props arrivalsApiProps) (*nurl.URL, error) {
	baseURL := "http://lapi.transitchicago.com/api/1.0/ttarrivals.aspx"

	// Parse the base URL first
	u, err := nurl.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	if len(props.mapId) != 5 {
		return nil, errors.New("mapId is invalid")
	}

	// Add query parameters
	params := nurl.Values{}

	params.Add("outputType", "JSON")

	if props.mapId != "" {
		params.Add("mapid", string(props.mapId))
	}

	if props.stpId != "" {
		params.Add("stpid", props.stpId)
	}

	params.Add("key", props.apiKey)

	// Attach encoded params to the URL
	u.RawQuery = params.Encode()

	_, err = fmt.Println(u.String())
	if err != nil {
		return nil, err
	}

	return u, nil
}

func callArrivalsApi(props arrivalsApiProps) error {

	apiKey := os.Getenv("CTA_API_KEY")
	if apiKey == "" {
		return errors.New("environment error: CTA_API_KEY is not set")
	}

	url, err := generateArrivalsUrl(arrivalsApiProps{
		mapId:  props.mapId,
		apiKey: apiKey,
	})
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
