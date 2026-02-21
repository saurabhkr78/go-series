package config

import (
	"fmt"
	"os"
)

func getDSN() (string, error) {
	dsn := os.Getenv("PGSQL_DATABASE_URL")
	if dsn == "" {
		return "", fmt.Errorf("PGSQL_DATABASE_URL not set")
	}
	return dsn, nil
}
