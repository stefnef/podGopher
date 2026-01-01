package show

import (
	repositoryDistribution "podGopher/adapter/outbound/repository/postgres/distribution"
	repositoryEpisode "podGopher/adapter/outbound/repository/postgres/episode"
	"podGopher/adapter/outbound/repository/postgres/postgresTestSetup"
	"podGopher/core/domain/model"
	forSaveShow "podGopher/core/port/outbound/show"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_should_implement_port(t *testing.T) {
	repository := NewPostgresShowRepository(nil)

	assert.NotNil(t, repository)
	assert.Implements(t, (*forSaveShow.SaveShowPort)(nil), repository)
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

	t.Run("should return false if forSaveShow with title or slug does not exist", func(t *testing.T) {
		exists := repository.ExistsByTitleOrSlug(showTitle, showSlug)
		assert.False(t, exists)
	})

	t.Run("should save a forSaveShow", func(t *testing.T) {
		err := repository.SaveShow(show)
		assert.Nil(t, err)
	})

	t.Run("should return true if forSaveShow with title exists", func(t *testing.T) {
		exists := repository.ExistsByTitleOrSlug(showTitle, "some-other-slug")
		assert.True(t, exists)
	})

	t.Run("should return true if forSaveShow with slug exists", func(t *testing.T) {
		exists := repository.ExistsByTitleOrSlug("some-other-title", showSlug)
		assert.True(t, exists)
	})

	t.Run("should return true if forSaveShow with title and slug exists", func(t *testing.T) {
		exists := repository.ExistsByTitleOrSlug(showTitle, showSlug)
		assert.True(t, exists)
	})

	t.Run("should return false if forSaveShow with title or slug does not exists", func(t *testing.T) {
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

func Test_should_reference_episodes_and_distributions(t *testing.T) {
	db := postgresTestSetup.StartTestcontainersPostgres(t, "../postgresTestSetup/")

	defer postgresTestSetup.Teardown(t, db)

	showRepository := NewPostgresShowRepository(db)

	episodeRepository := repositoryEpisode.NewPostgresEpisodeRepository(db)
	distributionRepository := repositoryDistribution.NewPostgresDistributionRepository(db)

	showWithEpisodesAndDistributions := &model.Show{
		Id:    uuid.NewString(),
		Title: "first show",
		Slug:  "first-show-Slug",
	}
	showWithoutEpisodesNorDistributions := &model.Show{
		Id:    uuid.NewString(),
		Title: "show",
		Slug:  "show-Slug",
	}

	err := showRepository.SaveShow(showWithEpisodesAndDistributions)
	assert.Nil(t, err)

	err = showRepository.SaveShow(showWithoutEpisodesNorDistributions)
	assert.Nil(t, err)

	err = episodeRepository.SaveEpisode(&model.Episode{
		Id:     uuid.NewString(),
		ShowId: showWithEpisodesAndDistributions.Id,
		Title:  "first episode",
	})
	assert.Nil(t, err)

	err = episodeRepository.SaveEpisode(&model.Episode{
		Id:     uuid.NewString(),
		ShowId: showWithEpisodesAndDistributions.Id,
		Title:  "second episode",
	})
	assert.Nil(t, err)

	err = distributionRepository.SaveDistribution(&model.Distribution{
		Id:     uuid.NewString(),
		ShowId: showWithEpisodesAndDistributions.Id,
		Title:  "1st distribution",
	})
	assert.Nil(t, err)

	err = distributionRepository.SaveDistribution(&model.Distribution{
		Id:     uuid.NewString(),
		ShowId: showWithEpisodesAndDistributions.Id,
		Title:  "2nd distribution",
	})
	assert.Nil(t, err)

	t.Run("should retrieve a show with episodes and distributions", func(t *testing.T) {
		foundShow, err := showRepository.GetShowOrNil(showWithEpisodesAndDistributions.Id)

		assert.Nil(t, err)
		assert.NotNil(t, foundShow)
		assert.Len(t, foundShow.Episodes, 2)
		assert.Len(t, foundShow.Distributions, 2)
	})

	t.Run("should not retrieve non-referenced episodes nor distributions", func(t *testing.T) {
		foundShow, err := showRepository.GetShowOrNil(showWithoutEpisodesNorDistributions.Id)

		assert.Nil(t, err)
		assert.NotNil(t, foundShow)
		assert.Len(t, foundShow.Episodes, 0)
		assert.Len(t, foundShow.Distributions, 0)
	})
}
