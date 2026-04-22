package traintracker

import "net/http"

const DefaultBaseUrl = "https://lapi.transitchicago.com"

type TrainTrackerProps struct {
	Key           string
	HttpClient   *http.Client
	BaseUrl     string
}

type TrainTracker struct {
	key        string
	httpClient *http.Client
	baseUrl   string
}

func NewTrainTracker(props TrainTrackerProps) TrainTracker {
	client := props.HttpClient
	if client == nil {
		client = http.DefaultClient
	}
	baseUrl := props.BaseUrl
	if baseUrl == "" {
		baseUrl = DefaultBaseUrl
	}
	return TrainTracker{
		key:        props.Key,
		httpClient: client,
		baseUrl:   baseUrl,
	}
}

func (tt TrainTracker) Do(req *http.Request) (*http.Response, error) {
	return tt.httpClient.Do(req)
}
