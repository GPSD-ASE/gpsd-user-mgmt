package db

import (
	"context"
	"fmt"
	"gpsd-user-mgmt/src/config"
	"log/slog"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

func Connect() bool {
	connString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		config.DB_USER,
		config.DB_PASS,
		config.DB_HOST,
		config.DB_PORT,
		config.DB_NAME,
	)

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

// For Testing
const (
	create_table = `CREATE TABLE IF NOT EXISTS users (
		id SERIAL NOT NULL PRIMARY KEY,
		name       varchar(40),
		deviceID    int,
		role 	   varchar(40),
		createdAt timestamp,
		updatedAt timestamp
	);`

	delete_users = "DELETE FROM users"
)

func CreateDatabase() {
	_, err := Pool.Exec(context.Background(), create_table)
	if err != nil {
		slog.Error(err.Error())
	}
}

func EmptyDatabase() {
	_, err := Pool.Exec(context.Background(), delete_users)
	if err != nil {
		slog.Error(err.Error())
	}
}

// func (p *connPool) Select(tableName string) {}
// func (p *connPool) Insert(tableName string) {}
// func (p *connPool) Update(tableName string) {}
// func (p *connPool) Delete(tableName string) {}
