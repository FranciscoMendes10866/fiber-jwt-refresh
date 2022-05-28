package services

import "github.com/alexedwards/argon2id"

type HashServiceInterface interface {
	HashPassword(password string) (string, error)
	VerifyPassword(password, hash string) (bool, error)
}

type hashService struct{}

func NewHashService() HashServiceInterface {
	return &hashService{}
}

func (*hashService) HashPassword(password string) (string, error) {
	hashed, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hashed, nil
}

func (*hashService) VerifyPassword(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, err
	}
	return match, nil
}

// This is the HashService instance
var HashService HashServiceInterface = NewHashService()
