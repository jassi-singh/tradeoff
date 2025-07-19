package service

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"time"
	"tradeoff/backend/internal/domain"

	"github.com/golang-jwt/jwt/v5"
)

type UserWithToken struct {
	User         domain.Player `json:"user"`
	Token        string        `json:"token"`
	RefreshToken string        `json:"refreshToken"`
}

type AuthService struct {
	playerRepository PlayerRepository
	jwtSecret        string
	jwtExpiration    int64
}

func NewAuthService(playerRepository PlayerRepository, jwtSecret string, jwtExpiration int64) *AuthService {
	return &AuthService{
		playerRepository: playerRepository,
		jwtSecret:        jwtSecret,
		jwtExpiration:    jwtExpiration,
	}
}

func (s *AuthService) Login(payload LoginParams) (UserWithToken, error) {
	player := domain.Player{
		Username: payload.Username,
	}

	refreshToken, err := s.generateRefreshToken()
	if err != nil {
		return UserWithToken{}, err
	}
	player.RefreshToken = refreshToken
	player.RefreshTokenExpiry = time.Now().Add(24 * time.Hour)

	// Create player in database first to get the UUID
	player, err = s.playerRepository.CreatePlayer(player)
	if err != nil {
		return UserWithToken{}, fmt.Errorf("failed to create player: %w", err)
	}

	// Now generate token with the player that has an ID
	token, err := s.generateToken(player)
	if err != nil {
		return UserWithToken{}, err
	}

	return UserWithToken{
		User:         player,
		Token:        token,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) RefreshToken(refreshToken string) (UserWithToken, error) {
	// Retrieve the player based on the refresh token
	player, err := s.playerRepository.FindPlayerByRefreshToken(refreshToken)
	if err != nil {
		return UserWithToken{}, fmt.Errorf("invalid refresh token")
	}

	// Check if the refresh token is still valid
	if player.RefreshTokenExpiry.Unix() < time.Now().Unix() {
		return UserWithToken{}, fmt.Errorf("refresh token has expired")
	}

	// Generate a new access token and refresh token
	token, err := s.generateToken(player)
	if err != nil {
		return UserWithToken{}, err
	}

	newRefreshToken, err := s.generateRefreshToken()
	if err != nil {
		return UserWithToken{}, err
	}

	// Update the player's refresh token and expiry
	player.RefreshToken = newRefreshToken
	player.RefreshTokenExpiry = time.Now().Add(24 * time.Hour)
	player, err = s.playerRepository.UpdatePlayer(player)
	if err != nil {
		return UserWithToken{}, err
	}

	return UserWithToken{
		User:         player,
		Token:        token,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *AuthService) generateToken(user domain.Player) (string, error) {
	now := time.Now().Unix()
	claims := jwt.MapClaims{
		"sub":  user.Id,
		"name": user.Username,
		"exp":  now + s.jwtExpiration, // Set expiration as current time + duration
		"iat":  now,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		log.Println("Failed to sign JWT token", err)
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) generateRefreshToken() (string, error) {
	b := make([]byte, 32)

	_, err := rand.Read(b)
	if err != nil {
		log.Println("Failed to generate random bytes for refresh token", err)
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}
