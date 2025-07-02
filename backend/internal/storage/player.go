package storage

import "tradeoff/backend/internal/domain"


func (s *PostgresStore) CreatePlayer(player domain.Player) (domain.Player, error) {
	query := `INSERT INTO public.players (username) VALUES ($1) RETURNING id, username`
	err := s.DB.QueryRow(query, player.Username).Scan(&player.Id, &player.Username)
	if err != nil {
		return domain.Player{}, err
	}
	return player, nil
}

func (s *PostgresStore) GetPlayer(id string) (domain.Player, error) {
	query := `SELECT id, username FROM public.players WHERE id = $1`
	var player domain.Player
	err := s.DB.QueryRow(query, id).Scan(&player.Id, &player.Username)
	if err != nil {
		return domain.Player{}, err
	}
	return player, nil
}