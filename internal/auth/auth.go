package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/11Petrov/gopherloyal/internal/logger"
	"github.com/golang-jwt/jwt/v4"
)

type key int

const (
	userIDKey key = iota
	TokenEXP      = time.Hour * 3
	SecretKEY     = "supersecretkey"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

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
		log.Errorf("error generating token string: %v", err)
		return err
	}
	http.SetCookie(rw, &http.Cookie{
		Name:    "Token",
		Value:   tokenString,
		Expires: time.Now().Add(TokenEXP),
	})

	return nil
}

func UserIDFromContext(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value(userIDKey).(int)
	return userID, ok
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := logger.FromContext(ctx)

		claims := &Claims{}
		cookie, err := r.Cookie("Token")
		if err != nil {
			log.Errorf("error getting cookie: %v", err)
			http.Error(rw, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token, err := jwt.ParseWithClaims(cookie.Value, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(SecretKEY), nil
		})

		if err != nil || !token.Valid {
			log.Errorf("invalid or expired token: %v", err)
			http.Error(rw, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ctx = context.WithValue(ctx, userIDKey, claims.UserID)
		next.ServeHTTP(rw, r.WithContext(ctx))
	})
}
