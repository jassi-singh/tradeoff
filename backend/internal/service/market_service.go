package service

import (
	"context"
	"log"
	"math/rand"
	"time"
	"tradeoff/backend/internal/domain"

	polygon "github.com/polygon-io/client-go/rest"
	"github.com/polygon-io/client-go/rest/models"
)

type MarketService struct {
	priceData     []models.Agg
	hub           *Hub
	polygonClient *polygon.Client
}

func NewMarketService(hub *Hub, apiKey string) *MarketService {
	polygonClient := polygon.New(apiKey)
	return &MarketService{
		hub:           hub,
		priceData:     []models.Agg{},
		polygonClient: polygonClient,
	}
}

func polygonPriceDataToDomain(p models.Agg) domain.PriceData {
	return domain.PriceData{
		Time:   time.Time(p.Timestamp).Unix(),
		Open:   p.Open,
		High:   p.High,
		Low:    p.Low,
		Close:  p.Close,
		Volume: p.Volume,
	}
}

func (m *MarketService) LoadPriceData() {
	from, err := time.Parse("2006-01-02", "2024-01-09")
	if err != nil {
		log.Fatalf("Error parsing 'from' date: %v", err)
	}
	to, err := time.Parse("2006-01-02", "2025-01-10")
	if err != nil {
		log.Fatalf("Error parsing 'to' date: %v", err)
	}

	params := models.ListAggsParams{
		Ticker:     "X:BTCUSD",
		Multiplier: 1,
		Timespan:   "day",
		From:       models.Millis(from),
		To:         models.Millis(to),
	}.
		WithAdjusted(true).
		WithOrder(models.Order("asc"))

	aggs := m.polygonClient.ListAggs(context.Background(), params)

	for aggs.Next() {
		agg := aggs.Item()
		m.priceData = append(m.priceData, agg)
	}
}

func (m *MarketService) StartPriceFeed() {
	for {
		priceData := []domain.PriceData{}
		for _, p := range m.priceData {
			priceData = append(priceData, polygonPriceDataToDomain(p))
		}
		m.hub.Broadcast <- domain.WsMessage{
			Type:    "price_data",
			Payload: priceData,
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
