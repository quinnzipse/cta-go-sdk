package traintracker

import (
	"testing"
)

type expectedEta struct {
	lat     *float64
	lon     *float64
	heading *int16
	isSch   bool
}

type toArrivalResponseTest struct {
	name     string
	input    []ApiEta
	expected []expectedEta
}

func TestToArrivalResponse(t *testing.T) {
	t.Parallel()

	tests := []toArrivalResponseTest{
		{
			name: "scheduled train with empty location",
			input: []ApiEta{
				{
					StaId:   "40380",
					StpId:   "30075",
					StaNm:   "Clark/Lake",
					StpDe:   "Service at Outer Loop platform",
					Rn:      "430",
					Rt:      "Brn",
					DestSt:  "30249",
					DestNm:  "Kimball",
					TrDr:    "1",
					Prdt:    "2026-01-09T16:55:29",
					ArrT:    "2026-01-09T16:56:29",
					IsApp:   "0",
					IsSch:   "1",
					IsDly:   "0",
					IsFlt:   "0",
					Flags:   nil,
					Lat:     "",
					Lon:     "",
					Heading: "",
				},
			},
			expected: []expectedEta{
				{lat: nil, lon: nil, heading: nil, isSch: true},
			},
		},
		{
			name: "active train with location",
			input: []ApiEta{
				{
					StaId:   "40380",
					StpId:   "30075",
					StaNm:   "Clark/Lake",
					StpDe:   "Service at Outer Loop platform",
					Rn:      "430",
					Rt:      "Brn",
					DestSt:  "30249",
					DestNm:  "Kimball",
					TrDr:    "1",
					Prdt:    "2026-01-09T16:55:29",
					ArrT:    "2026-01-09T16:56:29",
					IsApp:   "1",
					IsSch:   "0",
					IsDly:   "0",
					IsFlt:   "0",
					Flags:   nil,
					Lat:     "41.88443",
					Lon:     "-87.62622",
					Heading: "358",
				},
			},
			expected: []expectedEta{
				{lat: ptr(41.88443), lon: ptr(-87.62622), heading: ptrInt16(358), isSch: false},
			},
		},
		{
			name: "mixed scheduled and active trains",
			input: []ApiEta{
				{
					StaId:   "40380",
					StpId:   "30075",
					StaNm:   "Clark/Lake",
					StpDe:   "Service at Outer Loop platform",
					Rn:      "430",
					Rt:      "Brn",
					DestSt:  "30249",
					DestNm:  "Kimball",
					TrDr:    "1",
					Prdt:    "2026-01-09T16:55:29",
					ArrT:    "2026-01-09T16:56:29",
					IsApp:   "1",
					IsSch:   "0",
					IsDly:   "0",
					IsFlt:   "0",
					Flags:   nil,
					Lat:     "41.88443",
					Lon:     "-87.62622",
					Heading: "358",
				},
				{
					StaId:   "40380",
					StpId:   "30075",
					StaNm:   "Clark/Lake",
					StpDe:   "Service at Outer Loop platform",
					Rn:      "431",
					Rt:      "Brn",
					DestSt:  "30249",
					DestNm:  "Kimball",
					TrDr:    "1",
					Prdt:    "2026-01-09T17:00:00",
					ArrT:    "2026-01-09T17:05:00",
					IsApp:   "0",
					IsSch:   "1",
					IsDly:   "0",
					IsFlt:   "0",
					Flags:   nil,
					Lat:     "",
					Lon:     "",
					Heading: "",
				},
			},
			expected: []expectedEta{
				{lat: ptr(41.88443), lon: ptr(-87.62622), heading: ptrInt16(358), isSch: false},
				{lat: nil, lon: nil, heading: nil, isSch: true},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			apiResponse := ArrivalApiResponse{
				Ctatt: struct {
					Tmst  string   `json:"tmst"`
					ErrCd string   `json:"errCd"`
					ErrNm string   `json:"errNm"`
					Eta   []ApiEta `json:"eta"`
				}{
					Tmst:  "2026-01-09T16:55:45",
					ErrCd: "0",
					ErrNm: "",
					Eta:   tt.input,
				},
			}

			response, err := apiResponse.toArrivalResponse()
			if err != nil {
				t.Fatalf("toArrivalResponse failed: %v", err)
			}

			if len(response.Ctatt.Eta) != len(tt.expected) {
				t.Fatalf("expected %d etas, got %d", len(tt.expected), len(response.Ctatt.Eta))
			}

			for i, exp := range tt.expected {
				eta := response.Ctatt.Eta[i]

				if !ptrEqualFloat64(exp.lat, eta.Lat) {
					t.Errorf("eta[%d]: expected Lat to be %v, got %v", i, exp.lat, eta.Lat)
				}

				if !ptrEqualFloat64(exp.lon, eta.Lon) {
					t.Errorf("eta[%d]: expected Lon to be %v, got %v", i, exp.lon, eta.Lon)
				}

				if !ptrEqualInt16(exp.heading, eta.Heading) {
					t.Errorf("eta[%d]: expected Heading to be %v, got %v", i, exp.heading, eta.Heading)
				}

				if eta.IsSch != exp.isSch {
					t.Errorf("eta[%d]: expected IsSch to be %v, got %v", i, exp.isSch, eta.IsSch)
				}
			}
		})
	}
}

func ptr(v float64) *float64 {
	return &v
}

func ptrInt16(v int16) *int16 {
	return &v
}

func ptrEqualFloat64(a, b *float64) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}

func ptrEqualInt16(a, b *int16) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}
