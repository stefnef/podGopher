package show

// TODO analog zu get_show_handler_test.go: Hier einen get_all_shows_handler erstellen und testen
//var getShowHandler = NewGetAllShowsHandler(inbound.PortMap{
//	inbound.GetShow: mockGetShowService,
//})
//
//func Test_should_implement_handler_for_get_all_show(t *testing.T) {
//	assert.NotNil(t, getShowHandler)
//	assert.Implements(t, (*handler.Handler)(nil), getShowHandler)
//}
//
//func Test_should_panic_if_no_port_was_found_on_get_all_shows_handler(t *testing.T) {
//	invalidPortMap := inbound.PortMap{
//		inbound.PortInvalid: mockCreateShowService,
//	}
//
//	assert.Panics(t, func() {
//		NewGetShowHandler(invalidPortMap)
//	})
//}
//
//func Test_should_return_route_on_get_all_shows(t *testing.T) {
//	var route = getShowHandler.GetRoute()
//
//	var expectedRoute = &handler.Route{
//		Method: "GET",
//		Path:   "/show",
//	}
//
//	assert.Equal(t, expectedRoute, route)
//}
//
//func Test_should_propagate_error_on_get_all_shows(t *testing.T) {
//	defer mockGetShowService.init()
//	var context, _ = handlerTestSetup.GetTestGinContext(t)
//	expectedError := errors.New("some error")
//
//	mockGetShowService.failsWith = expectedError
//
//	context.Request = httptest.NewRequest("GET", "/show", bytes.NewBuffer([]byte("")))
//
//	getShowHandler.Handle(context)
//
//	assert.NotEmpty(t, context.Errors)
//	assert.Equal(t, expectedError, (*context.Errors[0]).Err)
//}
//
//func Test_should_call_service_on_get_all_shows(t *testing.T) {
//	defer mockGetShowService.init()
//	var getShowDto *allShowsResponseDto
//
//	type testParameterStruct struct {
//		title               string
//		mockedPortResponse  *show.GetAllShowsResponse
//		expectedWebResponse *allShowsResponseDto
//	}
//
//	tests := []testParameterStruct{
//		//TODO Beispiel mit nil und empty list füllen
//		{
//			"one show returned",
//			&show.GetAllShowsResponse{
//				Shows: []*model.Show{
//					{
//						Id:       "some-id",
//						Title:    "Mocked Title",
//						Slug:     "Mocked Slug",
//						Episodes: []string{},
//					},
//				},
//			},
//			&allShowsResponseDto{
//				//TODO Erwartung füllen
//			},
//		},
//	}
//
//	for _, tc := range tests {
//		var context, recorder = handlerTestSetup.GetTestGinContext(t)
//		mockGetShowService.init()
//
//		t.Run(tc.title, func(t *testing.T) {
//			mockGetShowService.returnsOnGetAllShow = tc.mockedPortResponse
//
//			context.Request = httptest.NewRequest("GET", "/show", bytes.NewBuffer([]byte("")))
//
//			getShowHandler.Handle(context)
//
//			var err = json.Unmarshal(recorder.Body.Bytes(), &getShowDto)
//
//			assert.Equal(t, 1, mockGetShowService.called)
//			assert.Nil(t, err)
//			assert.Empty(t, context.Errors)
//			assert.Equal(t, tc.expectedWebResponse, getShowDto)
//			assert.Equal(t, http.StatusOK, recorder.Code)
//		})
//	}
//}
