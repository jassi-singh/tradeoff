package storage

import (
	"tradeoff/backend/internal/domain"

	"gorm.io/gorm"
)

// Helper function to convert PlayerModel to domain.Player
func (pm *PlayerModel) ToDomain() domain.Player {
	return domain.Player{
		Id:                 pm.ID,
		Username:           pm.Username,
		RefreshToken:       pm.RefreshToken,
		RefreshTokenExpiry: pm.RefreshTokenExpiry,
	}
}

// Helper function to convert domain.Player to PlayerModel
func FromDomain(player domain.Player) PlayerModel {
	return PlayerModel{
		ID:                 player.Id,
		Username:           player.Username,
		RefreshToken:       player.RefreshToken,
		RefreshTokenExpiry: player.RefreshTokenExpiry,
	}
}

func (s *PostgresStore) CreatePlayer(player domain.Player) (domain.Player, error) {
	// Convert domain player to GORM model
	playerModel := FromDomain(player)

	// Let GORM handle the UUID generation
	playerModel.ID = ""

	// Create the player
	if err := s.DB.Create(&playerModel).Error; err != nil {
		return domain.Player{}, err
	}

	// Return the created player with generated ID
	return playerModel.ToDomain(), nil
}

func (s *PostgresStore) FindPlayerByRefreshToken(refreshToken string) (domain.Player, error) {
	var playerModel PlayerModel
	err := s.DB.Where("refresh_token = ?", refreshToken).First(&playerModel).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.Player{}, gorm.ErrRecordNotFound
		}
		return domain.Player{}, err
	}

	return playerModel.ToDomain(), nil
}

func (s *PostgresStore) UpdatePlayer(player domain.Player) (domain.Player, error) {
	var playerModel PlayerModel

	// Find the existing player
	if err := s.DB.Where("id = ?", player.Id).First(&playerModel).Error; err != nil {
		return domain.Player{}, err
	}

	// Update only non-zero fields
	updates := make(map[string]interface{})

	if player.Username != "" {
		updates["username"] = player.Username
	}
	if player.RefreshToken != "" {
		updates["refresh_token"] = player.RefreshToken
	}
	if !player.RefreshTokenExpiry.IsZero() {
		updates["refresh_token_expiry"] = player.RefreshTokenExpiry
	}

	// Perform the update
	if err := s.DB.Model(&playerModel).Updates(updates).Error; err != nil {
		return domain.Player{}, err
	}

	// Reload the updated player from database
	if err := s.DB.Where("id = ?", player.Id).First(&playerModel).Error; err != nil {
		return domain.Player{}, err
	}

	return playerModel.ToDomain(), nil
}

func (s *PostgresStore) GetPlayer(id string) (domain.Player, error) {
	var playerModel PlayerModel
	err := s.DB.Where("id = ?", id).First(&playerModel).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.Player{}, gorm.ErrRecordNotFound
		}
		return domain.Player{}, err
	}

	return playerModel.ToDomain(), nil
}
