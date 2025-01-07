package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type JwtGenerateRequest struct {
	Keys      map[string]interface{} `json:"keys"`
	JwtKey    string
	ExpiresAt int64 `json:"expires_at"`
}

func GenerateJWT(keys map[string]interface{}, jwtKey string) (string, error) {
	claims := jwt.MapClaims{}

	for key, value := range keys {
		claims[key] = value
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseJWT(tokenString string, jwtKey string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the algorithm
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Return the secret key
		return []byte(jwtKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
