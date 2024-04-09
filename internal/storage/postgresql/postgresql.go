package postgresql

import (
	repeatable "app/internal/utils"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type postgres struct {
	*pgxpool.Pool
}

func GetPgxPool(PGLD string, maxAttempts int) (pool *pgxpool.Pool, err error) {
	const el = "postgresql.postgresql.GetPgxPool"
	config, err := pgxpool.ParseConfig(PGLD)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", el, err)
	}

	err = repeatable.DoWithTries(func() error {

		pool, err = pgxpool.NewWithConfig(context.Background(), config)
		if err != nil {
			return fmt.Errorf("%s: %w", el, err)
		}

		err = pool.Ping(context.Background())
		if err != nil {
			return fmt.Errorf("%s: %w", el, err)
		}

		return nil
	}, maxAttempts, 5*time.Second)

	if err != nil {
		log.Fatal("error do with tries postgresql")
	}

	return pool, nil
}
