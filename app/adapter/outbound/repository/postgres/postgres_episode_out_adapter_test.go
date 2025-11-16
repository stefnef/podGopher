package repository

import (
	"podGopher/adapter/outbound/repository/postgres/postgresTestSetup"
	"podGopher/core/domain/model"
	"podGopher/core/port/outbound"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_episode_repository_should_implement_port(t *testing.T) {
	repository := NewPostgresEpisodeRepository(nil)

	assert.NotNil(t, repository)
	assert.Implements(t, (*outbound.SaveEpisodePort)(nil), repository)
}

func Test_should_not_save_episode_if_show_does_not_exist(t *testing.T) {
	db := postgresTestSetup.StartTestcontainersPostgres(t, "postgresTestSetup/")
	defer teardown(t, db)

	nonExistingShowId := uuid.NewString()

	repository := NewPostgresEpisodeRepository(db)
	episode := &model.Episode{
		Id:     uuid.NewString(),
		ShowId: nonExistingShowId,
		Title:  "test-Title",
	}
	err := repository.SaveEpisode(episode)
	assert.NotNil(t, err)
}

func Test_should_save_an_episode(t *testing.T) {
	db := postgresTestSetup.StartTestcontainersPostgres(t, "postgresTestSetup/")
	defer teardown(t, db)

	showUuid := uuid.NewString()

	showRepository := NewPostgresShowRepository(db)
	repository := NewPostgresEpisodeRepository(db)
	episodeTitle := "Some title"
	episode := &model.Episode{
		Id:     uuid.NewString(),
		ShowId: showUuid,
		Title:  episodeTitle,
	}

	if err := showRepository.SaveShow(&model.Show{Id: showUuid, Title: "test-show", Slug: "test-slug"}); err != nil {
		t.Fatal(err)
	}

	t.Run("should return false if episode with title does not exist", func(t *testing.T) {
		exists := repository.ExistsByTitle(episodeTitle)
		assert.False(t, exists)
	})

	t.Run("should save an episode", func(t *testing.T) {
		err := repository.SaveEpisode(episode)
		assert.Nil(t, err)
	})

	t.Run("should return true if episode with title exists", func(t *testing.T) {
		exists := repository.ExistsByTitle(episodeTitle)
		assert.True(t, exists)
	})

	t.Run("should query an episode", func(t *testing.T) {
		var id string
		var title string
		var showId string
		err := db.QueryRow("SELECT * FROM episode WHERE id = $1", episode.Id).
			Scan(&id, &showId, &title)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, episode.Id, id)
		assert.Equal(t, episode.Title, title)
		assert.Equal(t, episode.ShowId, showId)
	})
}
