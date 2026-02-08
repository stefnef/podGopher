package rss

import (
	repositoryDistribution "podGopher/adapter/outbound/repository/postgres/distribution"
	repositoryEpisode "podGopher/adapter/outbound/repository/postgres/episode"
	"podGopher/adapter/outbound/repository/postgres/postgresTestSetup"
	repositoryShow "podGopher/adapter/outbound/repository/postgres/show"
	"podGopher/core/domain/model"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_should_retrieve_rss(t *testing.T) {
	db := postgresTestSetup.StartTestcontainersPostgres(t, "../postgresTestSetup/")

	defer postgresTestSetup.Teardown(t, db)

	showUuid := uuid.NewString()
	firstEpisodeUuid := uuid.NewString()
	secondEpisodeUuid := uuid.NewString()

	showRepository := repositoryShow.NewPostgresShowRepository(db)

	distributionRepository := repositoryDistribution.NewPostgresDistributionRepository(db)
	episodeRepository := repositoryEpisode.NewPostgresEpisodeRepository(db)

	rssRepository := NewPostgresRSSRepository(db)

	show := &model.Show{
		Id:    showUuid,
		Title: "Some title",
		Slug:  "Some-Slug",
	}
	distribution := &model.Distribution{
		Id:     uuid.NewString(),
		ShowId: showUuid,
		Title:  "Some title",
		Slug:   ("Some title") + "-Slug",
	}
	firstEpisode := &model.Episode{Id: firstEpisodeUuid, ShowId: showUuid, Title: "first-episode"}
	secondEpisode := &model.Episode{Id: secondEpisodeUuid, ShowId: showUuid, Title: "first-episode"}

	err := showRepository.SaveShow(show)
	assert.Nil(t, err)

	err = distributionRepository.SaveDistribution(distribution)
	assert.Nil(t, err)

	if err := episodeRepository.SaveEpisode(firstEpisode); err != nil {
		t.Fatal(err)
	}
	if err := episodeRepository.SaveEpisode(secondEpisode); err != nil {
		t.Fatal(err)
	}

	t.Run("should return nil if show slug does not exist", func(t *testing.T) {
		foundRss, err := rssRepository.GetRSSOrNil(uuid.NewString(), distribution.Slug)
		assert.Nil(t, err)
		assert.Nil(t, foundRss)
	})

	t.Run("should return nil if distribution slug does not exist", func(t *testing.T) {
		foundRss, err := rssRepository.GetRSSOrNil(show.Slug, uuid.NewString())
		assert.Nil(t, err)
		assert.Nil(t, foundRss)
	})

	t.Run("should retrieve a rss", func(t *testing.T) {
		foundRss, err := rssRepository.GetRSSOrNil(show.Slug, distribution.Slug)
		assert.Nil(t, err)
		assert.NotNil(t, foundRss)

		assert.Equal(t, show, foundRss.Show)
		assert.Equal(t, distribution, foundRss.Distribution)
		assert.Empty(t, foundRss.Episodes)
	})

	t.Run("should include episodes", func(t *testing.T) {
		distribution.Episodes = []string{firstEpisodeUuid, secondEpisodeUuid}
		if err := distributionRepository.UpdateDistribution(distribution); err != nil {
			t.Fatal(err)
		}

		foundRss, err := rssRepository.GetRSSOrNil(show.Slug, distribution.Slug)
		assert.Nil(t, err)
		assert.NotNil(t, foundRss)

		assert.Equal(t, show, foundRss.Show)
		assert.Equal(t, distribution, foundRss.Distribution)
		assert.ElementsMatch(t, []*model.Episode{firstEpisode, secondEpisode}, foundRss.Episodes)
	})
}
