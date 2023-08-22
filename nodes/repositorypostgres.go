package nodes

import "github.com/jackc/pgx/v4/pgxpool"

type RepositoryPostgres struct {
	conn *pgxpool.Pool
}

func NewRepositoryPostgres(conn *pgxpool.Pool) *RepositoryPostgres {
	return &RepositoryPostgres{conn: conn}
}
