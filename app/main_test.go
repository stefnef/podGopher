package main

import (
	"bytes"
	"database/sql"
	"net/http"
	"podGopher/adapter/outbound/repository/postgres/postgresTestSetup"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var app *App
var db *sql.DB

func setup(t *testing.T) {
	db = postgresTestSetup.StartTestcontainersPostgres(t, "adapter/outbound/repository/postgres/postgresTestSetup/")
	app = NewApp("env/.testcontainers-env")
}

func teardown(t *testing.T, db *sql.DB) {
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			t.Fatalf("Failed to close database connection: %v", err)
		}
	}(db)

	defer postgresTestSetup.TeardownTestcontainersPostgres(t)
}

func Test_should_load_context(t *testing.T) {
	setup(t)

	defer teardown(t, db)
	defer app.Stop()

	go app.Start()
	time.Sleep(100 * time.Millisecond)

	t.Run("should add a show", func(t *testing.T) {
		postShowRequest := `{"Title":"some title", "Slug":"some slug"}`
		response, err := http.Post("http://localhost:3000/show", "application/json", bytes.NewBuffer([]byte(postShowRequest)))
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, http.StatusCreated, response.StatusCode)
	})
}
