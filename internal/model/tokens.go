package model

import "github.com/golang-jwt/jwt/v5"

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserToken struct {
	UserID       int64  `json:"user_id"`
	RefreshToken string `json:"refresh_token"`
}

type Payload struct {
	UserID int64
	jwt.RegisteredClaims
}

func (t *Tokens) ValidateAccessToken() error {
	return nil
}
