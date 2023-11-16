package postgre

import (
	"context"
	"database/sql"

	"github.com/11Petrov/gopherloyal/internal/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

type Database struct {
	db *pgxpool.Pool
}

func NewDBStore(databaseAddress string, ctx context.Context) (*Database, error) {
	log := logger.LoggerFromContext(ctx)

	// Открываем соединение для миграции
	migrationDB, err := sql.Open("pgx", databaseAddress)
	if err != nil {
		log.Errorf("failed to connect for migration: %s", err)
		return nil, err
	}
	defer migrationDB.Close()

	// Проводим миграцию
	log.Info("Start migrating database")
	err = goose.Up(migrationDB, "./internal/migrations")
	if err != nil {
		log.Errorf("error goose.Up: %s", err)
		return nil, err
	}

	// Открываем пул соединений для реальных операций
	db, err := pgxpool.New(context.Background(), databaseAddress)
	if err != nil {
		log.Errorf("failed to connect: %s", err)
		return nil, err
	}

	d := &Database{
		db: db,
	}

	return d, nil
}
