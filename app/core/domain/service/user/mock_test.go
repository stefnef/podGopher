package user

import (
	"podGopher/adapter/outbound/repository/mock"
	"testing"
)

var mockGetShowAdapter = mock.NewGetShowTestAdapter(nil)
var mockGetAndSaveUserAdapter = mock.NewGetAndSaveUserTestAdapter(nil)

func initAdapter(t *testing.T) {
	mockGetShowAdapter.Init(t)
	mockGetAndSaveUserAdapter.Init(t)
}
