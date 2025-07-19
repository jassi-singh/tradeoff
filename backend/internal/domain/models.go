package domain

import "time"

type Player struct {
	Id                 string    `json:"id"`
	Username           string    `json:"username"`
	RefreshToken       string    `json:"-"`
	RefreshTokenExpiry time.Time `json:"-"`
}

type PriceData struct {
	Time   int64   `json:"time"`
	Open   float64 `json:"open"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
	Volume float64 `json:"volume"`
}

type Phase string

const (
	Lobby  Phase = "lobby"
	Live   Phase = "live"
	Closed Phase = "closed"
)

type PositionType string

const (
	PositionTypeLong  PositionType = "long"
	PositionTypeShort PositionType = "short"
)

type Position struct {
	Type 			 PositionType
	EntryPrice       float64
	EntryTime        time.Time
	ExitPrice        float64
	ExitTime         time.Time
	Profit           float64
	ProfitPercentage float64
}

type PlayerSession struct {
	PlayerId        string
	RoundID         string
	Balance         float64
	ActivePosition  Position
	ClosedPositions []Position
}
