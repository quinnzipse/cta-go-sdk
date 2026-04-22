// Package traintracker is a good package
package traintracker

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	nurl "net/url"
	"strconv"
	"time"
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

func (tt TrainTracker) RawArrivals(props ArrivalsProps) (*ArrivalApiResponse, error) {

	if err := sanityCheck(props); err != nil {
		return nil, err
	}

	url, err := generateArrivalsUrl(props, tt.key, tt.baseUrl)
	if err != nil {
		return nil, err
	}

	resp, err := tt.httpClient.Get(url.String())
	if err != nil {
		return nil, err
	}

	arrivalApiResp, err := parseRawArrivalApiResponse(resp)
	if err != nil {
		return nil, err
	}

	return arrivalApiResp, nil
}

func (tt TrainTracker) Arrivals(props ArrivalsProps) (*ArrivalResponse, error) {
	res, err := tt.RawArrivals(props)
	if err != nil {
		return nil, err
	}

	fmt.Printf("RawArrivals -> res:\n%+v\n", res)

	return res.toArrivalResponse()
}

func decodeAPIResponse[T any, E any](resp *http.Response) (*T, *E, error) {
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var apiErr E
		if err := json.Unmarshal(body, &apiErr); err != nil {
			return nil, nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
		}
		return nil, &apiErr, nil
	}

	var result T
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	return &result, nil, nil
}

func parseRawArrivalApiResponse(resp *http.Response) (*ArrivalApiResponse, error) {
	if resp == nil {
		return nil, fmt.Errorf("No response was included")
	}

	type ApiErr struct {
		Err string `json:"err"`
	}

	apiRes, apiErr, err := decodeAPIResponse[ArrivalApiResponse, ApiErr](resp)
	if err != nil {
		return nil, err
	}

	if apiErr != nil {
		return nil, fmt.Errorf("Error while decoding %v", apiErr)
	}

	return apiRes, nil
}

func (r ArrivalApiResponse) toArrivalResponse() (*ArrivalResponse, error) {
	timeLoc, err := time.LoadLocation("America/Chicago")
	if err != nil {
		return nil, err
	}

	tmst, err := time.ParseInLocation("2006-01-02T15:04:05", r.Ctatt.Tmst, timeLoc)
	if err != nil {
		return nil, err
	}

	output := ArrivalResponse{
		Ctatt: Ctatt{
			Tmst:  tmst,
			ErrCd: r.Ctatt.ErrCd,
			ErrNm: r.Ctatt.ErrNm,
			Eta:   []Eta{},
		},
	}

	for i, eta := range r.Ctatt.Eta {

		var lat *float64
		if eta.Lat != "" {
			parsed, err := strconv.ParseFloat(eta.Lat, 64)
			if err != nil {
				return nil, fmt.Errorf("Error on eta[%d]\n\n %w", i, err)
			}
			lat = &parsed
		}

		var lon *float64
		if eta.Lon != "" {
			parsed, err := strconv.ParseFloat(eta.Lon, 64)
			if err != nil {
				return nil, fmt.Errorf("Error on eta[%d]\n\n %w", i, err)
			}
			lon = &parsed
		}

		var heading *int16
		if eta.Heading != "" {
			parsed, err := strconv.ParseInt(eta.Heading, 10, 16)
			if err != nil {
				return nil, fmt.Errorf("Error on eta[%d]\n\n %w", i, err)
			}
			h := int16(parsed)
			heading = &h
		}

		trDr, err := strconv.ParseInt(eta.TrDr, 10, 16)
		if err != nil {
			return nil, fmt.Errorf("Error on eta[%d]\n\n %w", i, err)
		}

		prdt, err := time.ParseInLocation("2006-01-02T15:04:05", eta.Prdt, timeLoc)
		if err != nil {
			return nil, fmt.Errorf("Error on eta[%d]\n\n %w", i, err)
		}

		arrT, err := time.ParseInLocation("2006-01-02T15:04:05", eta.ArrT, timeLoc)
		if err != nil {
			return nil, fmt.Errorf("Error on eta[%d]\n\n %w", i, err)
		}

		output.Ctatt.Eta = append(output.Ctatt.Eta, Eta{
			StaId:   eta.StaId,
			StpId:   eta.StpId,
			StaNm:   eta.StaNm,
			StpDe:   eta.StpDe,
			Rn:      eta.Rn,
			Rt:      eta.Rt,
			DestSt:  eta.DestSt,
			DestNm:  eta.DestNm,
			TrDr:    int16(trDr),
			Prdt:    prdt,
			ArrT:    arrT,
			IsApp:   eta.IsApp == "1",
			IsSch:   eta.IsSch == "1",
			IsDly:   eta.IsDly == "1",
			IsFlt:   eta.IsFlt == "1",
			Flags:   eta.Flags,
			Lat:     lat,
			Lon:     lon,
			Heading: heading,
		})

	}

	return &output, nil
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

func generateArrivalsUrl(props ArrivalsProps, key, baseUrl string) (*nurl.URL, error) {
	fullUrl := baseUrl + path

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
