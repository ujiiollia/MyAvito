package services

import (
	"app/internal/models"
	"context"
	"fmt"
)

var (
	_ models.BannerService = (*banner)(nil)
)

type banner struct {
	repo models.BannerRepository
}

func NewBanner(repo models.BannerRepository) *banner {
	return &banner{repo: repo}
}

func (s *banner) Ping(ctx context.Context) error {

	err := s.repo.Ping(ctx)
	if err != nil {
		return fmt.Errorf("service.Ping: %w", err)
	}

	return nil
}
