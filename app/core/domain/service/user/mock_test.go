package user

import (
	mockOAuth "podGopher/adapter/outbound/credentials/mock"
	mockRepository "podGopher/adapter/outbound/repository/mock"
	"testing"
)

var mockGetShowAdapter = mockRepository.NewGetShowTestAdapter(nil)
var mockGetAndSaveUserAdapter = mockRepository.NewGetAndSaveUserTestAdapter(nil)
var mockCreateUserCredentialsAdapter = mockOAuth.NewCreateUserTestAdapter(nil)

func initAdapter(t *testing.T) {
	mockGetShowAdapter.Init(t)
	mockGetAndSaveUserAdapter.Init(t)
	mockCreateUserCredentialsAdapter.Init(t)
}
