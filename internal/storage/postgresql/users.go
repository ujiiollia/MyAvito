package postgresql

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

type users struct {
	dbu *postgres
}

func NewUser(pgl *postgres) *users {
	return &users{dbu: pgl}
}

func GetRole(userID string, pool *pgxpool.Pool) (string, error) {
	if pool == nil {
		return "", errors.New("database pool was not initialized")
	}
	var is_admin string
	err := pool.QueryRow(context.Background(), "SELECT is_admin FROM users WHERE user_id=$1", userID).Scan(&is_admin)
	if err != nil {
		return "", err
	}
	return is_admin, nil
}
