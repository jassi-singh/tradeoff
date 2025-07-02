package domain

type Player struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

type PriceData struct {
	Time  string  `json:"time"`
	Open  float64 `json:"open"`
	High  float64 `json:"high"`
	Low   float64 `json:"low"`
	Close float64 `json:"close"`
}

type WsMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}
