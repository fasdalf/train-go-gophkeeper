package jwt

import (
	"fmt"
	jwt "github.com/golang-jwt/jwt/v4"
	"time"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID uint64
}

const (
	AuthHeader = "Authorization"
	AuthPrefix = "Bearer "
)

// BuildJWTString создаёт токен и возвращает его в виде строки
func BuildJWTString(id uint64, key *string, exp time.Duration) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// когда истекает токен
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(exp)),
		},
		// собственное утверждение
		UserID: id,
	})

	tokenString, _ := token.SignedString([]byte(*key))

	return tokenString
}

// GetUserID валидирует токен и возвращает ID пользователя
func GetUserID(tokenString, key *string) (uint64, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(*tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Method.Alg())
		}
		return []byte(*key), nil
	})
	if err != nil {
		return 0, err
	}

	return claims.UserID, nil
}
