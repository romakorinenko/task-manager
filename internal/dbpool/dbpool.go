package dbpool

import (
	"context"
	"github.com/romakorinenko/task-manager/internal/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDBPool(ctx context.Context, dBCfg *config.DB) (*pgxpool.Pool, error) {
	DBConfig, err := pgxpool.ParseConfig(dBCfg.ConnectionString)
	if err != nil {
		return nil, err
	}
	return pgxpool.NewWithConfig(ctx, DBConfig)
}
