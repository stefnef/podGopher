package show

import (
	"errors"
	domainError "podGopher/core/domain/error"
	"podGopher/core/domain/model"
	onGetShow "podGopher/core/port/inbound/show"
	"testing"

	"github.com/stretchr/testify/assert"
)

var getShowService = NewGetShowService(mockGetShowAdapter)

func Test_should_implement_GetShowInPort(t *testing.T) {
	assert.NotNil(t, getShowService)
	assert.Implements(t, (*onGetShow.GetShowPort)(nil), getShowService)
}

func Test_should_return_not_found_if_show_was_not_found(t *testing.T) {
	defer initAdapter()

	mockGetShowAdapter.returnsOnGetOrNilShow["non-existing-show-id"] = nil
	show := mockGetShowAdapter.returnsOnGetOrNilShow["non-existing-show-id"]
	assert.Nil(t, show)

	command := &onGetShow.GetShowCommand{Id: "non-existing-show-id"}
	result, err := getShowService.GetShow(command)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, domainError.NewShowNotFoundError("non-existing-show-id"), err)
	assert.Equal(t, 1, mockGetShowAdapter.called)
}

func Test_should_propagate_errors_from_adapter_on_get(t *testing.T) {
	defer initAdapter()

	expectedError := errors.New("some error")
	mockGetShowAdapter.withErrorOnGetOrNilShow = expectedError

	foundShow, err := getShowService.GetShow(&onGetShow.GetShowCommand{Id: "id-with-error"})

	assert.Nil(t, foundShow)
	assert.NotNil(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, 1, mockGetShowAdapter.called)
}

func Test_retrieve_show_from_repository_on_get(t *testing.T) {
	defer initAdapter()

	expectedShow := &model.Show{
		Id:       "some-id",
		Title:    "some title",
		Slug:     "some-slug",
		Episodes: []string{"some-episode-id"},
	}
	expectedShowResponse := &onGetShow.GetShowResponse{
		Id:       "some-id",
		Title:    "some title",
		Slug:     "some-slug",
		Episodes: []string{"some-episode-id"},
	}
	mockGetShowAdapter.withErrorOnGetOrNilShow = nil
	mockGetShowAdapter.returnsOnGetOrNilShow["some-id"] = expectedShow

	foundShow, err := getShowService.GetShow(&onGetShow.GetShowCommand{Id: "some-id"})

	assert.Nil(t, err)
	assert.NotNil(t, foundShow)
	assert.Equal(t, expectedShowResponse, foundShow)
	assert.Equal(t, 1, mockGetShowAdapter.called)
}
