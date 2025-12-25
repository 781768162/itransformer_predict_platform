package pkg

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ServerSecret    = "testSecretKey012345"
	ExpireDuration  = time.Hour * 24
	SecretAlgorithm = jwt.SigningMethodHS256
)

type Claims struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"user_name"`
	jwt.RegisteredClaims
}

func NewToken(userID int64, username string) (string, error) {
	nowTime := time.Now()
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(nowTime.Add(ExpireDuration)),
			IssuedAt:  jwt.NewNumericDate(nowTime),
			Issuer:    "test",
		},
	}

	token := jwt.NewWithClaims(SecretAlgorithm, claims)

	tokenStr, err := token.SignedString(ServerSecret)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}

func ParseToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected algorithm")
		}
		return ServerSecret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("token invalid")
}
