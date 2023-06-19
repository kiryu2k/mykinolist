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
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

func ParseAccessToken(tokenStr, secretKey string) (int64, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Payload{}, func(token *jwt.Token) (any, error) {
		return []byte(secretKey), nil
	})
	claims, ok := token.Claims.(*Payload)
	if !ok || !token.Valid {
		return 0, err
	}
	return claims.UserID, nil
}

// func ParseRefreshToken
