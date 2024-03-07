package database

import (
	"context"
	"github.com/jackc/pgx/v5"
	"log/slog"
	"os"
)

type PostgresDb struct {
	*pgx.Conn
}

func NewPostgresDb() *PostgresDb {
	conn, err := pgx.Connect(context.Background(), os.Getenv("POSTGRES_URI"))
	if err != nil {
		slog.Error("error connecting to database", "error", err)
		// lets just panic for now - don't want the application to run without a db connection
		panic(err)
	}
	slog.Info("connected to postgres...")
	return &PostgresDb{
		conn,
	}
}

func (pg *PostgresDb) Get(key string) (string, error) {
	// TODO: need to implement this
	return "", nil
}

func (pg *PostgresDb) Set(key string) error {
	// TODO: need to implement this
	return nil
}

func (pg *PostgresDb) CloseDb() {
	if err := pg.Close(context.Background()); err != nil {
		slog.Error("cannot close postgres connection", "error", err)
		// panic if closing does not work
		panic(err)
	}
}
