package service

import (
	"math/rand"
	"time"
	"tradeoff/backend/internal/domain"
)

type MarketService struct {
	priceData []domain.PriceData
	hub       *Hub
}

func NewMarketService(hub *Hub) *MarketService {
	return &MarketService{
		hub:       hub,
		priceData: []domain.PriceData{},
	}
}


func (m *MarketService) LoadPriceData() {
	m.priceData = []domain.PriceData{
		{Time: "2021-01-01", Open: 100, High: 100, Low: 100, Close: 100},
	}
}

func (m *MarketService) StartPriceFeed() {
	for {
		m.hub.Broadcast <- domain.WsMessage{
			Type: "price_data",
			Payload: m.priceData,
		}
		time.Sleep(1 * time.Second)

		// randomly increase or decrease the price
		lastPrice := &m.priceData[len(m.priceData)-1]
		priceChange := rand.Float64() * 2
		
		if rand.Float64() < 0.5 {
			lastPrice.Close += priceChange
			if lastPrice.Close > lastPrice.High {
				lastPrice.High = lastPrice.Close
			}
		} else {
			lastPrice.Close -= priceChange
			if lastPrice.Close < lastPrice.Low {
				lastPrice.Low = lastPrice.Close
			}
		}
	}
}
