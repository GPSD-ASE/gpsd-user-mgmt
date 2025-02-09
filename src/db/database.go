package db

import (
	"context"
	"fmt"
	"gpsd-user-mgmt/src/config"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func Connect(config *config.Config) bool {
	connString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		config.DB_USER,
		config.DB_PASS,
		config.DB_HOST,
		config.DB_PORT,
		config.DB_NAME)

	var err error
	Pool, err = pgxpool.New(context.Background(), connString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	}

	return err == nil
}

func Close() {
	if Pool != nil {
		Pool.Close()
	}
}

// func (p *connPool) Select(tableName string) {}
// func (p *connPool) Insert(tableName string) {}
// func (p *connPool) Update(tableName string) {}
// func (p *connPool) Delete(tableName string) {}
