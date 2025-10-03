package service

import (
	"errors"
	"podGopher/core/port/inbound"
	"testing"

	"github.com/stretchr/testify/assert"
)

type saveAndGetShowTestAdapter struct {
	calledSave           int
	onSave               map[string]interface{}
	returnsExistsByTitle map[string]bool
	returnsSaveShow      error
}

func newSaveAndGetShowTestAdapter() *saveAndGetShowTestAdapter {
	adapter := &saveAndGetShowTestAdapter{}
	adapter.init()
	return adapter
}

func newTestCreateShowCommand(title string) *inbound.CreateShowCommand {
	show := &inbound.CreateShowCommand{
		Title: title,
	}
	return show
}

func (adapter *saveAndGetShowTestAdapter) SaveShow(title string) error {
	(*adapter).calledSave++
	(*adapter).onSave["title"] = title
	return (*adapter).returnsSaveShow
}

func (adapter *saveAndGetShowTestAdapter) init() {
	(*adapter).calledSave = 0
	(*adapter).onSave = make(map[string]interface{})
	(*adapter).returnsExistsByTitle = make(map[string]bool)
	(*adapter).returnsSaveShow = nil
}

func (adapter *saveAndGetShowTestAdapter) everyExistsByTitleReturns(title string, returnValue bool) {
	(*adapter).returnsExistsByTitle[title] = returnValue
}

func (adapter *saveAndGetShowTestAdapter) ExistsByTitle(title string) bool {
	return (*adapter).returnsExistsByTitle[title]
}

var mockSaveAndGetShowAdapter = newSaveAndGetShowTestAdapter()
var createShowService = NewCreateShowService(mockSaveAndGetShowAdapter)

func Test_should_implement_CreateShowInPort(t *testing.T) {
	assert.NotNil(t, createShowService)
	assert.Implements(t, (*inbound.CreateShowPort)(nil), createShowService)
}

func Test_should_save_a_new_show(t *testing.T) {
	defer mockSaveAndGetShowAdapter.init()

	mockSaveAndGetShowAdapter.everyExistsByTitleReturns("Test", false)

	show := newTestCreateShowCommand("Test")
	err := createShowService.CreateShow(show)

	assert.Equal(t, 1, mockSaveAndGetShowAdapter.calledSave)
	assert.Equal(t, "Test", mockSaveAndGetShowAdapter.onSave["title"])
	assert.Nil(t, err)
}

func Test_should_throw_error_if_show_with_name_already_exists(t *testing.T) {
	defer mockSaveAndGetShowAdapter.init()

	mockSaveAndGetShowAdapter.everyExistsByTitleReturns("Test", true)

	show := newTestCreateShowCommand("Test")
	err := createShowService.CreateShow(show)

	assert.NotNil(t, err)
	assert.Equal(t, &inbound.ShowAlreadyExistsError{Name: "Test"}, err)
	assert.Equal(t, 0, mockSaveAndGetShowAdapter.calledSave)
}

func Test_should_propagate_errors_from_adapter(t *testing.T) {
	defer mockSaveAndGetShowAdapter.init()

	expectedError := errors.New("some error")
	mockSaveAndGetShowAdapter.returnsSaveShow = expectedError

	show := newTestCreateShowCommand("Fake")
	err := createShowService.CreateShow(show)

	assert.NotNil(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, 1, mockSaveAndGetShowAdapter.calledSave)
}
