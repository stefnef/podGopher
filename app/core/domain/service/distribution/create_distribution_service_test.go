package distribution

import (
	"errors"
	domainError "podGopher/core/domain/error"
	"podGopher/core/domain/model"
	onCreateDistribution "podGopher/core/port/inbound/distribution"
	"testing"

	"github.com/stretchr/testify/assert"
)

var createDistributionService = NewCreateDistributionService(mockGetShowAdapter, mockSaveAndGetDistributionAdapter)

func Test_should_implement_CreateDistributionInPort(t *testing.T) {
	assert.NotNil(t, createDistributionService)
	assert.Implements(t, (*onCreateDistribution.CreateDistributionPort)(nil), createDistributionService)
}

func Test_should_throw_error_if_distribution_with_name_already_exists(t *testing.T) {
	defer initAdapter()

	mockSaveAndGetDistributionAdapter.everyExistsByTitleReturns("Test", "Slug", true)

	command := newTestCreateDistributionCommand("Test")
	result, err := createDistributionService.CreateDistribution(command)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, domainError.NewDistributionAlreadyExistsError("Test"), err)
	assert.Equal(t, 0, mockSaveAndGetDistributionAdapter.calledSave)
}

func Test_should_propagate_errors_from_adapter_on_create_distribution(t *testing.T) {
	defer initAdapter()

	expectedError := errors.New("some error")
	mockSaveAndGetDistributionAdapter.withErrorOnSaveDistribution = expectedError
	mockGetShowAdapter.returnsOnGetOrNilShow["test-show-id"] = &model.Show{Id: "test-show-id"}

	command := newTestCreateDistributionCommand("Fake")
	result, err := createDistributionService.CreateDistribution(command)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, 1, mockSaveAndGetDistributionAdapter.calledSave)
}

func Test_should_throw_error_if_show_does_not_exist_on_save_distribution(t *testing.T) {
	defer initAdapter()
	command := newTestCreateDistributionCommand("Test")

	mockSaveAndGetDistributionAdapter.everyExistsByTitleReturns("Test", "Slug", false)
	mockGetShowAdapter.returnsOnGetOrNilShow["test-show-id"] = nil

	result, err := createDistributionService.CreateDistribution(command)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, domainError.NewShowNotFoundError("test-show-id"), err)
	assert.Equal(t, 0, mockSaveAndGetDistributionAdapter.calledSave)
}

func Test_should_save_a_new_distribution(t *testing.T) {
	defer initAdapter()

	mockSaveAndGetDistributionAdapter.everyExistsByTitleReturns("Test", "Slug", false)
	mockGetShowAdapter.returnsOnGetOrNilShow["test-show-id"] = &model.Show{Id: "mocked-show-id"}
	createDistributionCommand := newTestCreateDistributionCommand("Test")

	result, err := createDistributionService.CreateDistribution(createDistributionCommand)

	savedDistribution := mockSaveAndGetDistributionAdapter.onSaveCalledWith

	expectedSavedDistribution := &model.Distribution{
		Id:     savedDistribution.Id,
		ShowId: "test-show-id",
		Title:  "Test",
		Slug:   "Slug",
	}
	assert.NotNil(t, savedDistribution)
	assert.Equal(t, 1, mockSaveAndGetDistributionAdapter.calledSave)
	assert.Equal(t, expectedSavedDistribution, savedDistribution)
	assert.NotEmpty(t, savedDistribution.Id)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.IsType(t, (*onCreateDistribution.CreateDistributionResponse)(nil), result)

	expectedCreatedDistribution := &onCreateDistribution.CreateDistributionResponse{Id: savedDistribution.Id, ShowId: "test-show-id", Title: "Test", Slug: "Slug"}
	assert.Equal(t, expectedCreatedDistribution, result)
}
