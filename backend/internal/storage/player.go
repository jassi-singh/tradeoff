package storage

import (
	"fmt"
	"log"
	"tradeoff/backend/internal/domain"
)

func (s *PostgresStore) CreatePlayer(player domain.Player) (domain.Player, error) {
	// Let the database auto-generate the UUID
	query := `INSERT INTO public.players (username, refresh_token, refresh_token_expiry) VALUES ($1, $2, $3) RETURNING id`
	var id string
	err := s.DB.QueryRow(query, player.Username, player.RefreshToken, player.RefreshTokenExpiry).Scan(&id)
	if err != nil {
		log.Println("Error creating player:", err)
		return domain.Player{}, err
	}

	player.Id = id
	return player, nil
}

func (s *PostgresStore) FindPlayerByRefreshToken(refreshToken string) (domain.Player, error) {
	query := `SELECT id, username, refresh_token, refresh_token_expiry FROM public.players WHERE refresh_token = $1`
	var player domain.Player
	err := s.DB.QueryRow(query, refreshToken).Scan(&player.Id, &player.Username, &player.RefreshToken, &player.RefreshTokenExpiry)
	if err != nil {
		return domain.Player{}, fmt.Errorf("failed to find player by refresh token: %w", err)
	}
	return player, nil
}

func (s *PostgresStore) UpdatePlayer(player domain.Player) (domain.Player, error) {
	// Initialize the base query and parameters
	query := "UPDATE public.players SET"
	params := []any{}
	paramCounter := 1

	// Dynamically build the query based on provided fields
	if player.Username != "" {
		query += " username = $" + fmt.Sprint(paramCounter) + ","
		params = append(params, player.Username)
		paramCounter++
	}
	if player.RefreshToken != "" {
		query += " refresh_token = $" + fmt.Sprint(paramCounter) + ","
		params = append(params, player.RefreshToken)
		paramCounter++
	}
	if !player.RefreshTokenExpiry.IsZero() {
		query += " refresh_token_expiry = $" + fmt.Sprint(paramCounter) + ","
		params = append(params, player.RefreshTokenExpiry)
		paramCounter++
	}

	// Remove the trailing comma and add the WHERE clause
	query = query[:len(query)-1] + " WHERE id = $" + fmt.Sprint(paramCounter) + " RETURNING id, username, refresh_token, refresh_token_expiry"
	params = append(params, player.Id)

	// Execute the query
	updatedPlayer := domain.Player{}
	err := s.DB.QueryRow(query, params...).Scan(&updatedPlayer.Id, &updatedPlayer.Username, &updatedPlayer.RefreshToken, &updatedPlayer.RefreshTokenExpiry)
	if err != nil {
		return domain.Player{}, fmt.Errorf("failed to update player: %w", err)
	}

	return updatedPlayer, nil
}

func (s *PostgresStore) GetPlayer(id string) (domain.Player, error) {
	query := `SELECT id, username, refresh_token, refresh_token_expiry FROM public.players WHERE id = $1`
	var player domain.Player
	err := s.DB.QueryRow(query, id).Scan(&player.Id, &player.Username, &player.RefreshToken, &player.RefreshTokenExpiry)
	if err != nil {
		return domain.Player{}, err
	}
	return player, nil
}
