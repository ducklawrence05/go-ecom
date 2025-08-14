package main

import (
	"log"
	"os"

	"github.com/ducklawrence/go-ecom/config"
	"github.com/ducklawrence/go-ecom/db"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlserver"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	db, err := db.NewSQLServerStorage(&db.SQLServerConfig{
		User:     config.Envs.User,
		Password: config.Envs.Password,
		Host:     config.Envs.Host,
		Port:     config.Envs.Port,
		DBName:   config.Envs.DBName,
	})

	if err != nil {
		log.Fatal(err)
	}

	driver, err := sqlserver.WithInstance(db, &sqlserver.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://cmd/migrate/migrations",
		"sqlserver",
		driver,
	)
	if err != nil {
		log.Fatal(err)
	}

	cmd := os.Args[(len(os.Args) - 1)]
	if cmd == "up" {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}
	if cmd == "down" {
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}
}
