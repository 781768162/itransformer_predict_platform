package jwt

import (
	"errors"
	"time"

	"gateway/config"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ServerSecret    = config.Settings.JWT.Secret
	ExpireDuration  = config.Settings.JWT.ExpireDuration
	SecretAlgorithm = jwt.SigningMethodHS256
)

type Claims struct {
	UserID   int64  `json:"user_id"`
	UserName string `json:"user_name"`
	jwt.RegisteredClaims
}

func NewToken(userId int64, userName string) (string, int64, error) {
	nowTime := time.Now()
	expireAt := nowTime.Add(ExpireDuration)

	claims := Claims{
		UserID:   userId,
		UserName: userName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireAt),
			IssuedAt:  jwt.NewNumericDate(nowTime),
			Issuer:    "test",
		},
	}

	token := jwt.NewWithClaims(SecretAlgorithm, claims)

	tokenStr, err := token.SignedString(ServerSecret)
	if err != nil {
		return "", 0, err
	}

	return tokenStr, expireAt.Unix(), nil
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
