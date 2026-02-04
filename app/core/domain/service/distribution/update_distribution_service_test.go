package distribution

import (
	domainError "podGopher/core/domain/error"
	"podGopher/core/domain/model"
	onUpdateDistribution "podGopher/core/port/inbound/distribution"
	"testing"

	"github.com/stretchr/testify/assert"
)

var updateDistributionService = NewUpdateDistributionService(
	mockGetShowAdapter,
	mockGetEpisodeAdapter,
	mockSaveAndGetDistributionAdapter,
	mockSaveAndGetDistributionAdapter,
)

func Test_should_implement_UpdateDistributionInPort(t *testing.T) {
	assert.NotNil(t, updateDistributionService)
	assert.Implements(t, (*onUpdateDistribution.UpdateDistributionPort)(nil), updateDistributionService)
}

func Test_should_throw_error_if_distribution_with_name_already_exists_on_update(t *testing.T) {
	tests := map[string]struct {
		title                         *string
		slug                          *string
		mockEveryExistsByTitleReturns func()
	}{
		"title and slug given": {
			ptrOfString("Test"),
			ptrOfString("slug"),
			func() {
				mockSaveAndGetDistributionAdapter.everyExistsByTitleReturns("Test", "slug", true)
			},
		},
		"title given": {
			ptrOfString("Test"),
			nil,
			func() {
				mockSaveAndGetDistributionAdapter.everyExistsByTitleReturns("Test", "", true)
			},
		},
		"slug given": {
			nil,
			ptrOfString("slug"),
			func() {
				mockSaveAndGetDistributionAdapter.everyExistsByTitleReturns("", "slug", true)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			defer initAdapter()

			tc.mockEveryExistsByTitleReturns()
			mockGetShowAdapter.returnsOnGetOrNilShow["test-show-id"] = &model.Show{}
			mockSaveAndGetDistributionAdapter.returnsOnGetDistributionOrNil["distribution-id"] = &model.Distribution{ShowId: "test-show-id"}

			command := &onUpdateDistribution.UpdateDistributionCommand{
				DistributionId: "distribution-id",
				ShowId:         "test-show-id",
				Title:          tc.title,
				Slug:           tc.slug,
			}

			err := updateDistributionService.UpdateDistribution(command)

			assert.NotNil(t, err)
			assert.Equal(t, domainError.NewUpdateError("Distribution with title or slug already exists"), err)
			assert.Equal(t, 0, mockSaveAndGetDistributionAdapter.calledSave)
		})
	}
}

func Test_should_throw_error_if_show_does_not_exist_on_update(t *testing.T) {
	defer initAdapter()
	command := &onUpdateDistribution.UpdateDistributionCommand{
		ShowId:         "non-existing-show-id",
		DistributionId: "some-distribution-id",
	}

	mockSaveAndGetDistributionAdapter.returnsOnGetDistributionOrNil["some-distribution-id"] = &model.Distribution{ShowId: "non-existing-show-id"}
	mockGetShowAdapter.returnsOnGetOrNilShow["non-existing-show-id"] = nil

	err := updateDistributionService.UpdateDistribution(command)

	assert.NotNil(t, err)
	assert.Equal(t, domainError.NewShowNotFoundError("non-existing-show-id"), err)
	assert.Equal(t, 1, mockSaveAndGetDistributionAdapter.calledGet)
	assert.Equal(t, 0, mockSaveAndGetDistributionAdapter.calledUpdate)
}

func Test_should_throw_error_if_show_id_was_not_given_on_update(t *testing.T) {
	defer initAdapter()
	command := &onUpdateDistribution.UpdateDistributionCommand{
		ShowId:         "",
		DistributionId: "distribution-id",
	}

	mockSaveAndGetDistributionAdapter.returnsOnGetDistributionOrNil["distribution-id"] = &model.Distribution{}
	err := updateDistributionService.UpdateDistribution(command)

	assert.NotNil(t, err)
	assert.Equal(t, domainError.NewShowNotFoundError(""), err)
	assert.Equal(t, 1, mockSaveAndGetDistributionAdapter.calledGet)
	assert.Equal(t, 0, mockSaveAndGetDistributionAdapter.calledUpdate)
}

func Test_should_throw_error_if_distribution_does_not_exist_on_update(t *testing.T) {
	defer initAdapter()
	command := &onUpdateDistribution.UpdateDistributionCommand{
		ShowId:         "show-id",
		DistributionId: "non-existing-distribution-id",
	}

	mockGetShowAdapter.returnsOnGetOrNilShow["show-id"] = &model.Show{}
	mockSaveAndGetDistributionAdapter.returnsOnGetDistributionOrNil["non-existing-distribution-id"] = nil

	err := updateDistributionService.UpdateDistribution(command)

	assert.NotNil(t, err)
	assert.Equal(t, domainError.NewDistributionNotFoundError("non-existing-distribution-id"), err)
	assert.Equal(t, 1, mockSaveAndGetDistributionAdapter.calledGet)
	assert.Equal(t, 0, mockSaveAndGetDistributionAdapter.calledUpdate)
}

func Test_should_throw_error_if_distribution_id_was_not_given_on_update(t *testing.T) {
	defer initAdapter()
	command := &onUpdateDistribution.UpdateDistributionCommand{
		ShowId:         "show-id",
		DistributionId: "",
	}

	mockGetShowAdapter.returnsOnGetOrNilShow["show-id"] = &model.Show{}

	err := updateDistributionService.UpdateDistribution(command)

	assert.NotNil(t, err)
	assert.Equal(t, domainError.NewDistributionNotFoundError(""), err)
	assert.Equal(t, 1, mockSaveAndGetDistributionAdapter.calledGet)
	assert.Equal(t, 0, mockSaveAndGetDistributionAdapter.calledUpdate)
}

func Test_should_throw_error_if_episode_does_not_exist_on_update(t *testing.T) {
	defer initAdapter()
	command := &onUpdateDistribution.UpdateDistributionCommand{
		ShowId:         "show-id",
		DistributionId: "distribution-id",
		Episodes:       &[]string{"e-id-first", "e-id-second", "e-id-third"},
	}

	mockGetShowAdapter.returnsOnGetOrNilShow["show-id"] = &model.Show{}
	mockSaveAndGetDistributionAdapter.returnsOnGetDistributionOrNil["distribution-id"] = &model.Distribution{ShowId: "show-id"}
	mockGetEpisodeAdapter.returnsOnGetEpisodeOrNil["e-id-first"] = &model.Episode{ShowId: "show-id"}
	mockGetEpisodeAdapter.returnsOnGetEpisodeOrNil["e-id-second"] = nil
	mockGetEpisodeAdapter.returnsOnGetEpisodeOrNil["e-id-third"] = &model.Episode{ShowId: "show-id"}

	err := updateDistributionService.UpdateDistribution(command)

	assert.NotNil(t, err)
	assert.Equal(t, domainError.NewUpdateError("Episode with id 'e-id-second' does not exist"), err)
	assert.Equal(t, 1, mockSaveAndGetDistributionAdapter.calledGet)
	assert.Equal(t, 2, mockGetEpisodeAdapter.called)
	assert.Equal(t, 0, mockSaveAndGetDistributionAdapter.calledUpdate)
}

func Test_should_throw_error_if_episode_does_not_belong_to_updated_show_on_update(t *testing.T) {
	defer initAdapter()
	command := &onUpdateDistribution.UpdateDistributionCommand{
		ShowId:         "show-id",
		DistributionId: "distribution-id",
		Episodes:       &[]string{"e-id"},
	}

	mockGetShowAdapter.returnsOnGetOrNilShow["show-id"] = &model.Show{}
	mockSaveAndGetDistributionAdapter.returnsOnGetDistributionOrNil["distribution-id"] = &model.Distribution{ShowId: "show-id"}
	mockGetEpisodeAdapter.returnsOnGetEpisodeOrNil["e-id"] = &model.Episode{ShowId: "other-show-id"}

	err := updateDistributionService.UpdateDistribution(command)

	assert.NotNil(t, err)
	assert.Equal(t, domainError.NewUpdateError("Episode with id 'e-id' does not belong to show with id 'show-id'"), err)
	assert.Equal(t, 0, mockSaveAndGetDistributionAdapter.calledUpdate)
}

func Test_should_throw_error_if_show_does_not_belong_to_distribution_on_update(t *testing.T) {
	defer initAdapter()
	command := &onUpdateDistribution.UpdateDistributionCommand{
		ShowId:         "show-id",
		DistributionId: "distribution-id",
	}

	mockGetShowAdapter.returnsOnGetOrNilShow["show-id"] = &model.Show{}
	mockSaveAndGetDistributionAdapter.returnsOnGetDistributionOrNil["distribution-id"] = &model.Distribution{ShowId: "other-show-id"}

	err := updateDistributionService.UpdateDistribution(command)

	assert.NotNil(t, err)
	assert.Equal(t, domainError.NewDistributionNotFoundError("distribution-id"), err)
}

func Test_should_update_distribution(t *testing.T) {
	defer initAdapter()

	tests := map[string]struct {
		title            *string
		expectedTitle    string
		slug             *string
		expectedSlug     string
		episodes         *[]string
		expectedEpisodes []string
	}{
		"all fields updated": {
			ptrOfString("new title"),
			"new title",
			ptrOfString("new slug"),
			"new slug",
			&[]string{"new-episode-id"},
			[]string{"new-episode-id"},
		},

		"all nil": {
			nil,
			"old title",
			nil,
			"old slug",
			nil,
			[]string{"e-id"},
		},

		"episodes cleared": {
			nil,
			"old title",
			nil,
			"old slug",
			&[]string{},
			[]string{},
		},
	}

	oldDistribution := &model.Distribution{
		Id:       "distribution-id",
		ShowId:   "show-id",
		Title:    "old title",
		Slug:     "old slug",
		Episodes: []string{"e-id"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			initAdapter()

			command := &onUpdateDistribution.UpdateDistributionCommand{
				ShowId:         "show-id",
				DistributionId: "distribution-id",
				Title:          tc.title,
				Slug:           tc.slug,
				Episodes:       tc.episodes,
			}

			expectedDistribution := &model.Distribution{
				Id:       "distribution-id",
				ShowId:   "show-id",
				Title:    tc.expectedTitle,
				Slug:     tc.expectedSlug,
				Episodes: tc.expectedEpisodes,
			}

			mockGetShowAdapter.returnsOnGetOrNilShow["show-id"] = &model.Show{}
			mockSaveAndGetDistributionAdapter.returnsOnGetDistributionOrNil["distribution-id"] = oldDistribution
			mockGetEpisodeAdapter.returnsOnGetEpisodeOrNil["e-id"] = &model.Episode{ShowId: "show-id"}
			mockGetEpisodeAdapter.returnsOnGetEpisodeOrNil["new-episode-id"] = &model.Episode{ShowId: "show-id"}

			err := updateDistributionService.UpdateDistribution(command)

			assert.Nil(t, err)
			assert.Equal(t, 1, mockSaveAndGetDistributionAdapter.calledUpdate)
			assert.Equal(t, expectedDistribution, mockSaveAndGetDistributionAdapter.calledUpdateWith)
		})
	}
}
