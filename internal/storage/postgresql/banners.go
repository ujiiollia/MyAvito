package postgresql

import (
	"context"
	"fmt"
)

type BannerRepository interface {
	Ping(context.Context) error
}

type banner struct {
	db *postgres
}

func NewBanner(pgl *postgres) *banner {
	return &banner{db: pgl}
}

func (r *banner) Ping(ctx context.Context) error {
	err := r.db.Ping(ctx)
	if err != nil {
		return fmt.Errorf("repository.Ping: %w", err)
	}

	return nil
}
