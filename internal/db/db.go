package db

import (
	"context"

	"github.com/4serviceSoftware/tech-task/internal/config"
	"github.com/jackc/pgx/v4/pgxpool"
)

func GetPostgresConnection(ctx context.Context, config *config.Config) (*pgxpool.Pool, error) {
	dbUrl := "postgres://" + config.DbUser + ":" + config.DbPass + "@" + config.DbHost + ":" + config.DbPort + "/" + config.DbName
	db, err := pgxpool.Connect(ctx, dbUrl)
	if err != nil {
		return nil, err
	}
	return db, nil
}
