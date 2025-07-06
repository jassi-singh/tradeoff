package domain

type Player struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

type PriceData struct {
	Time   int64   `json:"time"`
	Open   float64 `json:"open"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
	Volume float64 `json:"volume"`
}

type WsMessageType string

const (
	WsMessageTypePriceData   WsMessageType = "price_data"
	WsMessageTypeRoundStatus WsMessageType = "round_status"
)

type WsMessage struct {
	Type WsMessageType `json:"type"`
	Data interface{}   `json:"payload"`
}

type RoundStatus string

const (
	RoundStatusLobby  RoundStatus = "lobby"
	RoundStatusLive   RoundStatus = "live"
	RoundStatusClosed RoundStatus = "closed"
)
