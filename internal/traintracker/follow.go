package traintracker

import (
	"fmt"
	"net/http"
	nurl "net/url"
	"strconv"
	"time"
)

func (tt TrainTracker) RawFollow(props FollowProps) (*FollowApiResponse, error) {
	if err := sanityCheckFollow(props); err != nil {
		return nil, err
	}

	url, err := generateFollowUrl(props, tt.key, tt.baseUrl)
	if err != nil {
		return nil, err
	}

	resp, err := tt.httpClient.Get(url.String())
	if err != nil {
		return nil, err
	}

	followResp, err := parseRawFollowApiResponse(resp)
	if err != nil {
		return nil, err
	}

	return followResp, nil
}

func (tt TrainTracker) Follow(props FollowProps) (*FollowResponse, error) {
	res, err := tt.RawFollow(props)
	if err != nil {
		return nil, err
	}

	return res.toFollowResponse()
}

func parseRawFollowApiResponse(resp *http.Response) (*FollowApiResponse, error) {
	if resp == nil {
		return nil, fmt.Errorf("No response was included")
	}

	type ApiErr struct {
		Err string `json:"err"`
	}

	apiRes, apiErr, err := decodeAPIResponse[FollowApiResponse, ApiErr](resp)
	if err != nil {
		return nil, err
	}

	if apiErr != nil {
		return nil, fmt.Errorf("Error while decoding %v", apiErr)
	}

	return apiRes, nil
}

func (r FollowApiResponse) toFollowResponse() (*FollowResponse, error) {
	timeLoc, err := time.LoadLocation("America/Chicago")
	if err != nil {
		return nil, err
	}

	tmst, err := time.ParseInLocation("2006-01-02T15:04:05", r.Ctatt.Tmst, timeLoc)
	if err != nil {
		return nil, err
	}

	var trainPos *FollowTrainPosition
	if r.Ctatt.Position.Lat != "" {
		lat, err := strconv.ParseFloat(r.Ctatt.Position.Lat, 64)
		if err != nil {
			return nil, fmt.Errorf("Error parsing lat: %w", err)
		}
		lon, err := strconv.ParseFloat(r.Ctatt.Position.Lon, 64)
		if err != nil {
			return nil, fmt.Errorf("Error parsing lon: %w", err)
		}
		heading, err := strconv.ParseInt(r.Ctatt.Position.Heading, 10, 16)
		if err != nil {
			return nil, fmt.Errorf("Error parsing heading: %w", err)
		}
		h := int16(heading)
		trainPos = &FollowTrainPosition{
			Lat:     &lat,
			Lon:     &lon,
			Heading: &h,
		}
	}

	stops := make([]FollowStopPrediction, 0, len(r.Ctatt.Eta))
	for i, eta := range r.Ctatt.Eta {
		var prdt *time.Time
		if eta.Prdt != "" && eta.Prdt != "0" && eta.Prdt != "?" {
			parsed, err := time.ParseInLocation("2006-01-02T15:04:05", eta.Prdt, timeLoc)
			if err != nil {
				return nil, fmt.Errorf("Error on eta[%d]: %w", i, err)
			}
			prdt = &parsed
		}

		var arrT *time.Time
		if eta.ArrT != "" && eta.ArrT != "0" && eta.ArrT != "?" {
			parsed, err := time.ParseInLocation("2006-01-02T15:04:05", eta.ArrT, timeLoc)
			if err != nil {
				return nil, fmt.Errorf("Error on eta[%d]: %w", i, err)
			}
			arrT = &parsed
		}

		trDr, err := strconv.ParseInt(eta.TrDr, 10, 16)
		if err != nil {
			return nil, fmt.Errorf("Error on eta[%d] trDr: %w", i, err)
		}

		stops = append(stops, FollowStopPrediction{
			StaId:  eta.StaId,
			StpId:  eta.StpId,
			StaNm:  eta.StaNm,
			StpDe:  eta.StpDe,
			Rn:     eta.Rn,
			Rt:     eta.Rt,
			DestSt: eta.DestSt,
			DestNm: eta.DestNm,
			TrDr:   int16(trDr),
			Prdt:   prdt,
			ArrT:   arrT,
			IsApp:  eta.IsApp == "1",
			IsSch:  eta.IsSch == "1",
			IsDly:  eta.IsDly == "1",
			IsFlt:  eta.IsFlt == "1",
		})
	}

	return &FollowResponse{
		Ctatt: FollowCtatt{
			Tmst:     tmst,
			ErrCd:    r.Ctatt.ErrCd,
			ErrNm:    r.Ctatt.ErrNm,
			Position: trainPos,
			Eta:      stops,
		},
	}, nil
}

func sanityCheckFollow(props FollowProps) error {
	if props.RunNumber == "" {
		return fmt.Errorf("RunNumber is required")
	}

	_, err := strconv.ParseUint(props.RunNumber, 10, 16)
	if err != nil {
		return fmt.Errorf("error converting RunNumber to uint16: %w", err)
	}

	return nil
}

func generateFollowUrl(props FollowProps, key, baseUrl string) (*nurl.URL, error) {
	fullUrl := baseUrl + followPath

	u, err := nurl.Parse(fullUrl)
	if err != nil {
		return nil, err
	}

	params := nurl.Values{}

	params.Add("outputType", "JSON")
	params.Add("runnumber", props.RunNumber)
	params.Add("key", key)

	u.RawQuery = params.Encode()

	return u, nil
}
