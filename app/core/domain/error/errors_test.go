package error

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_is_an_error(t *testing.T) {
	tests := map[string]struct {
		err                 error
		expectedErrorString string
		category            Category
	}{
		"ShowAlreadyExistsError": {
			NewShowAlreadyExistsError("some-name"),
			"show with title 'some-name' or given slug already exists",
			AlreadyExists,
		},

		"ShowNotFoundError": {
			NewShowNotFoundError("some-id"),
			"show with id 'some-id' does not exist",
			NotFound,
		},

		"EpisodeAlreadyExistsError": {
			NewEpisodeAlreadyExistsError("some-name"),
			"episode with title 'some-name' already exists",
			AlreadyExists,
		},

		"EpisodeNotFoundError": {
			NewEpisodeNotFoundError("some-id"),
			"episode with id 'some-id' does not exist",
			NotFound,
		},

		"DistributionAlreadyExistsError": {
			NewDistributionAlreadyExistsError("some-name"),
			"distribution with title 'some-name' or given slug already exists",
			AlreadyExists,
		},

		"DistributionNotFoundError": {
			NewDistributionNotFoundError("some-id"),
			"distribution with id 'some-id' does not exist",
			NotFound,
		},

		"UpdateError": {
			NewUpdateError("some test"),
			"update not possible, reason: some test",
			DataConflict,
		},

		"RSSFeedNotFoundError": {
			NewRSSFeedNotFoundError("RSS"),
			"rss feed 'RSS' does not exist",
			NotFound,
		},

		"UserAlreadyExistsError": {
			NewUserAlreadyExistsError("some-show", "any-user"),
			"user with username 'any-user' already exists for show 'some-show'",
			AlreadyExists,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Implements(t, (*error)(nil), test.err)
			assert.ErrorContains(t, test.err.(error), test.expectedErrorString)
			assert.Equal(t, test.category, test.err.(*ModelError).Category)
		})
	}

}
