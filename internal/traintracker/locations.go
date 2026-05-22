package traintracker

import (
	"fmt"
	"net/http"
	nurl "net/url"
	"strconv"
	"time"
)

type LocationsProps struct {
	Rt  string
	Max uint16
}

const positionsPath = "/api/1.0/ttpositions.aspx"

func (tt TrainTracker) RawLocations(props LocationsProps) (*PositionsApiResponse, error) {
	if err := sanityCheckLocations(props); err != nil {
		return nil, err
	}

	url, err := generateLocationsUrl(props, tt.key, tt.baseUrl)
	if err != nil {
		return nil, err
	}

	resp, err := tt.httpClient.Get(url.String())
	if err != nil {
		return nil, err
	}

	positionsResp, err := parseRawPositionsApiResponse(resp)
	if err != nil {
		return nil, err
	}

	return positionsResp, nil
}

func (tt TrainTracker) Locations(props LocationsProps) (*PositionsResponse, error) {
	res, err := tt.RawLocations(props)
	if err != nil {
		return nil, err
	}

	return res.toPositionsResponse()
}

func parseRawPositionsApiResponse(resp *http.Response) (*PositionsApiResponse, error) {
	if resp == nil {
		return nil, fmt.Errorf("No response was included")
	}

	type ApiErr struct {
		Err string `json:"err"`
	}

	apiRes, apiErr, err := decodeAPIResponse[PositionsApiResponse, ApiErr](resp)
	if err != nil {
		return nil, err
	}

	if apiErr != nil {
		return nil, fmt.Errorf("Error while decoding %v", apiErr)
	}

	return apiRes, nil
}

func (r PositionsApiResponse) toPositionsResponse() (*PositionsResponse, error) {
	timeLoc, err := time.LoadLocation("America/Chicago")
	if err != nil {
		return nil, err
	}

	tmst, err := time.ParseInLocation("2006-01-02T15:04:05", r.Ctatt.Tmst, timeLoc)
	if err != nil {
		return nil, err
	}

	routes := []RoutePosition{}

	for _, route := range r.Ctatt.Route {
		routePos := RoutePosition{
			Rt:       route.Rt,
			Position: []TrainPosition{},
		}

		for i, train := range route.Position {
			lat, err := strconv.ParseFloat(train.Lat, 64)
			if err != nil {
				return nil, fmt.Errorf("Error on route[%s] train[%d]: %w", route.Rt, i, err)
			}

			lon, err := strconv.ParseFloat(train.Lon, 64)
			if err != nil {
				return nil, fmt.Errorf("Error on route[%s] train[%d]: %w", route.Rt, i, err)
			}

			heading, err := strconv.ParseInt(train.Heading, 10, 16)
			if err != nil {
				return nil, fmt.Errorf("Error on route[%s] train[%d]: %w", route.Rt, i, err)
			}

			var nextSrt *time.Time
			if train.NextStpTm != "" && train.NextStpTm != "0" && train.NextStpTm != "?" {
				var nxt time.Time
				nxt, err = time.ParseInLocation("2006-01-02T15:04:05", train.NextStpTm, timeLoc)
				if err != nil {
					return nil, fmt.Errorf("Error on route[%s] train[%d]: %w", route.Rt, i, err)
				}
				nextSrt = &nxt
			}

			var prevSrt *time.Time
			if train.PrevStpTm != "" && train.PrevStpTm != "0" && train.PrevStpTm != "?" {
				var prv time.Time
				prv, err = time.ParseInLocation("2006-01-02T15:04:05", train.PrevStpTm, timeLoc)
				if err != nil {
					return nil, fmt.Errorf("Error on route[%s] train[%d]: %w", route.Rt, i, err)
				}
				prevSrt = &prv
			}

			trDr, err := strconv.ParseInt(train.TrDr, 10, 16)
			if err != nil {
				return nil, fmt.Errorf("Error on route[%s] train[%d]: %w", route.Rt, i, err)
			}

			routePos.Position = append(routePos.Position, TrainPosition{
				TrainId:   train.TrainId,
				RunNumber: train.RunNumber,
				Rt:        train.Rt,
				TrDr:      int16(trDr),
				Lat:       lat,
				Lon:       lon,
				Heading:   int16(heading),
				NextStpId: train.NextStpId,
				NextStpNm: train.NextStpNm,
				NextStpTm: nextSrt,
				PrevStpId: train.PrevStpId,
				PrevStpNm: train.PrevStpNm,
				PrevStpTm: prevSrt,
				IsSch:     train.IsSch == "1",
				IsDly:     train.IsDly == "1",
				IsFlt:     train.IsFlt == "1",
				IsBoe:     train.IsBoe == "1",
				IsDet:     train.IsDet == "1",
				Flags:     train.Flags,
			})
		}

		routes = append(routes, routePos)
	}

	return &PositionsResponse{
		Ctatt: CtattPos{
			Tmst:  tmst,
			ErrCd: r.Ctatt.ErrCd,
			ErrNm: r.Ctatt.ErrNm,
			Route: routes,
		},
	}, nil
}

func sanityCheckLocations(props LocationsProps) error {
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

func generateLocationsUrl(props LocationsProps, key, baseUrl string) (*nurl.URL, error) {
	fullUrl := baseUrl + positionsPath

	u, err := nurl.Parse(fullUrl)
	if err != nil {
		return nil, err
	}

	params := nurl.Values{}

	params.Add("outputType", "JSON")

	if props.Rt != "" {
		params.Add("rt", props.Rt)
	}

	if props.Max != 0 {
		params.Add("max", strconv.FormatUint(uint64(props.Max), 10))
	}

	params.Add("key", key)

	u.RawQuery = params.Encode()

	return u, nil
}
