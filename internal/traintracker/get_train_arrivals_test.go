package traintracker_test

import (
	"testing"

	"git.quinnzipse.dev/cta-go-sdk/internal/traintracker"
)

type want struct {
	hasErr    bool
	hasResult bool
	station   string
	direction string
}

func TestGetTrainArrivalsByStation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input traintracker.GetTrainArrivalsByStationProps
		want  want
	}{
		{
			name: "valid station - division",
			input: traintracker.GetTrainArrivalsByStationProps{
				Station: traintracker.Station{
					Name: "division",
					Id:   "test",
				},
			},
			want: want{
				hasErr:    false,
				hasResult: true,
				station:   "division",
			},
		},
		{
			name: "valid station/direction - division -> ohare",
			input: traintracker.GetTrainArrivalsByStationProps{
				Station: traintracker.Station{
					Name: "division",
					Id:   "test",
				},
				Direction: "ohare",
			},
			want: want{
				hasErr:    false,
				hasResult: true,
				station:   "division",
				direction: "ohare",
			},
		},
		{
			name: "valid station/direction - division -> forest park",
			input: traintracker.GetTrainArrivalsByStationProps{
				Station: traintracker.Station{
					Name: "division",
					Id:   "test",
				},
				Direction: "forest park",
			},
			want: want{
				hasErr:    false,
				hasResult: true,
				station:   "division",
				direction: "forest park",
			},
		},
		{
			name: "invalid station ID",
			input: traintracker.GetTrainArrivalsByStationProps{
				Station: traintracker.Station{
					Name: "division",
					Id:   "bogus",
				},
			},
			want: want{
				hasErr:    true,
				hasResult: false,
			},
		},
		{
			name: "empty station input",
			input: traintracker.GetTrainArrivalsByStationProps{
				Station: traintracker.Station{},
			},
			want: want{
				hasErr:    true,
				hasResult: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			res, err := traintracker.GetTrainArrivalsByStation(tt.input)

			if tt.want.hasErr {
				if err == nil {
					t.Fatalf("want error\ngot nil\n(input: %+v)", tt.input)
				}

				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if tt.want.hasResult && res == nil {
				t.Fatalf("want non-nil result\ngot nil\ninput %+v", tt.input)
			}

			if !tt.want.hasResult && res != nil {
				t.Fatalf("want nil result\ngot %#v\ninput %+v", res, tt.input)
			}

			if tt.want.station != "" {
				for i, arrival := range res {
					if arrival.Station.Name != tt.want.station {
						t.Fatalf("want station == %s result, got %s (i=%d)", tt.want.station, arrival.Station.Name, i)
					}
				}
			}

			if tt.want.direction != "" {
				for i, arrival := range res {
					if arrival.Train.Direction != tt.want.direction {
						t.Fatalf("want direction == %s result, got %s (i=%d)", tt.want.direction, arrival.Train.Direction, i)
					}
				}
			}
		})
	}
}
