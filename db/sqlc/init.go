package db

import (
	"context"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v4"

	"conduit/config"
)

func Connect(config config.Config) (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), config.DBUrl)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func Close(conn *pgx.Conn) error {
	err := conn.Close(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func AutoMigrate(config config.Config) error {
	m, err := migrate.New(config.DBMigrationPath, config.DBUrl)
	if err != nil {
		return err
	}

	if config.DBRecreate {
		if err := m.Down(); err != nil {
			if err != migrate.ErrNoChange {
				return err
			}
		}
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}
