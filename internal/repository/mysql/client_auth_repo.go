package mysql

import (
	"context"
	"database/sql"

	"github.com/withzeus/mm_connect/internal/domain"
)

type ClientAuthRepository struct {
	db *sql.DB
}

func NewClientAuthRepository(db *sql.DB) *ClientAuthRepository {
	return &ClientAuthRepository{db: db}
}

func (r *ClientAuthRepository) GetByClientID(ctx context.Context, clientID string) (*domain.Client, error) {
	query := `SELECT id, BIN_TO_UUID(client_id, 1) as client_id, client_name, client_secret, enabled, domain FROM clients WHERE client_id=UUID_TO_BIN(?, 1)`

	var client domain.Client

	err := r.db.QueryRowContext(ctx, query, clientID).Scan(
		&client.ID,
		&client.UUID,
		&client.Name,
		&client.Secret,
		&client.Enabled,
		&client.Domain,
	)

	if err != nil {
		return nil, err
	}

	queryScope := `SELECT scope FROM client_scopes WHERE client_id=?`

	rows, err := r.db.QueryContext(ctx, queryScope, client.ID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var scope string
		if err := rows.Scan(&scope); err != nil {
			return nil, err
		}

		client.Scopes = append(client.Scopes, scope)
	}

	return &client, nil
}

func (r *ClientAuthRepository) RegisterClient(ctx context.Context, client *domain.Client) (*domain.Client, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	query := `INSERT INTO clients (client_id, client_name, client_secret, domain, enabled) VALUES (UUID_TO_BIN(?, 1), ?, ?, ?, ?)`

	result, err := tx.ExecContext(
		ctx, query, client.UUID, client.Name, client.Secret, client.Domain, client.Enabled,
	)
	if err != nil {
		return nil, err
	}

	client_id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	qScope := `INSERT INTO client_scopes (client_id, scope) VALUES (?, ?)`

	_, err = tx.ExecContext(
		ctx, qScope, client_id, "iam",
	)

	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return client, nil
}
