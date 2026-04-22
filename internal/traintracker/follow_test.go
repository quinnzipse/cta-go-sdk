package traintracker_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"git.quinnzipse.dev/cta-go-sdk/internal/traintracker"
)

func TestFollow_Validation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input  traintracker.FollowProps
		wantErr bool
	}{
		{
			name:    "empty RunNumber",
			input:  traintracker.FollowProps{},
			wantErr: true,
		},
		{
			name: "valid RunNumber",
			input: traintracker.FollowProps{
				RunNumber: "123",
			},
			wantErr: false,
		},
		{
			name: "RunNumber with leading zeros",
			input: traintracker.FollowProps{
				RunNumber: "001",
			},
			wantErr: false,
		},
		{
			name: "RunNumber with max value",
			input: traintracker.FollowProps{
				RunNumber: "999",
			},
			wantErr: false,
		},
		{
			name: "RunNumber non-numeric",
			input: traintracker.FollowProps{
				RunNumber: "abc",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"err": "Invalid"}`))
			}))
			defer server.Close()

			tracker := traintracker.NewTrainTracker(traintracker.TrainTrackerProps{
				Key:        "test-key",
				HttpClient: server.Client(),
				BaseUrl:    server.URL,
			})

			_, err := tracker.Follow(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}

			if err != nil {
				errStr := err.Error()
				if contains(errStr, "RunNumber") {
					t.Errorf("unexpected validation error: %v", err)
				}
			}
		})
	}
}

func TestFollow_Mocking(t *testing.T) {
	t.Parallel()

	mockResponse := `{
		"ctatt": {
			"tmst": "2024-01-15T10:30:00",
			"errCd": "0",
			"errNm": "",
			"run": "401",
			"rt": "Red",
			"destNm": "Howard",
			"route": {
				"position": [
					{
						"stpId": "40090",
						"stpNm": "Lake",
						"stpDe": "Service toward Howard",
						"prdt": "2024-01-15T10:32:00",
						"arrT": "2024-01-15T10:32:00",
						"isApp": "1",
						"isSch": "0",
						"isDly": "0",
						"isFlt": "0"
					},
					{
						"stpId": "40100",
						"stpNm": "Chicago",
						"stpDe": "Service toward Howard",
						"prdt": "2024-01-15T10:35:00",
						"arrT": "2024-01-15T10:35:00",
						"isApp": "0",
						"isSch": "0",
						"isDly": "0",
						"isFlt": "0"
					}
				]
			}
		}
	}`

	tests := []struct {
		name         string
		mockResponse string
		props       traintracker.FollowProps
		wantErr     bool
		checkFn     func(*testing.T, *traintracker.FollowResponse)
	}{
		{
			name:         "valid response with two stops",
			mockResponse: mockResponse,
			props: traintracker.FollowProps{
				RunNumber: "401",
			},
			wantErr: false,
			checkFn: func(t *testing.T, resp *traintracker.FollowResponse) {
				if resp.Ctatt.Run != "401" {
					t.Errorf("expected run 401, got %s", resp.Ctatt.Run)
				}
				if resp.Ctatt.Rt != "Red" {
					t.Errorf("expected route Red, got %s", resp.Ctatt.Rt)
				}
				if resp.Ctatt.DestNm != "Howard" {
					t.Errorf("expected destination Howard, got %s", resp.Ctatt.DestNm)
				}
				if len(resp.Ctatt.Route) != 2 {
					t.Errorf("expected 2 stops, got %d", len(resp.Ctatt.Route))
				}
				if resp.Ctatt.Route[0].StpNm != "Lake" {
					t.Errorf("expected first stop Lake, got %s", resp.Ctatt.Route[0].StpNm)
				}
			},
		},
		{
			name:         "error response from API",
			mockResponse: `{"err": "Invalid Run Number"}`,
			props: traintracker.FollowProps{
				RunNumber: "999",
			},
			wantErr: true,
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

			tracker := traintracker.NewTrainTracker(traintracker.TrainTrackerProps{
				Key:        "test-key",
				HttpClient: server.Client(),
				BaseUrl:    server.URL,
			})

			resp, err := tracker.Follow(tt.props)

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