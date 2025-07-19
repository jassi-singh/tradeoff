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
	Type             PositionType `json:"type"`
	EntryPrice       float64      `json:"entryPrice"`
	EntryTime        time.Time    `json:"entryTime"`
	ExitPrice        float64      `json:"exitPrice"`
	ExitTime         time.Time    `json:"exitTime"`
	Profit           float64      `json:"profit"`
	ProfitPercentage float64      `json:"profitPercentage"`
}

type PlayerSession struct {
	PlayerId        string     `json:"playerId"`
	RoundID         string     `json:"roundId"`
	Balance         float64    `json:"balance"`
	ActivePosition  Position   `json:"activePosition"`
	ClosedPositions []Position `json:"closedPositions"`
}
