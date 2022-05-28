package services

import (
	"go-refresh/entities"
	"time"

	"github.com/golang-jwt/jwt"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

type TokenClaims struct {
	UserID uint
	jwt.StandardClaims
}

type TokenServiceInterface interface {
	GenerateToken(payload TokenClaims) (string, error)
	VerifyToken(tokenString string) (TokenClaims, error)
	GenerateRefreshToken(userId uint) (entities.RefreshToken, error)
	IsRefreshTokenExpired(expiresAt int64) bool
}

type tokenService struct{}

func NewTokenService() TokenServiceInterface {
	return &tokenService{}
}

var signingKey = []byte("AllYourBase")

func (*tokenService) GenerateToken(payload TokenClaims) (string, error) {
	fiveMinutes := time.Now().Add(5 * time.Minute).Unix()

	claims := TokenClaims{
		UserID: payload.UserID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: fiveMinutes,
			Issuer:    "github.com/golang-jwt/jwt",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(signingKey)
}

func (*tokenService) VerifyToken(tokenString string) (TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})

	if err != nil {
		return TokenClaims{}, err
	}

	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return *claims, nil
	}

	return TokenClaims{}, nil
}

func (*tokenService) GenerateRefreshToken(userId uint) (entities.RefreshToken, error) {
	uniqueId, err := gonanoid.New()
	if err != nil {
		return entities.RefreshToken{}, err
	}

	sevenDays := time.Now().Add(time.Hour * 24 * 7).Unix()

	token := entities.RefreshToken{
		UserID:    userId,
		Token:     uniqueId,
		ExpiresAt: sevenDays,
	}

	return token, nil
}

func (*tokenService) IsRefreshTokenExpired(expiresAt int64) bool {
	return expiresAt < time.Now().Unix()
}

// This is the TokenService instance
var TokenService TokenServiceInterface = NewTokenService()
