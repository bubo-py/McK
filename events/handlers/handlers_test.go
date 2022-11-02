package handlers

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bubo-py/McK/customErrors"
	"github.com/bubo-py/McK/events/repositories/mocks"
	"github.com/bubo-py/McK/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetEventsHandler(t *testing.T) {
	testCases := []struct {
		testName          string
		strConvErrPresent bool
		r                 *http.Request
		w                 *httptest.ResponseRecorder
		mockDataReturn    []types.Event
		mockErrReturn     error
		expFilters        types.Filters
		expJSONReturn     string
		expStatusCode     int
	}{
		{
			testName: "GetEvents_with_two_events_return",
			r:        httptest.NewRequest("GET", "/api/events", nil),
			w:        httptest.NewRecorder(),
			mockDataReturn: []types.Event{
				{
					ID:        1,
					Name:      "Daily meeting",
					StartTime: time.Date(2020, 5, 15, 20, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 5, 15, 20, 30, 0, 0, time.UTC),
				},
				{
					ID:        2,
					Name:      "Weekly meeting",
					StartTime: time.Date(2020, 5, 15, 20, 30, 0, 0, time.UTC),
					EndTime:   time.Date(2020, 5, 15, 20, 30, 0, 0, time.UTC),
				},
			},
			mockErrReturn: nil,
			expJSONReturn: `[{"id":1,"name":"Daily meeting","startTime":"2020-05-15T20:30:00Z","endTime":"2020-05-15T20:30:00Z","alertTime":"0001-01-01T00:00:00Z"},{"id":2,"name":"Weekly meeting","startTime":"2020-05-15T20:30:00Z","endTime":"2020-05-15T20:30:00Z","alertTime":"0001-01-01T00:00:00Z"}]`,
			expStatusCode: 200,
		},
		{
			testName:      "GetEvents_with_two_filters",
			r:             httptest.NewRequest("GET", "/api/events?day=15&month=10", nil),
			w:             httptest.NewRecorder(),
			mockErrReturn: nil,
			expFilters:    types.Filters{Day: 15, Month: 10, Year: 0},
			expJSONReturn: "null\n",
			expStatusCode: 200,
		},
		{
			testName:      "GetEvents_with_one_filter",
			r:             httptest.NewRequest("GET", "/api/events?year=2222", nil),
			w:             httptest.NewRecorder(),
			mockErrReturn: nil,
			expFilters:    types.Filters{Day: 0, Month: 0, Year: 2222},
			expJSONReturn: "null\n",
			expStatusCode: 200,
		},
		{
			testName:      "GetEvents_BadRequest",
			r:             httptest.NewRequest("GET", "/api/events?day=100", nil),
			w:             httptest.NewRecorder(),
			mockErrReturn: customErrors.ErrBadRequest,
			expFilters:    types.Filters{Day: 100, Month: 0, Year: 0},
			expJSONReturn: `{"ErrorType":"BadRequest","ErrorMessage":"the server cannot process the request"}`,
			expStatusCode: 400,
		},
		{
			testName:          "GetEvents_DayStrConvErr_BadRequest",
			strConvErrPresent: true,
			r:                 httptest.NewRequest("GET", "/api/events?day=1.5", nil),
			w:                 httptest.NewRecorder(),
			expJSONReturn:     `{"ErrorType":"BadRequest","ErrorMessage":"the server cannot process the request"}`,
			expStatusCode:     400,
		},
		{
			testName:          "GetEvents_MonthStrConvErr_BadRequest",
			strConvErrPresent: true,
			r:                 httptest.NewRequest("GET", "/api/events?month=3.14", nil),
			w:                 httptest.NewRecorder(),
			expJSONReturn:     `{"ErrorType":"BadRequest","ErrorMessage":"the server cannot process the request"}`,
			expStatusCode:     400,
		},
		{
			testName:          "GetEvents_YearStrConvErr_BadRequest",
			strConvErrPresent: true,
			r:                 httptest.NewRequest("GET", "/api/events?year=20.22", nil),
			w:                 httptest.NewRecorder(),
			expJSONReturn:     `{"ErrorType":"BadRequest","ErrorMessage":"the server cannot process the request"}`,
			expStatusCode:     400,
		},
		{
			testName:          "GetEvents_StrConvErr_BadRequest",
			strConvErrPresent: true,
			r:                 httptest.NewRequest("GET", "/api/events?day=1.5&year=20.22&month=3.14", nil),
			w:                 httptest.NewRecorder(),
			expJSONReturn:     `{"ErrorType":"BadRequest","ErrorMessage":"the server cannot process the request"}`,
			expStatusCode:     400,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {

			// mock business logic
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockBL := mocks.NewMockBusinessLogicInterface(mockCtrl)

			if !tc.strConvErrPresent {
				mockBL.EXPECT().GetEvents(tc.r.Context(), tc.expFilters).Return(tc.mockDataReturn, tc.mockErrReturn).Times(1)
			}

			// create handler with mocks
			handler := InitHandler(mockBL)
			handler.GetEventsHandler(tc.w, tc.r)

			resp := tc.w.Result()

			data, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Error(err)
			}

			require.JSONEq(t, tc.expJSONReturn, string(data), "JSON data should be equal")

			require.Equal(t, tc.expStatusCode, resp.StatusCode, "Wrong status code returned")
		})
	}
}

func TestGetEventHandler(t *testing.T) {
	testCases := []struct {
		testName          string
		r                 *http.Request
		w                 *httptest.ResponseRecorder
		strConvErrPresent bool
		expID             int64
		eventToMock       types.Event
		mockErrReturn     error
		expJSONReturn     string
		expStatusCode     int
	}{
		{
			testName: "GetEvent_with_one_event_return",
			r:        httptest.NewRequest("GET", "/5", nil),
			w:        httptest.NewRecorder(),
			expID:    5,
			eventToMock: types.Event{
				ID:        1,
				Name:      "Onboarding Meeting",
				StartTime: time.Date(2022, 9, 14, 9, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2022, 9, 14, 9, 0, 0, 0, time.UTC),
			},
			expJSONReturn: `{"id":1,"name":"Onboarding Meeting","startTime":"2022-09-14T09:00:00Z","endTime":"2022-09-14T09:00:00Z", "alertTime":"0001-01-01T00:00:00Z"}`,
			expStatusCode: 200,
		},
		{
			testName:      "GetEvent_param_value_no_return",
			r:             httptest.NewRequest("GET", "/450", nil),
			w:             httptest.NewRecorder(),
			expID:         450,
			expStatusCode: 200,
		},
		{
			testName:      "GetEvent_param_value_Unexpected",
			r:             httptest.NewRequest("GET", "/450", nil),
			w:             httptest.NewRecorder(),
			expID:         450,
			mockErrReturn: customErrors.ErrUnexpected,
			expJSONReturn: `{"ErrorType":"Unexpected","ErrorMessage":"an unexpected error occurred"}`,
			expStatusCode: 500,
		},
		{
			testName:      "GetEvent_StrConvErr_NotFound",
			r:             httptest.NewRequest("GET", "/450", nil),
			w:             httptest.NewRecorder(),
			expID:         450,
			mockErrReturn: customErrors.ErrNotFound,
			expJSONReturn: `{"ErrorType":"NotFound","ErrorMessage":"the server cannot find the requested resource"}`,
			expStatusCode: 404,
		},
		{
			testName:          "GetEvent_StrConvErr_BadRequest",
			r:                 httptest.NewRequest("GET", "/4.50", nil),
			w:                 httptest.NewRecorder(),
			strConvErrPresent: true,
			mockErrReturn:     customErrors.ErrBadRequest,
			expJSONReturn:     `{"ErrorType":"BadRequest","ErrorMessage":"the server cannot process the request"}`,
			expStatusCode:     400,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {

			// mock business logic
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockBL := mocks.NewMockBusinessLogicInterface(mockCtrl)

			if !tc.strConvErrPresent {
				mockBL.EXPECT().GetEvent(gomock.Any(), tc.expID).Return(tc.eventToMock, tc.mockErrReturn)
			}

			// create handler with mocks
			handler := InitHandler(mockBL)
			handler.Mux.ServeHTTP(tc.w, tc.r)

			resp := tc.w.Result()

			data, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Error(err)
			}

			if tc.expJSONReturn != "" {
				require.JSONEq(t, tc.expJSONReturn, string(data), "JSON data should to be equal")

			}

			require.Equal(t, tc.expStatusCode, resp.StatusCode, "Wrong status code returned")

		})
	}
}

func TestAddEventHandler(t *testing.T) {
	testCases := []struct {
		testName         string
		decodeErrPresent bool
		jsonStr          string
		eventToMock      types.Event
		mockErrReturn    error
		expJSONReturn    string
		expStatusCode    int
	}{
		{
			testName: "AddEvent_positive_response_with_event_return",
			jsonStr:  `{"id":1,"name":"Onboarding Meeting","startTime":"2022-09-14T09:00:00Z","endTime":"2022-09-14T09:00:00Z", "alertTime":"0001-01-01T00:00:00Z"}`,
			eventToMock: types.Event{
				ID:        1,
				Name:      "Onboarding Meeting",
				StartTime: time.Date(2022, 9, 14, 9, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2022, 9, 14, 9, 0, 0, 0, time.UTC),
			},
			expJSONReturn: `{"id":1,"name":"Onboarding Meeting","startTime":"2022-09-14T09:00:00Z","endTime":"2022-09-14T09:00:00Z", "alertTime":"0001-01-01T00:00:00Z"}`,
			expStatusCode: 200,
		},
		{
			testName: "AddEvent_BadRequest",
			jsonStr:  `{"id":1,"name":"Supposedly too long Meeting Name","startTime":"2022-09-14T09:00:00Z","endTime":"2022-09-14T09:00:00Z", "alertTime":"0001-01-01T00:00:00Z"}`,
			eventToMock: types.Event{
				ID:        1,
				Name:      "Supposedly too long Meeting Name",
				StartTime: time.Date(2022, 9, 14, 9, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2022, 9, 14, 9, 0, 0, 0, time.UTC),
			},
			mockErrReturn: customErrors.ErrBadRequest,
			expJSONReturn: `{"ErrorType":"BadRequest","ErrorMessage":"the server cannot process the request"}`,
			expStatusCode: 400,
		},
		{
			testName: "AddEvent_Unexpected",
			jsonStr:  `{"id":1,"name":"Meeting Name","startTime":"2022-09-14T09:00:00Z","endTime":"2022-09-14T09:00:00Z", "alertTime":"0001-01-01T00:00:00Z"}`,
			eventToMock: types.Event{
				ID:        1,
				Name:      "Meeting Name",
				StartTime: time.Date(2022, 9, 14, 9, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2022, 9, 14, 9, 0, 0, 0, time.UTC),
			},
			mockErrReturn: customErrors.ErrUnexpected,
			expJSONReturn: `{"ErrorType":"Unexpected","ErrorMessage":"an unexpected error occurred"}`,
			expStatusCode: 500,
		},
		{
			testName:         "AddEvent_DecodeErr1",
			decodeErrPresent: true,
			jsonStr:          `{json string}`,
			expJSONReturn:    `{"ErrorType":"BadRequest","ErrorMessage":"the server cannot process the request"}`,
			expStatusCode:    400,
		},
		{
			testName:         "AddEvent_DecodeErr2",
			decodeErrPresent: true,
			jsonStr:          `{"hello": world}`,
			expJSONReturn:    `{"ErrorType":"BadRequest","ErrorMessage":"the server cannot process the request"}`,
			expStatusCode:    400,
		},
		{
			testName:         "AddEvent_DecodeErr3",
			decodeErrPresent: true,
			jsonStr:          `{"name": false}`,
			expJSONReturn:    `{"ErrorType":"BadRequest","ErrorMessage":"the server cannot process the request"}`,
			expStatusCode:    400,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {

			// mock request
			r := httptest.NewRequest("POST", "/api/events", bytes.NewBuffer([]byte(tc.jsonStr)))
			w := httptest.NewRecorder()

			// mock business logic
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockBL := mocks.NewMockBusinessLogicInterface(mockCtrl)

			if !tc.decodeErrPresent {
				mockBL.EXPECT().AddEvent(r.Context(), tc.eventToMock).Return(tc.mockErrReturn)
			}

			// create handler with mocks
			handler := InitHandler(mockBL)
			handler.AddEventHandler(w, r)

			resp := w.Result()

			data, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Error(err)
			}

			require.JSONEqf(t, tc.expJSONReturn, string(data), "JSON data should be equal")

			require.Equal(t, tc.expStatusCode, resp.StatusCode, "Wrong status code returned")

		})
	}
}

func TestDeleteEventHandler(t *testing.T) {
	testCases := []struct {
		testName      string
		r             *http.Request
		w             *httptest.ResponseRecorder
		expID         int64
		mockErrReturn error
		expJSONReturn string
		expStatusCode int
	}{
		{
			testName:      "DeleteEvent_positive_paramValue_5",
			r:             httptest.NewRequest("DELETE", "/5", nil),
			w:             httptest.NewRecorder(),
			expID:         5,
			expStatusCode: 204,
		},
		{
			testName:      "DeleteEvent_positive_paramValue_450",
			r:             httptest.NewRequest("DELETE", "/450", nil),
			w:             httptest.NewRecorder(),
			expID:         450,
			expStatusCode: 204,
		},
		{
			testName:      "DeleteEvent_Unexpected",
			r:             httptest.NewRequest("DELETE", "/450", nil),
			w:             httptest.NewRecorder(),
			expID:         450,
			mockErrReturn: customErrors.ErrUnexpected,
			expJSONReturn: `{"ErrorType":"Unexpected","ErrorMessage":"an unexpected error occurred"}`,
			expStatusCode: 500,
		},
		{
			testName:      "DeleteEvent_NotFound",
			r:             httptest.NewRequest("DELETE", "/450", nil),
			w:             httptest.NewRecorder(),
			expID:         450,
			mockErrReturn: customErrors.ErrNotFound,
			expJSONReturn: `{"ErrorType":"NotFound","ErrorMessage":"the server cannot find the requested resource"}`,
			expStatusCode: 404,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {

			// mock business logic
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockBL := mocks.NewMockBusinessLogicInterface(mockCtrl)
			mockBL.EXPECT().DeleteEvent(gomock.Any(), tc.expID).Return(tc.mockErrReturn)

			// create handler with mocks
			handler := InitHandler(mockBL)
			handler.Mux.ServeHTTP(tc.w, tc.r)

			resp := tc.w.Result()

			data, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Error(err)
			}

			if tc.expJSONReturn != "" {
				require.JSONEq(t, tc.expJSONReturn, string(data), "JSON data should to be equal")

			}

			require.Equal(t, tc.expStatusCode, resp.StatusCode, "Wrong status code returned")

		})
	}
}

func TestUpdateEventHandler(t *testing.T) {
	testCases := []struct {
		testName         string
		r                *http.Request
		w                *httptest.ResponseRecorder
		decodeErrPresent bool
		expID            int64
		eventToMock      types.Event
		mockErrReturn    error
		expJSONReturn    string
		expStatusCode    int
	}{
		{
			testName: "UpdateEvent_with_event_return",
			r:        httptest.NewRequest("PUT", "/5", bytes.NewBuffer([]byte(`{"id":1,"name":"Onboarding Meeting","startTime":"2022-09-14T09:00:00Z","endTime":"2022-09-14T09:00:00Z", "alertTime":"0001-01-01T00:00:00Z"}`))),
			w:        httptest.NewRecorder(),
			eventToMock: types.Event{
				ID:        1,
				Name:      "Onboarding Meeting",
				StartTime: time.Date(2022, 9, 14, 9, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2022, 9, 14, 9, 0, 0, 0, time.UTC),
			},
			expID:         5,
			expJSONReturn: `{"id":1,"name":"Onboarding Meeting","startTime":"2022-09-14T09:00:00Z","endTime":"2022-09-14T09:00:00Z", "alertTime":"0001-01-01T00:00:00Z"}`,
			expStatusCode: 200,
		},
		{
			testName: "UpdateEvent_Unexpected",
			r:        httptest.NewRequest("PUT", "/100", bytes.NewBuffer([]byte(`{"id":1,"name":"Supposedly too long Meeting Name","startTime":"2022-09-14T09:00:00Z","endTime":"2022-09-14T09:00:00Z", "alertTime":"0001-01-01T00:00:00Z"}`))),
			w:        httptest.NewRecorder(),
			eventToMock: types.Event{
				ID:        1,
				Name:      "Supposedly too long Meeting Name",
				StartTime: time.Date(2022, 9, 14, 9, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2022, 9, 14, 9, 0, 0, 0, time.UTC),
			},
			expID:         100,
			mockErrReturn: customErrors.ErrUnexpected,
			expJSONReturn: `{"ErrorType":"Unexpected","ErrorMessage":"an unexpected error occurred"}`,
			expStatusCode: 500,
		},
		{
			testName: "UpdateEvent_NotFound",
			r:        httptest.NewRequest("PUT", "/1000", bytes.NewBuffer([]byte(`{"id":1,"name":"Meeting Name","startTime":"2022-09-14T09:00:00Z","endTime":"2022-09-14T09:00:00Z", "alertTime":"0001-01-01T00:00:00Z"}`))),
			w:        httptest.NewRecorder(),
			eventToMock: types.Event{
				ID:        1,
				Name:      "Meeting Name",
				StartTime: time.Date(2022, 9, 14, 9, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2022, 9, 14, 9, 0, 0, 0, time.UTC),
			},
			expID:         1000,
			mockErrReturn: customErrors.ErrNotFound,
			expJSONReturn: `{"ErrorType":"NotFound","ErrorMessage":"the server cannot find the requested resource"}`,
			expStatusCode: 404,
		},
		{
			testName: "UpdateEvent_BadRequest",
			r:        httptest.NewRequest("PUT", "/1000", bytes.NewBuffer([]byte(`{"id":1,"name":"Meeting Name","startTime":"2022-09-14T09:00:00Z","endTime":"2022-09-14T09:00:00Z", "alertTime":"0001-01-01T00:00:00Z"}`))),
			w:        httptest.NewRecorder(),
			eventToMock: types.Event{
				ID:        1,
				Name:      "Meeting Name",
				StartTime: time.Date(2022, 9, 14, 9, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2022, 9, 14, 9, 0, 0, 0, time.UTC),
			},
			expID:         1000,
			mockErrReturn: customErrors.ErrBadRequest,
			expJSONReturn: `{"ErrorType":"BadRequest","ErrorMessage":"the server cannot process the request"}`,
			expStatusCode: 400,
		},
		{
			testName:         "UpdateEvent_DecodeErr1",
			r:                httptest.NewRequest("PUT", "/5", bytes.NewBuffer([]byte(`{json string}`))),
			w:                httptest.NewRecorder(),
			decodeErrPresent: true,
			expJSONReturn:    `{"ErrorType":"BadRequest","ErrorMessage":"the server cannot process the request"}`,
			expStatusCode:    400,
		},
		{
			testName:         "UpdateEvent_DecodeErr2",
			r:                httptest.NewRequest("PUT", "/10", bytes.NewBuffer([]byte(`{"hello": world}`))),
			w:                httptest.NewRecorder(),
			decodeErrPresent: true,
			expJSONReturn:    `{"ErrorType":"BadRequest","ErrorMessage":"the server cannot process the request"}`,
			expStatusCode:    400,
		},
		{
			testName:         "UpdateEvent_DecodeErr3",
			r:                httptest.NewRequest("PUT", "/15", bytes.NewBuffer([]byte(`{"name": false}`))),
			w:                httptest.NewRecorder(),
			decodeErrPresent: true,
			expJSONReturn:    `{"ErrorType":"BadRequest","ErrorMessage":"the server cannot process the request"}`,
			expStatusCode:    400,
		},
		{
			testName:         "UpdateEvent_StrconvErr",
			r:                httptest.NewRequest("PUT", "/1.5", nil),
			w:                httptest.NewRecorder(),
			decodeErrPresent: true,
			expJSONReturn:    `{"ErrorType":"BadRequest","ErrorMessage":"the server cannot process the request"}`,
			expStatusCode:    400,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {

			// mock business logic
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockBL := mocks.NewMockBusinessLogicInterface(mockCtrl)

			if !tc.decodeErrPresent {
				mockBL.EXPECT().UpdateEvent(gomock.Any(), tc.eventToMock, tc.expID).Return(tc.mockErrReturn)
			}

			// create handler with mocks
			handler := InitHandler(mockBL)
			handler.Mux.ServeHTTP(tc.w, tc.r)

			resp := tc.w.Result()

			data, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Error(err)
			}

			require.JSONEq(t, tc.expJSONReturn, string(data), "JSON data should be equal")

			require.Equal(t, tc.expStatusCode, resp.StatusCode, "Wrong status code returned")

		})
	}
}
