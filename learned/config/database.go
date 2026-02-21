package config

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

func ConnectDB() error {
	dsn, err := getDSN()
	if err != nil {
		return err
	}
	db, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		return err
	}
	if err = db.Ping(context.Background()); err != nil {
		return err
	}
	return nil
}

func GetDB() *pgxpool.Pool {
	if db == nil {
		panic("You MUST call ConnectDB() first then call GetDB()")
	}
	return db
}

func CloseDB() {
	if db != nil {
		db.Close()
	}
}
