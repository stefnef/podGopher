package rss

import (
	"errors"
	"podGopher/core/domain/model"
	onGetRSS "podGopher/core/port/inbound/rss"
	"testing"

	"github.com/stretchr/testify/assert"
)

var getRSSService = NewGetRSSService(mockGetRSSAdapter)

func Test_should_implement_GetRSSInPort(t *testing.T) {
	assert.NotNil(t, getRSSService)
	assert.Implements(t, (*onGetRSS.GetRSSPort)(nil), getRSSService)
}

func Test_should_propagate_errors_from_adapter_on_get(t *testing.T) {
	defer initAdapter()

	expectedError := errors.New("some error")
	mockGetRSSAdapter.returnsOnGetOrNilError["showSlugRSSSlug"] = expectedError

	foundRss, err := getRSSService.GetRSS(&onGetRSS.GetRSSCommand{ShowSlug: "showSlug", DistributionSlug: "RSSSlug"})

	assert.Nil(t, foundRss)
	assert.NotNil(t, err)
	assert.Equal(t, expectedError, err)
	assert.Equal(t, 1, mockGetRSSAdapter.called)
}

func Test_retrieve_rss_from_repository_on_get(t *testing.T) {
	defer initAdapter()

	expectedRSS := &model.RSS{
		Show: &model.Show{
			Id:            "some-show-id",
			Title:         "show-title",
			Slug:          "show-slug",
			Episodes:      []string{"some-episode-id", "some-other-episode-id"},
			Distributions: []string{"some-distribution-id", "some-other-distribution-id"},
		},
		Distribution: &model.Distribution{
			Id:       "some-distribution-id",
			ShowId:   "some-show-id",
			Title:    "distribution-title",
			Slug:     "distribution-slug",
			Episodes: []string{"episode1-id"},
		},
		Episodes: []*model.Episode{{
			Id:     "some-episode-id",
			ShowId: "some-show-id",
			Title:  "episode-title",
		}},
	}
	var expectedRSSResponse = &onGetRSS.GetRSSResponse{
		ShowTitle:         "show-title",
		DistributionTitle: "distribution-title",
		Episodes: []*onGetRSS.GetRSSEpisodeResponse{{
			Title: "episode-title",
		}},
	}
	mockGetRSSAdapter.returnsOnGetOrNilRSS["show-slugdistribution-slug"] = expectedRSS

	foundRSS, err := getRSSService.GetRSS(&onGetRSS.GetRSSCommand{ShowSlug: "show-slug", DistributionSlug: "distribution-slug"})

	assert.Nil(t, err)
	assert.NotNil(t, foundRSS)
	assert.Equal(t, expectedRSSResponse, foundRSS)
	assert.Equal(t, 1, mockGetRSSAdapter.called)
}
