package traintracker_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"git.quinnzipse.dev/cta-go-sdk/internal/traintracker"
)

func TestLocations_Validation(t *testing.T) {
	t.Parallel()

	tracker := traintracker.NewTrainTracker(traintracker.TrainTrackerProps{
		Key: "test-key",
	})

	tests := []struct {
		name    string
		input   traintracker.LocationsProps
		wantErr bool
	}{
		{
			name:    "empty input - valid (returns all routes)",
			input:   traintracker.LocationsProps{},
			wantErr: false,
		},
		{
			name: "valid route - Red",
			input: traintracker.LocationsProps{
				Rt: "Red",
			},
			wantErr: false,
		},
		{
			name: "valid route - Blue",
			input: traintracker.LocationsProps{
				Rt: "Blue",
			},
			wantErr: false,
		},
		{
			name: "valid route - Brn",
			input: traintracker.LocationsProps{
				Rt: "Brn",
			},
			wantErr: false,
		},
		{
			name: "valid route - G",
			input: traintracker.LocationsProps{
				Rt: "G",
			},
			wantErr: false,
		},
		{
			name: "valid route - Org",
			input: traintracker.LocationsProps{
				Rt: "Org",
			},
			wantErr: false,
		},
		{
			name: "valid route - P",
			input: traintracker.LocationsProps{
				Rt: "P",
			},
			wantErr: false,
		},
		{
			name: "valid route - Pink",
			input: traintracker.LocationsProps{
				Rt: "Pink",
			},
			wantErr: false,
		},
		{
			name: "valid route - Y",
			input: traintracker.LocationsProps{
				Rt: "Y",
			},
			wantErr: false,
		},
		{
			name: "valid route with max",
			input: traintracker.LocationsProps{
				Rt:  "Red",
				Max: 5,
			},
			wantErr: false,
		},
		{
			name: "invalid route",
			input: traintracker.LocationsProps{
				Rt: "InvalidRoute",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := tracker.Locations(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil (input: %+v)", tt.input)
				}
				return
			}

			if err != nil {
				errStr := err.Error()
				validationErrors := []string{
					"invalid Rt",
				}

				for _, ve := range validationErrors {
					if contains(errStr, ve) {
						t.Errorf("unexpected validation error: %v (input: %+v)", err, tt.input)
						return
					}
				}
			}
		})
	}
}

func TestLocations_UrlGeneration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		props      traintracker.LocationsProps
		key        string
		wantErr    bool
		wantParams []string
	}{
		{
			name:    "empty props",
			props:  traintracker.LocationsProps{},
			key:    "test-key",
			wantErr: false,
			wantParams: []string{
				"outputType=JSON",
				"key=test-key",
			},
		},
		{
			name: "with route",
			props: traintracker.LocationsProps{
				Rt: "Red",
			},
			key:     "test-key",
			wantErr: false,
			wantParams: []string{
				"rt=Red",
				"outputType=JSON",
				"key=test-key",
			},
		},
		{
			name: "with max",
			props: traintracker.LocationsProps{
				Max: 10,
			},
			key:     "test-key",
			wantErr: false,
			wantParams: []string{
				"max=10",
				"outputType=JSON",
				"key=test-key",
			},
		},
		{
			name: "full params",
			props: traintracker.LocationsProps{
				Rt:  "Blue",
				Max: 5,
			},
			key:     "my-api-key",
			wantErr: false,
			wantParams: []string{
				"rt=Blue",
				"max=5",
				"outputType=JSON",
				"key=my-api-key",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tracker := traintracker.NewTrainTracker(traintracker.TrainTrackerProps{
				Key: tt.key,
			})

			_, err := tracker.RawLocations(tt.props)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}

			if err != nil {
				return
			}
		})
	}
}

func TestLocations_Mocking(t *testing.T) {
	t.Parallel()

	mockResponse := `
{ "ctatt": { "tmst":"2024-01-15T10:30:00", "errCd":"0", "errNm":"",
  "route": [
   { "rt":"Red", "position": [
    { "trainId":"401", "rn":"401", "rt":"Red", "trDr":"1", "lat":"41.8781", "lon":"-87.6298", "heading":"180",
      "nextStpId":"40100", "nextStpNm":"Chicago", "nextStpTm":"2024-01-15T10:35:00",
      "prevStpId":"40090", "prevStpNm":"Lake", "prevStpTm":"2024-01-15T10:25:00",
      "isSch":"0","isDly":"0","isFlt":"0","isBoe":"0","isDet":"0","flags":[] }
   ] }
  ]
 }
}`

	tests := []struct {
		name         string
		mockResponse string
		props        traintracker.LocationsProps
		wantRoutes   int
		wantErr      bool
		checkFn      func(*testing.T, *traintracker.PositionsResponse)
	}{
		{
			name:         "valid response with two routes",
			mockResponse: mockResponse,
props:        traintracker.LocationsProps{},
		wantRoutes:   1,
		wantErr:     false,
		checkFn: func(t *testing.T, resp *traintracker.PositionsResponse) {
			if len(resp.Ctatt.Route) != 1 {
				t.Errorf("expected 1 route, got %d", len(resp.Ctatt.Route))
			}
			if resp.Ctatt.Route[0].Rt != "Red" {
				t.Errorf("expected first route Red, got %s", resp.Ctatt.Route[0].Rt)
			}
			if resp.Ctatt.Route[0].Position[0].Lat != 41.8781 {
				t.Errorf("expected lat 41.8781, got %f", resp.Ctatt.Route[0].Position[0].Lat)
			}
		},
	},
		{
			name:         "error response",
			mockResponse: `{"err": "Invalid API Key"}`,
			props:       traintracker.LocationsProps{},
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(tt.mockResponse))
			}))
			defer server.Close()

			client := server.Client()

			tracker := traintracker.NewTrainTracker(traintracker.TrainTrackerProps{
				Key:        "test-key",
				HttpClient: client,
				BaseUrl:   server.URL,
			})

			resp, err := tracker.Locations(tt.props)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.checkFn != nil {
				tt.checkFn(t, resp)
			}
		})
	}
}