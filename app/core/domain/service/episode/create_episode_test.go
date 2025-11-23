package episode

import (
	"errors"
	error2 "podGopher/core/domain/error"
	"podGopher/core/domain/model"
	"podGopher/core/port/inbound"
	"testing"

	"github.com/stretchr/testify/assert"
)

var createEpisodeService = NewCreateEpisodeService(mockGetShowAdapter, mockSaveAndGetEpisodeAdapter)

func Test_should_implement_CreateEpisodeInPort(t *testing.T) {
	assert.NotNil(t, createEpisodeService)
	assert.Implements(t, (*inbound.CreateEpisodePort)(nil), createEpisodeService)
}

func Test_should_throw_error_if_episode_with_name_already_exists(t *testing.T) {
	defer initAdapter()

	mockSaveAndGetEpisodeAdapter.everyExistsByTitleReturns("Test", true)

	command := newTestCreateEpisodeCommand("Test")
	result, err := createEpisodeService.CreateEpisode(command)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, &error2.EpisodeAlreadyExistsError{Name: "Test"}, err)
	assert.Equal(t, 0, mockSaveAndGetEpisodeAdapter.calledSave)
}

func Test_should_propagate_errors_from_adapter_on_create_episode(t *testing.T) {
	defer initAdapter()

	expectedError := errors.New("some error")
	mockSaveAndGetEpisodeAdapter.withErrorOnSaveEpisode = expectedError
	mockGetShowAdapter.returnsOnGetOrNilShow["test-show-id"] = &model.Show{Id: "test-show-id"}

	command := newTestCreateEpisodeCommand("Fake")
	result, err := createEpisodeService.CreateEpisode(command)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, 1, mockSaveAndGetEpisodeAdapter.calledSave)
}

func Test_should_throw_error_if_show_does_not_exist_on_save_episode(t *testing.T) {
	defer initAdapter()
	command := newTestCreateEpisodeCommand("Test")

	mockSaveAndGetEpisodeAdapter.everyExistsByTitleReturns("Test", false)
	mockGetShowAdapter.returnsOnGetOrNilShow["test-show-id"] = nil

	result, err := createEpisodeService.CreateEpisode(command)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, &error2.ShowNotFoundError{Id: "test-show-id"}, err)
	assert.Equal(t, 0, mockSaveAndGetEpisodeAdapter.calledSave)
}

func Test_should_save_a_new_episode(t *testing.T) {
	defer initAdapter()

	mockSaveAndGetEpisodeAdapter.everyExistsByTitleReturns("Test", false)
	mockGetShowAdapter.returnsOnGetOrNilShow["test-show-id"] = &model.Show{Id: "mocked-show-id"}
	createEpisodeCommand := newTestCreateEpisodeCommand("Test")

	result, err := createEpisodeService.CreateEpisode(createEpisodeCommand)

	savedEpisode := mockSaveAndGetEpisodeAdapter.onSaveCalledWith

	expectedSavedEpisode := &model.Episode{
		Id:     savedEpisode.Id,
		ShowId: "test-show-id",
		Title:  "Test",
	}
	assert.NotNil(t, savedEpisode)
	assert.Equal(t, 1, mockSaveAndGetEpisodeAdapter.calledSave)
	assert.Equal(t, expectedSavedEpisode, savedEpisode)
	assert.NotEmpty(t, savedEpisode.Id)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.IsType(t, (*inbound.CreateEpisodeResponse)(nil), result)

	expectedCreatedEpisode := &inbound.CreateEpisodeResponse{Id: savedEpisode.Id, ShowId: "test-show-id", Title: "Test"}
	assert.Equal(t, expectedCreatedEpisode, result)
}
