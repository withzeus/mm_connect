package service

import (
	"context"
	"errors"
	"slices"

	"github.com/withzeus/mm_connect/internal/auth/hashing"
	"github.com/withzeus/mm_connect/internal/auth/jwt"
	"github.com/withzeus/mm_connect/internal/auth/secret"
	"github.com/withzeus/mm_connect/internal/auth/uuid"
	"github.com/withzeus/mm_connect/internal/domain"
	"github.com/withzeus/mm_connect/internal/repository"
)

type AuthService struct {
	clientRepo repository.ClientAuthRepository
	jwt        *jwt.Manager
}

func NewAuthService(clientRepo repository.ClientAuthRepository, jwt *jwt.Manager) *AuthService {
	return &AuthService{
		clientRepo: clientRepo,
		jwt:        jwt,
	}
}

func (s *AuthService) RegisterClient(ctx context.Context, clientName string, website string) (*domain.Client, string, error) {
	ps, err := secret.Generate()
	if err != nil {
		return nil, "", errors.New("Internal server error")
	}

	sh, err := hashing.Hash(ps)
	if err != nil {
		return nil, "", errors.New("Internal server error")
	}

	uuidb := uuid.GenerateV4()

	client := &domain.Client{
		UUID:    uuidb,
		Name:    clientName,
		Domain:  website,
		Secret:  sh,
		Enabled: true,
	}

	created, err := s.clientRepo.RegisterClient(ctx, client)
	if err != nil {
		return nil, "", errors.New(err.Error())
	}

	return created, ps, nil
}

func (s *AuthService) IssueToken(ctx context.Context, clientID string, clientSecret string) (string, error) {
	client, err := s.clientRepo.GetByClientID(ctx, clientID)

	if err != nil {
		return "", err
	}

	if !client.Enabled {
		return "", errors.New("Invalid client")
	}

	if err := hashing.Verify(client.Secret, clientSecret); err != nil {
		return "", errors.New("Invalid credentials")
	}

	token, err := s.jwt.GenerateClientToken(
		client.Name,
		client.Scopes,
	)

	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) ValidateToken(ctx context.Context, token string) (*jwt.Claims, error) {
	return s.jwt.Verify(token)
}

func (s *AuthService) Authorize(ctx context.Context, token string, requiredScope string) (bool, error) {
	claims, err := s.jwt.Verify(token)

	if err != nil {
		return false, err
	}

	if slices.Contains(claims.Scopes, requiredScope) {
		return true, nil
	}

	return false, nil
}
