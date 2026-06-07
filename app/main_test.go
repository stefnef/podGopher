package main

import (
	"bytes"
	"database/sql"
	"net/http"
	"os"
	oktaAuth "podGopher/adapter/outbound/credentials/okta/user"
	"podGopher/adapter/outbound/repository/postgres/postgresTestSetup"
	"podGopher/env"
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

func Test_should_load_context(t *testing.T) {
	setup(t)

	defer postgresTestSetup.Teardown(t, db)
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

func Test_should_load_credential_service(t *testing.T) {
	defer func(key string) {
		err := os.Unsetenv(key)
		assert.Nil(t, err)
	}(string(env.CredentialService))

	t.Run("should load OAuth2", func(t *testing.T) {
		err := os.Setenv(string(env.CredentialService), "OAuth2")
		assert.Nil(t, err)

		var servicePorts = &servicePorts{}
		servicePorts.initCredentialService()
		assert.IsType(t, &oktaAuth.OktaUserOutAdapter{}, servicePorts.createUserCredentialsPort)
	})

	t.Run("should load mock", func(t *testing.T) {
		err := os.Setenv(string(env.CredentialService), "None")
		assert.Nil(t, err)

		var servicePorts = &servicePorts{}
		servicePorts.initCredentialService()

		assert.Nil(t, servicePorts.createUserCredentialsPort)
	})

	t.Run("should throw error on misconfiguration", func(t *testing.T) {
		_ = os.Setenv(string(env.CredentialService), "Invalid")
		assert.Panics(t, func() {
			new(servicePorts).initCredentialService()
		})
	})
}
