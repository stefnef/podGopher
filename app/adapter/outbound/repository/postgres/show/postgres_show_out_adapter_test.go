package show

import (
	repositoryEpisode "podGopher/adapter/outbound/repository/postgres/episode"
	"podGopher/adapter/outbound/repository/postgres/postgresTestSetup"
	"podGopher/core/domain/model"
	"podGopher/core/port/outbound"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_should_implement_port(t *testing.T) {
	repository := NewPostgresShowRepository(nil)

	assert.NotNil(t, repository)
	assert.Implements(t, (*outbound.SaveShowPort)(nil), repository)
}

func Test_should_save_a_show(t *testing.T) {
	db := postgresTestSetup.StartTestcontainersPostgres(t, "../postgresTestSetup/")

	defer postgresTestSetup.Teardown(t, db)

	repository := NewPostgresShowRepository(db)
	showTitle := "Some title"
	showSlug := showTitle + "-Slug"
	show := &model.Show{
		Id:       uuid.NewString(),
		Title:    showTitle,
		Slug:     showSlug,
		Episodes: []string{},
	}

	t.Run("should return false if show with title or slug does not exist", func(t *testing.T) {
		exists := repository.ExistsByTitleOrSlug(showTitle, showSlug)
		assert.False(t, exists)
	})

	t.Run("should save a show", func(t *testing.T) {
		err := repository.SaveShow(show)
		assert.Nil(t, err)
	})

	t.Run("should return true if show with title exists", func(t *testing.T) {
		exists := repository.ExistsByTitleOrSlug(showTitle, "some-other-slug")
		assert.True(t, exists)
	})

	t.Run("should return true if show with slug exists", func(t *testing.T) {
		exists := repository.ExistsByTitleOrSlug("some-other-title", showSlug)
		assert.True(t, exists)
	})

	t.Run("should return true if show with title and slug exists", func(t *testing.T) {
		exists := repository.ExistsByTitleOrSlug(showTitle, showSlug)
		assert.True(t, exists)
	})

	t.Run("should return false if show with title or slug does not exists", func(t *testing.T) {
		exists := repository.ExistsByTitleOrSlug("some-other-title", "some-other-slug")
		assert.False(t, exists)
	})

	t.Run("should query a show", func(t *testing.T) {
		var id string
		var title string
		var slug string
		err := db.QueryRow("SELECT * FROM show WHERE id = $1", show.Id).
			Scan(&id, &title, &slug)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, show.Id, id)
		assert.Equal(t, show.Title, title)
		assert.Equal(t, show.Slug, slug)
	})
}

func Test_should_retrieve_a_show(t *testing.T) {
	db := postgresTestSetup.StartTestcontainersPostgres(t, "../postgresTestSetup/")

	defer postgresTestSetup.Teardown(t, db)

	repository := NewPostgresShowRepository(db)
	show := &model.Show{
		Id:    uuid.NewString(),
		Title: "Some title",
		Slug:  ("Some title") + "-Slug",
	}

	err := repository.SaveShow(show)
	assert.Nil(t, err)

	t.Run("should return nil if show does not exist", func(t *testing.T) {
		foundShow, err := repository.GetShowOrNil(uuid.NewString())
		assert.Nil(t, err)
		assert.Nil(t, foundShow)
	})

	t.Run("should retrieve a show", func(t *testing.T) {
		foundShow, err := repository.GetShowOrNil(show.Id)
		assert.Nil(t, err)
		assert.NotNil(t, foundShow)
		assert.Equal(t, show.Id, foundShow.Id)
		assert.Empty(t, foundShow.Episodes)
	})

}

func Test_should_reference_episodes(t *testing.T) {
	db := postgresTestSetup.StartTestcontainersPostgres(t, "../postgresTestSetup/")

	defer postgresTestSetup.Teardown(t, db)

	showRepository := NewPostgresShowRepository(db)

	episodeRepository := repositoryEpisode.NewPostgresEpisodeRepository(db)
	showWithEpisodes := &model.Show{
		Id:    uuid.NewString(),
		Title: "first show",
		Slug:  "first-show-Slug",
	}
	showWithoutEpisodes := &model.Show{
		Id:    uuid.NewString(),
		Title: "show",
		Slug:  "show-Slug",
	}

	err := showRepository.SaveShow(showWithEpisodes)
	assert.Nil(t, err)

	err = showRepository.SaveShow(showWithoutEpisodes)
	assert.Nil(t, err)

	err = episodeRepository.SaveEpisode(&model.Episode{
		Id:     uuid.NewString(),
		ShowId: showWithEpisodes.Id,
		Title:  "first episode",
	})
	assert.Nil(t, err)

	err = episodeRepository.SaveEpisode(&model.Episode{
		Id:     uuid.NewString(),
		ShowId: showWithEpisodes.Id,
		Title:  "first episode",
	})
	assert.Nil(t, err)

	t.Run("should retrieve a show with episodes", func(t *testing.T) {
		foundShow, err := showRepository.GetShowOrNil(showWithEpisodes.Id)

		assert.Nil(t, err)
		assert.NotNil(t, foundShow)
		assert.Len(t, foundShow.Episodes, 2)
	})

	t.Run("should not retrieve non-referenced episodes", func(t *testing.T) {
		foundShow, err := showRepository.GetShowOrNil(showWithoutEpisodes.Id)

		assert.Nil(t, err)
		assert.NotNil(t, foundShow)
		assert.Len(t, foundShow.Episodes, 0)
	})
}
