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
	dailyPriceData  []models.Agg
	hourlyPriceData []models.Agg
	hub             *Hub
	polygonClient   *polygon.Client
}

func NewMarketService(hub *Hub, apiKey string) *MarketService {
	polygonClient := polygon.New(apiKey)
	return &MarketService{
		hub:             hub,
		dailyPriceData:  []models.Agg{},
		hourlyPriceData: []models.Agg{},
		polygonClient:   polygonClient,
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
	from := time.Now().AddDate(-2, 0, 0)
	randomDecrease := -5 - int(rand.Float64()*14)
	to := time.Now().AddDate(0, randomDecrease, 0)

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
		m.dailyPriceData = append(m.dailyPriceData, agg)
	}

	// get the next 10 days days as hourly
	from = time.Now().AddDate(0, randomDecrease, 1)
	to = time.Now().AddDate(0, randomDecrease, 10)

	params = models.ListAggsParams{
		Ticker:     "X:BTCUSD",
		Multiplier: 1,
		Timespan:   "hour",
		From:       models.Millis(from),
		To:         models.Millis(to),
	}.
		WithAdjusted(true).
		WithOrder(models.Order("asc"))

	aggs = m.polygonClient.ListAggs(context.Background(), params)

	for aggs.Next() {
		agg := aggs.Item()
		m.hourlyPriceData = append(m.hourlyPriceData, agg)
	}

}

func (m *MarketService) StartPriceFeed() {
	priceData := []domain.PriceData{}
	for _, p := range m.dailyPriceData {
		priceData = append(priceData, polygonPriceDataToDomain(p))
	}

	m.hub.Broadcast <- domain.WsMessage{
		Type:    "price_data",
		Payload: priceData,
	}

	var lastPrice *domain.PriceData
	log.Println(len(m.hourlyPriceData))

	for _, p := range m.hourlyPriceData {
		currentHourlyPrice := polygonPriceDataToDomain(p)

		if lastPrice == nil {
			// First price data point
			lastPrice = &currentHourlyPrice
			priceData = append(priceData, *lastPrice)
		} else {
			// Check if we've moved to a new day (86400 seconds = 24 hours)
			isNewDay := lastPrice.Time/86400 != currentHourlyPrice.Time/86400
			
			if isNewDay {
				// Start new day with current hourly price
				lastPrice = &currentHourlyPrice
				priceData = append(priceData, *lastPrice)
			} else {
				// Update existing day's OHLCV data
				lastPrice.High = max(lastPrice.High, currentHourlyPrice.Close)
				lastPrice.Low = min(lastPrice.Low, currentHourlyPrice.Close)
				lastPrice.Close = currentHourlyPrice.Close
				lastPrice.Volume += currentHourlyPrice.Volume
			}
		}

		time.Sleep(500 * time.Millisecond)

		priceData[len(priceData)-1] = *lastPrice

		m.hub.Broadcast <- domain.WsMessage{
			Type:    "price_data",
			Payload: priceData,
		}
	}
}
