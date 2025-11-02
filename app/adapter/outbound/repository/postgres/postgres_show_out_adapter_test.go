package repository

import (
	"database/sql"
	"podGopher/adapter/outbound/repository/postgres/postgresTestSetup"
	"podGopher/core/domain/model"
	"podGopher/core/port/outbound"
	"testing"

	"github.com/google/uuid"
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
	showSlug := showTitle + "-Slug"
	show := &model.Show{
		Id:    uuid.NewString(),
		Title: showTitle,
		Slug:  showSlug,
	}

	t.Run("should should return false if show with title or slug does not exist", func(t *testing.T) {
		exists := repository.ExistsByTitleOrSlug(showTitle, showSlug)
		assert.False(t, exists)
	})

	t.Run("should save a show", func(t *testing.T) {
		err := repository.SaveShow(show)
		assert.Nil(t, err)
	})

	t.Run("should should return true if show with title exists", func(t *testing.T) {
		exists := repository.ExistsByTitleOrSlug(showTitle, "some-other-slug")
		assert.True(t, exists)
	})

	t.Run("should should return true if show with slug exists", func(t *testing.T) {
		exists := repository.ExistsByTitleOrSlug("some-other-title", showSlug)
		assert.True(t, exists)
	})

	t.Run("should should return true if show with title and slug exists", func(t *testing.T) {
		exists := repository.ExistsByTitleOrSlug(showTitle, showSlug)
		assert.True(t, exists)
	})

	t.Run("should should return false if show with title or slug does not exists", func(t *testing.T) {
		exists := repository.ExistsByTitleOrSlug("some-other-title", "some-other-slug")
		assert.False(t, exists)
	})

	t.Run("should query a show", func(t *testing.T) {
		var id string
		var title string
		var slug string
		err := db.QueryRow("SELECT * FROM shows WHERE id = $1", show.Id).
			Scan(&id, &title, &slug)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, show.Id, id)
		assert.Equal(t, show.Title, title)
		assert.Equal(t, show.Slug, slug)
	})
}
