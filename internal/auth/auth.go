package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/11Petrov/gopherloyal/internal/logger"
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

const (
	TokenEXP  = time.Hour * 3
	SecretKEY = "supersecretkey"
)

func WriteToken(ctx context.Context, rw http.ResponseWriter, userID int) error {
	log := logger.FromContext(ctx)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenEXP)),
		},
		UserID: userID,
	})

	tokenString, err := token.SignedString([]byte(SecretKEY))
	if err != nil {
		log.Errorf("error tokenString in BuildJWTString()... ", err)
		return err
	}
	http.SetCookie(rw, &http.Cookie{
		Name:    "Token",
		Value:   tokenString,
		Expires: time.Now().Add(TokenEXP),
	})

	return nil
}

func GetUserID(ctx context.Context, r *http.Request) (int, error) {
	log := logger.FromContext(ctx)
	claims := &Claims{}
	cookie, err := r.Cookie("Token")
	if err != nil {
		log.Errorf("error in GetUserID, ", err)
		return 0, err
	}

	token, err := jwt.ParseWithClaims(cookie.Value, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(SecretKEY), nil
	})

	if !token.Valid {
		log.Errorf("no valid token ...", err)
		return 0, err
	}

	return claims.UserID, err
}
