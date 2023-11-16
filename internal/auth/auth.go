package auth

import (
	"context"
	"time"

	"github.com/11Petrov/gopherloyal/internal/logger"
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

const (
	TokenEXP  = time.Hour * 1
	SecretKEY = "supersecretkey"
)

func GenerateToken(ctx context.Context, userID string) (string, error) {
	log := logger.LoggerFromContext(ctx)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenEXP)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(SecretKEY))
	if err != nil {
		log.Errorf("error tokenString in BuildJWTString()... ", err)
		return "", err
	}
	return tokenString, nil
}

func GetUserID(ctx context.Context, tokenString string) (string, error) {
	log := logger.LoggerFromContext(ctx)
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(SecretKEY), nil
		})
	if err != nil {
		log.Errorf("error in GetUserID, ", err)
		return "", err
	}

	if !token.Valid {
		log.Errorf("no valid token ...", err)
		return "", err
	}

	return claims.UserID, err
}
