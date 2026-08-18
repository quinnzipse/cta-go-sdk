package traintracker_test

import (
	"testing"

	"github.com/quinnzipse/cta-go-sdk/traintracker"
)

func TestArrivals_Validation(t *testing.T) {
	t.Parallel()

	tracker := traintracker.NewTrainTracker(traintracker.TrainTrackerProps{
		Key: "test-key",
	})

	tests := []struct {
		name    string
		input   traintracker.ArrivalsProps
		wantErr bool
	}{
		{
			name:    "empty input - no MapId or StpId",
			input:   traintracker.ArrivalsProps{},
			wantErr: true,
		},
		{
			name: "valid MapId - Clark/Lake",
			input: traintracker.ArrivalsProps{
				MapId: "40380",
			},
			wantErr: false,
		},
		{
			name: "valid StpId",
			input: traintracker.ArrivalsProps{
				StpId: "30070",
			},
			wantErr: false,
		},
		{
			name: "invalid MapId - non-numeric",
			input: traintracker.ArrivalsProps{
				MapId: "bogus",
			},
			wantErr: true,
		},
		{
			name: "invalid StpId - non-numeric",
			input: traintracker.ArrivalsProps{
				StpId: "bogus",
			},
			wantErr: true,
		},
		{
			name: "valid MapId with route filter",
			input: traintracker.ArrivalsProps{
				MapId: "40380",
				Rt:    "Blue",
			},
			wantErr: false,
		},
		{
			name: "valid MapId with invalid route",
			input: traintracker.ArrivalsProps{
				MapId: "40380",
				Rt:    "InvalidRoute",
			},
			wantErr: true,
		},
		{
			name: "valid MapId with max results",
			input: traintracker.ArrivalsProps{
				MapId: "40380",
				Max:   5,
			},
			wantErr: false,
		},
		{
			name: "all valid routes - Red",
			input: traintracker.ArrivalsProps{
				MapId: "40380",
				Rt:    "Red",
			},
			wantErr: false,
		},
		{
			name: "all valid routes - Brn",
			input: traintracker.ArrivalsProps{
				MapId: "40380",
				Rt:    "Brn",
			},
			wantErr: false,
		},
		{
			name: "all valid routes - G",
			input: traintracker.ArrivalsProps{
				MapId: "40380",
				Rt:    "G",
			},
			wantErr: false,
		},
		{
			name: "all valid routes - Org",
			input: traintracker.ArrivalsProps{
				MapId: "40380",
				Rt:    "Org",
			},
			wantErr: false,
		},
		{
			name: "all valid routes - P",
			input: traintracker.ArrivalsProps{
				MapId: "40380",
				Rt:    "P",
			},
			wantErr: false,
		},
		{
			name: "all valid routes - Pink",
			input: traintracker.ArrivalsProps{
				MapId: "40380",
				Rt:    "Pink",
			},
			wantErr: false,
		},
		{
			name: "all valid routes - Y",
			input: traintracker.ArrivalsProps{
				MapId: "40380",
				Rt:    "Y",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := tracker.Arrivals(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil (input: %+v)", tt.input)
				}
				return
			}

			// For valid inputs, we expect either success or an HTTP/API error
			// (not a validation error). Since we're using a fake key,
			// we'll get an API error, but that's expected.
			// The test passes if we don't get a validation error.
			if err != nil {
				// Check if it's a validation error vs API error
				// Validation errors contain specific messages
				errStr := err.Error()
				validationErrors := []string{
					"must specify either MapId or StpId",
					"error converting mapId",
					"error converting StpId",
					"invalid mapId",
					"invalid StpId",
					"invalid Rt",
				}

				for _, ve := range validationErrors {
					if contains(errStr, ve) {
						t.Errorf("unexpected validation error: %v (input: %+v)", err, tt.input)
						return
					}
				}
				// Non-validation errors (HTTP, API key, etc.) are acceptable
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
