package service

import "tradeoff/backend/internal/domain"

type PlayerRepository interface {
	CreatePlayer(player domain.Player) (domain.Player, error)
	GetPlayer(id string) (domain.Player, error)
}