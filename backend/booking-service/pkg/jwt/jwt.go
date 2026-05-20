package jwt

import (
	"booking-service/pkg/config"
	"time"

	jwtLib "github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID uint) (string, error) {
	claims := jwtLib.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwtLib.NewWithClaims(jwtLib.SigningMethodHS256, claims)
	cfg := config.Get();

	return token.SignedString([]byte(cfg.JWT.Secret))
}