package helpers

import (
	"encoding/json"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type CustomError struct {
	Message string `json:"message"`
	status  int    `json:"-"`
}

func (e *CustomError) Error() string {
	return e.Message
}

func NewCustomError(message string, status int) *CustomError {
	return &CustomError{
		Message: message,
		status:  status,
	}
}

func RespondWithError(w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case *CustomError:
		RespondWithJSON(w, e.status, map[string]string{"error": e.Message})
	default:
		RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal Server Error"})
	}
}

func RespondWithJSON(w http.ResponseWriter, status int, data interface{}) {
	response, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}

// ValidateJWTAndGetPlayerID validates a JWT token and returns the player ID
func ValidateJWTAndGetPlayerID(tokenStr, jwtSecret string) (string, error) {
	// Parse and validate the token
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", jwt.ErrTokenMalformed
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", jwt.ErrInvalidKey
	}

	// Get player ID from claims
	playerID, ok := claims["sub"].(string)
	if !ok || playerID == "" {
		return "", jwt.ErrInvalidKey
	}

	return playerID, nil
}
