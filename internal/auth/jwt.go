package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userId uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  &jwt.NumericDate{Time: time.Now()},
		ExpiresAt: &jwt.NumericDate{Time: time.Now().Add(expiresIn)},
		Subject:   userId.String(),
	}).SignedString([]byte(tokenSecret))

	if err != nil {
		return "", err
	}

	return token, nil
}

func ValidateJWT(tokenstr string, tokenSecret string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenstr, claims, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecret), nil
	})

	if err != nil {
		return uuid.UUID{}, err
	}

	if token.Valid {
		userId, _ := token.Claims.GetSubject()
		return uuid.MustParse(userId), nil
	} else {
		return uuid.UUID{}, errors.New("invalid token")
	}
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")

	if authHeader == "" {
		return "", errors.New("missing authorization header")
	}

	splits := strings.Split(authHeader, "Bearer ")

	if len(splits) != 2 {
		return "", errors.New("invalid authorization header format")
	}

	return splits[1], nil
}
