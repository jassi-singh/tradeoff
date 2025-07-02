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

type WsMessage struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}
