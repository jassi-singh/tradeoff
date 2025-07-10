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
