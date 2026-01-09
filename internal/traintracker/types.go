package traintracker

import "time"

const (
	CtaBaseUrl = "https://lapi.transitchicago.com"
)

type Station struct {
	Name string
	Id   string
}

type LatLon struct {
	Latitude  float64
	Longitude float64
}

type Train struct {
	RunNumber   uint16
	Line        string
	Destination Station
	Direction   string
	Position    LatLon
}

type TrainArrival struct {
	Station     Station
	Train       Train
	Time        time.Time
	IsDue       bool
	IsScheduled bool
}

type ArrivalApiResponse struct {
	Ctatt struct {
		Tmst  string   `json:"tmst"`
		ErrCd string   `json:"errCd"`
		ErrNm string   `json:"errNm"`
		Eta   []ApiEta `json:"eta"`
	} `json:"ctatt"`
}

type ApiEta struct {
	StaId   string   `json:"staId"`
	StpId   string   `json:"stpId"`
	StaNm   string   `json:"staNm"`
	StpDe   string   `json:"stpDe"`
	Rn      string   `json:"rn"`
	Rt      string   `json:"rt"`
	DestSt  string   `json:"destSt"`
	DestNm  string   `json:"destNm"`
	TrDr    string   `json:"trDr"`
	Prdt    string   `json:"prdt"`
	ArrT    string   `json:"arrT"`
	IsApp   string   `json:"isApp"`
	IsSch   string   `json:"isSch"`
	IsDly   string   `json:"isDly"`
	IsFlt   string   `json:"isFlt"`
	Flags   []string `json:"flags"`
	Lat     string   `json:"lat"`
	Lon     string   `json:"lon"`
	Heading string   `json:"heading"`
}

type ArrivalResponse struct {
	Ctatt Ctatt
}

type Ctatt struct {
	Tmst  time.Time
	ErrCd string
	ErrNm string
	Eta   []Eta
}

type Eta struct {
	StaId   string
	StpId   string
	StaNm   string
	StpDe   string
	Rn      string
	Rt      string
	DestSt  string
	DestNm  string
	TrDr    int16
	Prdt    time.Time
	ArrT    time.Time
	IsApp   bool
	IsSch   bool
	IsDly   bool
	IsFlt   bool
	Flags   []string
	Lat     float64
	Lon     float64
	Heading int16
}
