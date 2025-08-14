package main

import (
	"log"

	"github.com/ducklawrence/go-ecom/cmd/api"
	"github.com/ducklawrence/go-ecom/config"
	"github.com/ducklawrence/go-ecom/db"
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

	server := api.NewAPIServer(":8080", db)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
