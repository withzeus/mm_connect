package repository

import (
	"context"

	"github.com/withzeus/mm_connect/internal/domain"
)

type ClientAuthRepository interface {
	GetByClientID(ctx context.Context, clientID string) (*domain.Client, error)

	RegisterClient(ctx context.Context, client *domain.Client) (*domain.Client, error)
}
