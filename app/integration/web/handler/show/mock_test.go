package show

import "podGopher/core/port/inbound/show"

type getShowTestService struct {
	called              int
	command             *show.GetShowCommand
	returnsOnGetShow    *show.GetShowResponse
	returnsOnGetAllShow *show.GetAllShowsResponse
	failsWith           error
}

func (s *getShowTestService) init() {
	s.called = 0
	s.command = nil
	s.returnsOnGetShow = nil
	s.returnsOnGetAllShow = nil
	s.failsWith = nil
}

func (s *getShowTestService) GetShow(command *show.GetShowCommand) (show *show.GetShowResponse, err error) {
	s.called++
	s.command = command
	return s.returnsOnGetShow, s.failsWith
}

func (s *getShowTestService) GetAllShows() (shows *show.GetAllShowsResponse, err error) {
	s.called++
	return s.returnsOnGetAllShow, s.failsWith
}

var mockGetShowService = new(getShowTestService)
