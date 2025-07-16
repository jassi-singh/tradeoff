package middleware

import (
    "context"
    "net/http"
    "strings"

    "github.com/golang-jwt/jwt/v5"
    "tradeoff/backend/internal/config"
)

const PlayerIDKey = "player_id"

func AuthMiddleware(next http.Handler) http.Handler {
    // Load config once (could use sync.Once for production)
    conf, err := config.LoadConfig()
    var jwtSecret string
    if err == nil {
        jwtSecret = conf.JWT.Secret
    }
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
            return
        }

        parts := strings.Split(authHeader, " ")
        if len(parts) != 2 || parts[0] != "Bearer" {
            http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
            return
        }

        tokenStr := parts[1]
        token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
            return []byte(jwtSecret), nil
        })
        if err != nil || !token.Valid {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok {
            http.Error(w, "Invalid token claims", http.StatusUnauthorized)
            return
        }

        playerID, ok := claims["player_id"].(string)
        if !ok || playerID == "" {
            http.Error(w, "Missing player_id in token", http.StatusUnauthorized)
            return
        }

        ctx := context.WithValue(r.Context(), PlayerIDKey, playerID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

