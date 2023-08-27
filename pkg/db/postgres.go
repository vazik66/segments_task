package pkg

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

func NewPostgresClient(ctx context.Context, dsn string) (*pgx.Conn, error) {
	pgConn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}

	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err = pgConn.Ping(ctxTimeout)
	if err != nil {
		return nil, fmt.Errorf("connection time exceeded, %v", err)
	}
	return pgConn, nil
}
