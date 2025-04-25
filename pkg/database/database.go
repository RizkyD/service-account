package database

import (
	"context"
	zlog "github.com/rs/zerolog/log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func ConnectDB(databaseURL string) *pgxpool.Pool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		zlog.Fatal().Msgf("Unable to create connection pool: %v", err)
	}

	ctxPing, cancelPing := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelPing()
	if err = pool.Ping(ctxPing); err != nil {
		zlog.Fatal().Msgf("Unable to ping database: %v", err)
	}

	zlog.Info().Msg("Successfully connected to the database!")
	return pool
}
