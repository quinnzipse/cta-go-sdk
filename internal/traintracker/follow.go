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

	var positions []FollowPosition
	for i, pos := range r.Ctatt.Route.Position {
		var prdt *time.Time
		if pos.Prdt != "" && pos.Prdt != "0" && pos.Prdt != "?" {
			parsed, err := time.ParseInLocation("2006-01-02T15:04:05", pos.Prdt, timeLoc)
			if err != nil {
				return nil, fmt.Errorf("Error on position[%d]: %w", i, err)
			}
			prdt = &parsed
		}

		var arrT *time.Time
		if pos.ArrT != "" && pos.ArrT != "0" && pos.ArrT != "?" {
			parsed, err := time.ParseInLocation("2006-01-02T15:04:05", pos.ArrT, timeLoc)
			if err != nil {
				return nil, fmt.Errorf("Error on position[%d]: %w", i, err)
			}
			arrT = &parsed
		}

		positions = append(positions, FollowPosition{
			StpId: pos.StpId,
			StpNm: pos.StpNm,
			StpDe: pos.StpDe,
			Prdt:  prdt,
			ArrT:  arrT,
			IsApp: pos.IsApp == "1",
			IsSch: pos.IsSch == "1",
			IsDly: pos.IsDly == "1",
			IsFlt: pos.IsFlt == "1",
		})
	}

	return &FollowResponse{
		Ctatt: FollowCtatt{
			Tmst:    tmst,
			ErrCd:   r.Ctatt.ErrCd,
			ErrNm:   r.Ctatt.ErrNm,
			Run:    r.Ctatt.Run,
			Rt:     r.Ctatt.Rt,
			DestNm: r.Ctatt.DestNm,
			Route:  positions,
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