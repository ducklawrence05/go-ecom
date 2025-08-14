package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/microsoft/go-mssqldb"
)

type SQLServerConfig struct {
	User     string
	Password string
	Host     string
	Port     int
	DBName   string
}

type DBTX interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func (c *SQLServerConfig) FormatDSN() string {
	return fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s",
		c.User, c.Password, c.Host, c.Port, c.DBName)
}

func NewSQLServerStorage(cfg *SQLServerConfig) (*sql.DB, error) {
	dsn := cfg.FormatDSN()
	db, err := sql.Open("sqlserver", dsn)

	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	return db, err
}
