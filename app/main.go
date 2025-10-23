package main

import (
	"log"
	"podGopher/adapter/outbound/repository/postgres/migration"
	"podGopher/env"
)

func main() {
	loadEnvironment()
	startMigration()
}

func loadEnvironment() {
	if err := env.Load("env/.env"); err != nil {
		log.Fatal(err)
	}
}

func startMigration() {
	dbMigration, err := migration.NewMigration()
	if err != nil {
		log.Fatal(err)
	}
	if err := dbMigration.Migrate(); err != nil {
		log.Printf("WARNING on migration: %s", err)
	}
}
