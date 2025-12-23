package show

import (
	"errors"
	error2 "podGopher/core/domain/error"
	"podGopher/core/domain/model"
	"podGopher/core/port/inbound"
	"testing"

	"github.com/stretchr/testify/assert"
)

var createShowService = NewCreateShowService(mockSaveAndGetShowAdapter)

func Test_should_implement_CreateShowInPort(t *testing.T) {
	assert.NotNil(t, createShowService)
	assert.Implements(t, (*inbound.CreateShowPort)(nil), createShowService)
}

func Test_should_save_a_new_show(t *testing.T) {
	defer initAdapter()

	mockSaveAndGetShowAdapter.everyExistsByTitleOrSlugReturns("Test", "Test-Slug", false)
	createShowCommand := newTestCreateShowCommand("Test")

	result, err := createShowService.CreateShow(createShowCommand)

	savedShow := mockSaveAndGetShowAdapter.onSave["show"]

	expectedSavedShow := &model.Show{
		Id:    savedShow.Id,
		Title: "Test",
		Slug:  "Test-Slug",
	}
	assert.NotNil(t, savedShow)
	assert.Equal(t, 1, mockSaveAndGetShowAdapter.calledSave)
	assert.Equal(t, expectedSavedShow, savedShow)
	assert.NotEmpty(t, savedShow.Id)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.IsType(t, (*inbound.CreateShowResponse)(nil), result)

	expectedCreatedShow := &inbound.CreateShowResponse{Id: savedShow.Id, Title: "Test", Slug: "Test-Slug"}
	assert.Equal(t, expectedCreatedShow, result)
}

func Test_should_throw_error_if_show_with_name_already_exists(t *testing.T) {
	defer initAdapter()

	mockSaveAndGetShowAdapter.everyExistsByTitleOrSlugReturns("Test", "Test-Slug", true)

	show := newTestCreateShowCommand("Test")
	result, err := createShowService.CreateShow(show)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, error2.NewShowAlreadyExistsError("Test"), err)
	assert.Equal(t, 0, mockSaveAndGetShowAdapter.calledSave)
}

func Test_should_propagate_errors_from_adapter_on_create_show(t *testing.T) {
	defer initAdapter()

	expectedError := errors.New("some error")
	mockSaveAndGetShowAdapter.withErrorOnSaveShow = expectedError

	show := newTestCreateShowCommand("Fake")
	result, err := createShowService.CreateShow(show)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, 1, mockSaveAndGetShowAdapter.calledSave)
}
