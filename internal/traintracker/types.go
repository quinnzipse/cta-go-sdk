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
	Lat     *float64
	Lon     *float64
	Heading *int16
}

type PositionsApiResponse struct {
	Ctatt struct {
		Tmst  string          `json:"tmst"`
		ErrCd string          `json:"errCd"`
		ErrNm string          `json:"errNm"`
		Route []ApiRoutePos  `json:"route"`
	} `json:"ctatt"`
}

type ApiRoutePos struct {
	Rt       string         `json:"rt"`
	Position []ApiTrainPos  `json:"position"`
}

type ApiTrainPos struct {
	TrainId   string   `json:"trainId"`
	RunNumber string   `json:"rn"`
	Rt        string   `json:"rt"`
	TrDr      string   `json:"trDr"`
	Lat       string   `json:"lat"`
	Lon       string   `json:"lon"`
	Heading   string   `json:"heading"`
	NextStpId string   `json:"nextStpId"`
	NextStpNm string   `json:"nextStpNm"`
	NextStpTm string   `json:"nextStpTm"`
	PrevStpId string   `json:"prevStpId"`
	PrevStpNm string   `json:"prevStpNm"`
	PrevStpTm string   `json:"prevStpTm"`
	IsSch     string   `json:"isSch"`
	IsDly     string   `json:"isDly"`
	IsFlt     string   `json:"isFlt"`
	IsBoe     string   `json:"isBoe"`
	IsDet     string   `json:"isDet"`
	Flags     []string `json:"flags"`
}

type PositionsResponse struct {
	Ctatt CtattPos
}

type CtattPos struct {
	Tmst  time.Time
	ErrCd string
	ErrNm string
	Route []RoutePosition
}

type RoutePosition struct {
	Rt       string
	Position []TrainPosition
}

type TrainPosition struct {
	TrainId   string
	RunNumber string
	Rt        string
	TrDr      int16
	Lat       float64
	Lon       float64
	Heading   int16
	NextStpId string
	NextStpNm string
	NextStpTm *time.Time
	PrevStpId string
	PrevStpNm string
	PrevStpTm *time.Time
	IsSch     bool
	IsDly     bool
	IsFlt     bool
	IsBoe     bool
	IsDet     bool
	Flags     []string
}

type FollowProps struct {
	RunNumber string
}

const followPath = "/api/1.0/ttfollow.aspx"

type FollowApiResponse struct {
	Ctatt struct {
		Tmst   string        `json:"tmst"`
		ErrCd  string        `json:"errCd"`
		ErrNm  string        `json:"errNm"`
		Run   string        `json:"run"`
		Rt    string        `json:"rt"`
		DestNm string        `json:"destNm"`
		Route ApiFollowRoute `json:"route"`
	} `json:"ctatt"`
}

type ApiFollowRoute struct {
	Position []ApiFollowPosition `json:"position"`
}

type ApiFollowPosition struct {
	StpId   string `json:"stpId"`
	StpNm   string `json:"stpNm"`
	StpDe   string `json:"stpDe"`
	Prdt    string `json:"prdt"`
	ArrT    string `json:"arrT"`
	IsApp   string `json:"isApp"`
	IsSch   string `json:"isSch"`
	IsDly   string `json:"isDly"`
	IsFlt   string `json:"isFlt"`
}

type FollowResponse struct {
	Ctatt FollowCtatt
}

type FollowCtatt struct {
	Tmst    time.Time
	ErrCd   string
	ErrNm   string
	Run    string
	Rt     string
	DestNm string
	Route  []FollowPosition
}

type FollowPosition struct {
	StpId   string
	StpNm   string
	StpDe   string
	Prdt   *time.Time
	ArrT   *time.Time
	IsApp  bool
	IsSch  bool
	IsDly  bool
	IsFlt  bool
}
