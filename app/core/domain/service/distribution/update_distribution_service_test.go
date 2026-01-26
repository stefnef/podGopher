package distribution

import (
	onUpdateDistribution "podGopher/core/port/inbound/distribution"
	"testing"

	"github.com/stretchr/testify/assert"
)

var updateDistributionService = NewUpdateDistributionService(mockGetShowAdapter, mockSaveAndGetDistributionAdapter)

func Test_should_implement_UpdateDistributionInPort(t *testing.T) {
	assert.NotNil(t, createDistributionService)
	assert.Implements(t, (*onUpdateDistribution.UpdateDistributionPort)(nil), updateDistributionService)
}

// TODO hier weiter
