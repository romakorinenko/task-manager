package dbpool

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/romakorinenko/task-manager/internal/config"
)

func NewDBPool(ctx context.Context, dBCfg *config.DB) (*pgxpool.Pool, error) {
	DBConfig, err := pgxpool.ParseConfig(dBCfg.ConnectionString)
	if err != nil {
		return nil, err
	}
	return pgxpool.NewWithConfig(ctx, DBConfig)
}
