package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("super-secret-intern-key")

var ErrInvalidToken = errors.New("invalid token")

func GenerateToken(subject string) (string, error) {
	expiresAt := time.Now().Add(24 * time.Hour)
	claims := &jwt.RegisteredClaims{
		Subject:   subject,
		ExpiresAt: jwt.NewNumericDate(expiresAt),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ParseToken(tokenStr string) (string, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return "", ErrInvalidToken
	}
	return claims.Subject, nil
}

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		userID, err := ParseToken(tokenStr)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		r.Header.Set("X-User-ID", userID)
		next.ServeHTTP(w, r)
	}
}
