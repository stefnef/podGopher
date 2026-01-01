package distribution

import (
	"errors"
	domainError "podGopher/core/domain/error"
	"podGopher/core/domain/model"
	onGetDistribution "podGopher/core/port/inbound/distribution"
	"testing"

	"github.com/stretchr/testify/assert"
)

var getDistributionService = NewGetDistributionService(mockGetShowAdapter, mockSaveAndGetDistributionAdapter)

func Test_should_implement_GetDistributionInPort(t *testing.T) {
	assert.NotNil(t, getDistributionService)
	assert.Implements(t, (*onGetDistribution.GetDistributionPort)(nil), getDistributionService)
}

func Test_should_throw_error_if_show_does_not_exist_on_get_distribution(t *testing.T) {
	defer initAdapter()
	command := &onGetDistribution.GetDistributionCommand{
		DistributionId: "some-distribution-id",
		ShowId:         "i-do-not-exist",
	}

	mockGetShowAdapter.returnsOnGetOrNilShow["i-do-not-exist"] = nil

	result, err := getDistributionService.GetDistribution(command)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, domainError.NewShowNotFoundError("i-do-not-exist"), err)
	assert.Equal(t, 0, mockSaveAndGetDistributionAdapter.calledGet)
}

func Test_should_propagate_errors_from_adapter_on_get(t *testing.T) {
	defer initAdapter()

	expectedError := errors.New("some error")
	expectedShow := &model.Show{Id: "mocked-show-id"}
	mockGetShowAdapter.returnsOnGetOrNilShow["some-show-id"] = expectedShow
	mockSaveAndGetDistributionAdapter.withErrorOnGetDistributionOrNil = expectedError

	foundDistribution, err := getDistributionService.GetDistribution(&onGetDistribution.GetDistributionCommand{DistributionId: "id-with-error", ShowId: "some-show-id"})

	assert.Nil(t, foundDistribution)
	assert.NotNil(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, 1, mockSaveAndGetDistributionAdapter.calledGet)
}

func Test_should_return_not_found_if_distribution_was_not_found_on_get(t *testing.T) {
	defer initAdapter()

	expectedShow := &model.Show{Id: "mocked-show-id"}
	mockGetShowAdapter.returnsOnGetOrNilShow["some-show-id"] = expectedShow

	foundShow, err := getDistributionService.GetDistribution(&onGetDistribution.GetDistributionCommand{DistributionId: "id-with-error", ShowId: "some-show-id"})

	assert.Nil(t, foundShow)
	assert.NotNil(t, err)
	assert.Equal(t, domainError.NewDistributionNotFoundError("id-with-error"), err)
	assert.Equal(t, 1, mockSaveAndGetDistributionAdapter.calledGet)
}

func Test_retrieve_distribution_from_repository_on_get(t *testing.T) {
	defer initAdapter()

	expectedDistribution := &model.Distribution{
		Id:     "some-id",
		ShowId: "some-show-id",
		Title:  "some title",
		Slug:   "some-slug",
	}
	expectedDistributionResponse := &onGetDistribution.GetDistributionResponse{
		Id:     "some-id",
		ShowId: "some-show-id",
		Title:  "some title",
		Slug:   "some-slug",
	}
	expectedShow := &model.Show{Id: "mocked-show-id"}
	mockGetShowAdapter.returnsOnGetOrNilShow["some-show-id"] = expectedShow
	mockSaveAndGetDistributionAdapter.returnsOnGetDistributionOrNil["some-id"] = expectedDistribution

	foundDistribution, err := getDistributionService.GetDistribution(&onGetDistribution.GetDistributionCommand{DistributionId: "some-id", ShowId: "some-show-id"})

	assert.Nil(t, err)
	assert.NotNil(t, foundDistribution)
	assert.Equal(t, expectedDistributionResponse, foundDistribution)
	assert.Equal(t, 1, mockSaveAndGetDistributionAdapter.calledGet)
}
