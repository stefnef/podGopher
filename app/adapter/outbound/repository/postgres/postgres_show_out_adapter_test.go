package repository

import (
	"database/sql"
	"podGopher/adapter/outbound/repository/postgres/postgresTestSetup"
	"podGopher/core/port/outbound"
	"testing"

	"github.com/stretchr/testify/assert"
)

func teardown(t *testing.T, db *sql.DB) {
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			t.Fatalf("Failed to close database connection: %v", err)
		}
	}(db)

	defer postgresTestSetup.TeardownTestcontainersPostgres(t)
}

func Test_should_implement_port(t *testing.T) {
	repository := NewPostgresShowRepository(nil)

	assert.NotNil(t, repository)
	assert.Implements(t, (*outbound.SaveShowPort)(nil), repository)
}

func Test_should_save_a_show(t *testing.T) {
	db := postgresTestSetup.StartTestcontainersPostgres(t, "postgresTestSetup/")

	defer teardown(t, db)

	repository := NewPostgresShowRepository(db)
	showTitle := "Some title"

	t.Run("should should return false if show does not exist", func(t *testing.T) {
		exists := repository.ExistsByTitle(showTitle)
		assert.False(t, exists)
	})

	t.Run("should save a show", func(t *testing.T) {
		err := repository.SaveShow(showTitle)
		assert.Nil(t, err)
	})

	t.Run("should should return true if show exists", func(t *testing.T) {
		exists := repository.ExistsByTitle(showTitle)
		assert.True(t, exists)
	})
}
