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
	Type       PositionType `json:"type"`
	EntryPrice float64      `json:"entryPrice"`
	EntryTime  time.Time    `json:"entryTime"`
}

type ClosedPosition struct {
	Position
	ExitPrice        float64 `json:"exitPrice"`
	ExitTime         time.Time `json:"exitTime"`
	ProfitLoss       float64 `json:"profitLoss"`
	ProfitLossPercentage float64 `json:"profitLossPercentage"`
}

type ActivePosition struct {
	Position
	PnL           float64 `json:"pnl"`
	PnlPercentage float64 `json:"pnlPercentage"`
}

type PlayerSession struct {
	PlayerId        string          `json:"playerId"`
	RoundID         string          `json:"roundId"`
	Balance         float64         `json:"balance"`
	ActivePosition  *ActivePosition `json:"activePosition"`
	ClosedPositions []ClosedPosition      `json:"closedPositions"`
}
