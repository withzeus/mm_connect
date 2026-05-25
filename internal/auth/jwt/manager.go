package jwt

import (
	"errors"
	"time"

	jwtv5 "github.com/golang-jwt/jwt/v5"
)

type Manager struct {
	secret string
}

func New(secret string) *Manager {
	return &Manager{secret: secret}
}

func (m *Manager) GenerateClientToken(clientName string, scopes []string) (string, error) {
	claims := Claims{
		ClientName: clientName,
		Scopes:     scopes,
		Type:       "service",
		RegisteredClaims: jwtv5.RegisteredClaims{
			ExpiresAt: jwtv5.NewNumericDate(
				time.Now().Add(24 * time.Hour),
			),
			Issuer: "mm_connect",
		},
	}

	token := jwtv5.NewWithClaims(
		jwtv5.SigningMethodHS256,
		claims,
	)

	return token.SignedString([]byte(m.secret))
}

func (m *Manager) Verify(ts string) (*Claims, error) {
	token, err := jwtv5.ParseWithClaims(
		ts,
		&Claims{},
		func(token *jwtv5.Token) (any, error) {
			return []byte(m.secret), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("Invalid token")
	}

	return claims, nil
}
