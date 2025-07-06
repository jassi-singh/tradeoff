package service

import (
	"context"
	"time"
	"tradeoff/backend/internal/domain"

	polygon "github.com/polygon-io/client-go/rest"
	"github.com/polygon-io/client-go/rest/models"
)

type MarketService struct {
	hub           *Hub
	polygonClient *polygon.Client
}

func NewMarketService(hub *Hub, apiKey string) *MarketService {
	polygonClient := polygon.New(apiKey)
	return &MarketService{
		hub:           hub,
		polygonClient: polygonClient,
	}
}

func polygonPriceDataToDomain(p models.Agg) domain.PriceData {
	return domain.PriceData{
		Time:   time.Time(p.Timestamp),
		Open:   p.Open,
		High:   p.High,
		Low:    p.Low,
		Close:  p.Close,
		Volume: p.Volume,
	}
}

func (m *MarketService) LoadPriceData(ctx context.Context, ticker string, from time.Time, to time.Time, timespan models.Timespan) ([]domain.PriceData, error) {
	priceData := []domain.PriceData{}
	params := models.ListAggsParams{
		Ticker:     ticker,
		Multiplier: 1,
		Timespan:   timespan,
		From:       models.Millis(from),
		To:         models.Millis(to),
	}.
		WithAdjusted(true).
		WithOrder(models.Order("asc"))

	iter := m.polygonClient.ListAggs(ctx, params)

	for iter.Next() {
		agg := iter.Item()
		priceData = append(priceData, polygonPriceDataToDomain(agg))
	}
	if iter.Err() != nil {
		return nil, iter.Err()
	}

	return priceData, nil
}
