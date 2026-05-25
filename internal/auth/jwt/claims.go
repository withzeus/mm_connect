package jwt

import (
	jwtv5 "github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	ClientName string   `json:"client_name"`
	Scopes     []string `json:"scopes"`
	Type       string   `json:"Type"`
	jwtv5.RegisteredClaims
}
