package distribution

import (
	"podGopher/adapter/outbound/repository/postgres/postgresTestSetup"
	repositoryShow "podGopher/adapter/outbound/repository/postgres/show"
	"podGopher/core/domain/model"
	"podGopher/core/port/outbound"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_distribution_repository_should_implement_port(t *testing.T) {
	repository := NewPostgresDistributionRepository(nil)

	assert.NotNil(t, repository)
	assert.Implements(t, (*outbound.SaveDistributionPort)(nil), repository)
}

func Test_should_not_save_distribution_if_show_does_not_exist(t *testing.T) {
	db := postgresTestSetup.StartTestcontainersPostgres(t, "../postgresTestSetup/")
	defer postgresTestSetup.Teardown(t, db)

	nonExistingShowId := uuid.NewString()

	repository := NewPostgresDistributionRepository(db)
	distribution := &model.Distribution{
		Id:     uuid.NewString(),
		ShowId: nonExistingShowId,
		Title:  "test-Title",
	}
	err := repository.SaveDistribution(distribution)
	assert.NotNil(t, err)
}

func Test_should_save_a_distribution(t *testing.T) {
	db := postgresTestSetup.StartTestcontainersPostgres(t, "../postgresTestSetup/")
	defer postgresTestSetup.Teardown(t, db)

	showUuid := uuid.NewString()

	showRepository := repositoryShow.NewPostgresShowRepository(db)
	repository := NewPostgresDistributionRepository(db)
	distributionTitle := "Some title"
	distributionSlug := distributionTitle + "-Slug"
	distribution := &model.Distribution{
		Id:     uuid.NewString(),
		ShowId: showUuid,
		Slug:   distributionSlug,
		Title:  distributionTitle,
	}

	if err := showRepository.SaveShow(&model.Show{Id: showUuid, Title: "test-show", Slug: "test-slug"}); err != nil {
		t.Fatal(err)
	}

	t.Run("should return false if distribution with title or slug does not exist", func(t *testing.T) {
		exists := repository.ExistsByTitleOrSlug(distributionTitle, distributionSlug)
		assert.False(t, exists)
	})

	t.Run("should save a distribution", func(t *testing.T) {
		err := repository.SaveDistribution(distribution)
		assert.Nil(t, err)
	})

	t.Run("should return true if distribution with title exists", func(t *testing.T) {
		exists := repository.ExistsByTitleOrSlug(distributionTitle, "some-other-slug")
		assert.True(t, exists)
	})

	t.Run("should return true if distribution with slug exists", func(t *testing.T) {
		exists := repository.ExistsByTitleOrSlug("some-other-title", distributionSlug)
		assert.True(t, exists)
	})

	t.Run("should return true if distribution with title and slug exists", func(t *testing.T) {
		exists := repository.ExistsByTitleOrSlug(distributionTitle, distributionSlug)
		assert.True(t, exists)
	})

	t.Run("should return false if distribtion with title or slug does not exists", func(t *testing.T) {
		exists := repository.ExistsByTitleOrSlug("some-other-title", "some-other-slug")
		assert.False(t, exists)
	})

	t.Run("should query a distribution", func(t *testing.T) {
		var id string
		var title string
		var showId string
		var slug string
		err := db.QueryRow("SELECT * FROM distribution WHERE id = $1", distribution.Id).
			Scan(&id, &showId, &title, &slug)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, distribution.Id, id)
		assert.Equal(t, distribution.Title, title)
		assert.Equal(t, distribution.ShowId, showId)
	})
}

func Test_should_retrieve_a_distribution(t *testing.T) {
	db := postgresTestSetup.StartTestcontainersPostgres(t, "../postgresTestSetup/")

	defer postgresTestSetup.Teardown(t, db)

	showUuid := uuid.NewString()
	showRepository := repositoryShow.NewPostgresShowRepository(db)

	repository := NewPostgresDistributionRepository(db)
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

	err := showRepository.SaveShow(show)
	assert.Nil(t, err)

	err = repository.SaveDistribution(distribution)
	assert.Nil(t, err)

	t.Run("should return nil if distribution does not exist", func(t *testing.T) {
		foundDistribution, err := repository.GetDistributionOrNil(uuid.NewString())
		assert.Nil(t, err)
		assert.Nil(t, foundDistribution)
	})

	t.Run("should retrieve a distribution", func(t *testing.T) {
		foundDistribution, err := repository.GetDistributionOrNil(distribution.Id)
		assert.Nil(t, err)
		assert.NotNil(t, foundDistribution)
		assert.Equal(t, show.Id, foundDistribution.ShowId)
		assert.Equal(t, distribution.Id, foundDistribution.Id)
		assert.Equal(t, distribution.Title, foundDistribution.Title)
		assert.Equal(t, distribution.Slug, foundDistribution.Slug)
	})
}
