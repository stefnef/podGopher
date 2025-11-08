package service

import (
	"errors"
	error2 "podGopher/core/domain/error"
	"podGopher/core/domain/model"
	"podGopher/core/port/inbound"
	"testing"

	"github.com/stretchr/testify/assert"
)

type getShowTestAdapter struct {
	called                  int
	returnsOnGetOrNilShow   map[string]*model.Show
	withErrorOnGetOrNilShow error
}

func (a *getShowTestAdapter) GetShowOrNil(Id string) (*model.Show, error) {
	a.called++
	show := a.returnsOnGetOrNilShow[Id]
	return show, a.withErrorOnGetOrNilShow
}

var mockGetShowAdapter = newGetShowTestAdapter()
var getShowService = NewGetShowService(mockGetShowAdapter)

func newGetShowTestAdapter() *getShowTestAdapter {
	adapter := &getShowTestAdapter{}
	adapter.init()
	return adapter
}
func (a *getShowTestAdapter) init() {
	a.called = 0
	a.returnsOnGetOrNilShow = make(map[string]*model.Show)
	a.withErrorOnGetOrNilShow = nil
}

func Test_should_implement_GetShowInPort(t *testing.T) {
	assert.NotNil(t, getShowService)
	assert.Implements(t, (*inbound.GetShowPort)(nil), getShowService)
}

func Test_should_return_not_found_if_show_was_not_found(t *testing.T) {
	defer mockGetShowAdapter.init()

	mockGetShowAdapter.returnsOnGetOrNilShow["non-existing-show-id"] = nil
	show := mockGetShowAdapter.returnsOnGetOrNilShow["non-existing-show-id"]
	assert.Nil(t, show)

	command := &inbound.GetShowCommand{Id: "non-existing-show-id"}
	result, err := getShowService.GetShow(command)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, &error2.ShowNotFoundError{Id: "non-existing-show-id"}, err)
	assert.Equal(t, 1, mockGetShowAdapter.called)
}

func Test_should_propagate_errors_from_adapter_on_get(t *testing.T) {
	defer mockGetShowAdapter.init()

	expectedError := errors.New("some error")
	mockGetShowAdapter.withErrorOnGetOrNilShow = expectedError

	foundShow, err := getShowService.GetShow(&inbound.GetShowCommand{Id: "id-with-error"})

	assert.Nil(t, foundShow)
	assert.NotNil(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, 1, mockGetShowAdapter.called)
}

func Test_retrieve_show_from_repository_on_get(t *testing.T) {
	defer mockGetShowAdapter.init()

	expectedShow := &model.Show{
		Id:    "some-id",
		Title: "some title",
		Slug:  "some-slug",
	}
	expectedShowResponse := &inbound.GetShowResponse{
		Id:    "some-id",
		Title: "some title",
		Slug:  "some-slug",
	}
	mockGetShowAdapter.withErrorOnGetOrNilShow = nil
	mockGetShowAdapter.returnsOnGetOrNilShow["some-id"] = expectedShow

	foundShow, err := getShowService.GetShow(&inbound.GetShowCommand{Id: "some-id"})

	assert.Nil(t, err)
	assert.NotNil(t, foundShow)
	assert.Equal(t, expectedShowResponse, foundShow)
	assert.Equal(t, 1, mockGetShowAdapter.called)
}
