package episode

import (
	"errors"
	error2 "podGopher/core/domain/error"
	"podGopher/core/domain/model"
	"podGopher/core/port/inbound"
	"testing"

	"github.com/stretchr/testify/assert"
)

var getEpisodeService = NewGetEpisodeService(mockGetShowAdapter, mockSaveAndGetEpisodeAdapter)

func Test_should_implement_GetEpisodeInPort(t *testing.T) {
	assert.NotNil(t, getEpisodeService)
	assert.Implements(t, (*inbound.GetEpisodePort)(nil), getEpisodeService)
}

func Test_should_throw_error_if_show_does_not_exist_on_get_episode(t *testing.T) {
	defer initAdapter()
	command := &inbound.GetEpisodeCommand{
		EpisodeId: "some-episode-id",
		ShowId:    "i-do-not-exist",
	}

	mockGetShowAdapter.returnsOnGetOrNilShow["i-do-not-exist"] = nil

	result, err := getEpisodeService.GetEpisode(command)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, &error2.ShowNotFoundError{Id: "i-do-not-exist"}, err)
	assert.Equal(t, 0, mockSaveAndGetEpisodeAdapter.calledGet)
}

func Test_should_propagate_errors_from_adapter_on_get(t *testing.T) {
	defer initAdapter()

	expectedError := errors.New("some error")
	expectedShow := &model.Show{Id: "mocked-show-id"}
	mockGetShowAdapter.returnsOnGetOrNilShow["some-show-id"] = expectedShow
	mockSaveAndGetEpisodeAdapter.withErrorOnGetEpisodeOrNil = expectedError

	foundEpisode, err := getEpisodeService.GetEpisode(&inbound.GetEpisodeCommand{EpisodeId: "id-with-error", ShowId: "some-show-id"})

	assert.Nil(t, foundEpisode)
	assert.NotNil(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, 1, mockSaveAndGetEpisodeAdapter.calledGet)
}

func Test_should_return_not_found_if_episode_was_not_found_on_get(t *testing.T) {
	defer initAdapter()

	expectedShow := &model.Show{Id: "mocked-show-id"}
	mockGetShowAdapter.returnsOnGetOrNilShow["some-show-id"] = expectedShow

	foundShow, err := getEpisodeService.GetEpisode(&inbound.GetEpisodeCommand{EpisodeId: "id-with-error", ShowId: "some-show-id"})

	assert.Nil(t, foundShow)
	assert.NotNil(t, err)
	assert.Equal(t, &error2.EpisodeNotFoundError{Id: "id-with-error"}, err)
	assert.Equal(t, 1, mockSaveAndGetEpisodeAdapter.calledGet)
}

func Test_retrieve_episode_from_repository_on_get(t *testing.T) {
	defer initAdapter()

	expectedEpisode := &model.Episode{
		Id:     "some-id",
		ShowId: "some-show-id",
		Title:  "some title",
	}
	expectedEpisodeResponse := &inbound.GetEpisodeResponse{
		Id:     "some-id",
		ShowId: "some-show-id",
		Title:  "some title",
	}
	expectedShow := &model.Show{Id: "mocked-show-id"}
	mockGetShowAdapter.returnsOnGetOrNilShow["some-show-id"] = expectedShow
	mockSaveAndGetEpisodeAdapter.returnsOnGetEpisodeOrNil["some-id"] = expectedEpisode

	foundEpisode, err := getEpisodeService.GetEpisode(&inbound.GetEpisodeCommand{EpisodeId: "some-id", ShowId: "some-show-id"})

	assert.Nil(t, err)
	assert.NotNil(t, foundEpisode)
	assert.Equal(t, expectedEpisodeResponse, foundEpisode)
	assert.Equal(t, 1, mockSaveAndGetEpisodeAdapter.calledGet)
}
