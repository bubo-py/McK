package handlers

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bubo-py/McK/customErrors"
	"github.com/bubo-py/McK/events/repositories/mocks"
	"github.com/bubo-py/McK/types"
	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetEventsHandler(t *testing.T) {
	testCases := []struct {
		r              *http.Request
		w              *httptest.ResponseRecorder
		mockDataReturn []types.Event
		mockErrReturn  error
		expFilters     types.Filters
		expJSONReturn  string
		expStatusCode  int
	}{
		{
			r: httptest.NewRequest("GET", "/api/events", nil),
			w: httptest.NewRecorder(),
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
			r:             httptest.NewRequest("GET", "/api/events?day=15&month=10", nil),
			w:             httptest.NewRecorder(),
			mockErrReturn: nil,
			expFilters:    types.Filters{Day: 15, Month: 10, Year: 0},
			expJSONReturn: "null\n",
			expStatusCode: 200,
		},
		{
			r:             httptest.NewRequest("GET", "/api/events?year=2222", nil),
			w:             httptest.NewRecorder(),
			mockErrReturn: nil,
			expFilters:    types.Filters{Day: 0, Month: 0, Year: 2222},
			expJSONReturn: "null\n",
			expStatusCode: 200,
		},
		{
			r:             httptest.NewRequest("GET", "/api/events?day=100", nil),
			w:             httptest.NewRecorder(),
			mockErrReturn: customErrors.BadRequest,
			expFilters:    types.Filters{Day: 100, Month: 0, Year: 0},
			expJSONReturn: `{"ErrorType":"BadRequest","ErrorMessage":"the server cannot process the request"}`,
			expStatusCode: 400,
		},
	}
	for i, tc := range testCases {
		testName := fmt.Sprintf("Test %d", i+1)
		t.Run(testName, func(t *testing.T) {

			// mock business logic
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockBL := mocks.NewMockBusinessLogicInterface(mockCtrl)
			mockBL.EXPECT().GetEvents(tc.r.Context(), tc.expFilters).Return(tc.mockDataReturn, tc.mockErrReturn).Times(1)

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

func TestGetEventsHandlerWithStrconvError(t *testing.T) {
	testCases := []struct {
		r             *http.Request
		w             *httptest.ResponseRecorder
		expJSONReturn string
		expStatusCode int
	}{
		{
			r:             httptest.NewRequest("GET", "/api/events?day=1.5", nil),
			w:             httptest.NewRecorder(),
			expJSONReturn: `{"ErrorType":"BadRequest","ErrorMessage":"the server cannot process the request"}`,
			expStatusCode: 400,
		},
		{
			r:             httptest.NewRequest("GET", "/api/events?month=3.14", nil),
			w:             httptest.NewRecorder(),
			expJSONReturn: `{"ErrorType":"BadRequest","ErrorMessage":"the server cannot process the request"}`,
			expStatusCode: 400,
		},
		{
			r:             httptest.NewRequest("GET", "/api/events?year=20.22", nil),
			w:             httptest.NewRecorder(),
			expJSONReturn: `{"ErrorType":"BadRequest","ErrorMessage":"the server cannot process the request"}`,
			expStatusCode: 400,
		},
		{
			r:             httptest.NewRequest("GET", "/api/events?day=1.5&year=20.22&month=3.14", nil),
			w:             httptest.NewRecorder(),
			expJSONReturn: `{"ErrorType":"BadRequest","ErrorMessage":"the server cannot process the request"}`,
			expStatusCode: 400,
		},
	}
	for i, tc := range testCases {
		testName := fmt.Sprintf("Test %d", i+1)
		t.Run(testName, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockBL := mocks.NewMockBusinessLogicInterface(mockCtrl)

			handler := InitHandler(mockBL)
			handler.GetEventsHandler(tc.w, tc.r)

			resp := tc.w.Result()

			data, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Error(err)
			}

			require.JSONEq(t, tc.expJSONReturn, string(data), "JSON data should to be equal")

			require.Equal(t, tc.expStatusCode, resp.StatusCode, "Wrong status code returned")
		})
	}
}

func TestGetEventHandler(t *testing.T) {
	testCases := []struct {
		URLParamValue string
		expID         int64
		eventToMock   types.Event
		mockErrReturn error
		expJSONReturn string
		expStatusCode int
	}{
		{
			URLParamValue: "5",
			expID:         5,
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
			URLParamValue: "450",
			expID:         450,
			expStatusCode: 200,
		},
		{
			URLParamValue: "450",
			expID:         450,
			mockErrReturn: customErrors.ErrUnexpected,
			expJSONReturn: `{"ErrorType":"Unexpected","ErrorMessage":"an unexpected error occurred"}`,
			expStatusCode: 404,
		},
	}
	for i, tc := range testCases {
		testName := fmt.Sprintf("Test %d", i+1)
		t.Run(testName, func(t *testing.T) {

			r := httptest.NewRequest("GET", "/api/events", nil)
			w := httptest.NewRecorder()

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tc.URLParamValue)

			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			// mock business logic
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockBL := mocks.NewMockBusinessLogicInterface(mockCtrl)
			mockBL.EXPECT().GetEvent(r.Context(), tc.expID).Return(tc.eventToMock, tc.mockErrReturn)

			// create handler with mocks
			handler := InitHandler(mockBL)
			handler.GetEventHandler(w, r)

			resp := w.Result()

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

func TestGetEventHandlerWithDecodeError(t *testing.T) {
	testCases := []struct {
		jsonStr       string
		expJSONReturn string
		expStatusCode int
	}{
		{
			jsonStr:       `{json string}`,
			expJSONReturn: `{"ErrorType":"BadRequest","ErrorMessage":"the server cannot process the request"}`,
			expStatusCode: 400,
		},
		{
			jsonStr:       `{"hello": world}`,
			expJSONReturn: `{"ErrorType":"BadRequest","ErrorMessage":"the server cannot process the request"}`,
			expStatusCode: 400,
		},
		{
			jsonStr:       `{"name": false}`,
			expJSONReturn: `{"ErrorType":"BadRequest","ErrorMessage":"the server cannot process the request"}`,
			expStatusCode: 400,
		},
	}
	for i, tc := range testCases {
		testName := fmt.Sprintf("Test %d", i+1)
		t.Run(testName, func(t *testing.T) {

			// mock request
			r := httptest.NewRequest("POST", "/api/events", bytes.NewBuffer([]byte(tc.jsonStr)))
			w := httptest.NewRecorder()

			// mock business logic
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockBL := mocks.NewMockBusinessLogicInterface(mockCtrl)

			// create handler with mocks
			handler := InitHandler(mockBL)
			handler.GetEventHandler(w, r)

			resp := w.Result()

			data, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Error(err)
			}

			require.JSONEqf(t, tc.expJSONReturn, string(data), "JSON data should to be equal")

			require.Equal(t, tc.expStatusCode, resp.StatusCode, "Wrong status code returned")

		})
	}
}

func TestAddEventHandler(t *testing.T) {
	testCases := []struct {
		jsonStr       string
		eventToMock   types.Event
		mockErrReturn error
		expJSONReturn string
		expStatusCode int
	}{
		{
			jsonStr: `{"id":1,"name":"Onboarding Meeting","startTime":"2022-09-14T09:00:00Z","endTime":"2022-09-14T09:00:00Z", "alertTime":"0001-01-01T00:00:00Z"}`,
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
			jsonStr: `{"id":1,"name":"Supposedly too long Meeting Name","startTime":"2022-09-14T09:00:00Z","endTime":"2022-09-14T09:00:00Z", "alertTime":"0001-01-01T00:00:00Z"}`,
			eventToMock: types.Event{
				ID:        1,
				Name:      "Supposedly too long Meeting Name",
				StartTime: time.Date(2022, 9, 14, 9, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2022, 9, 14, 9, 0, 0, 0, time.UTC),
			},
			mockErrReturn: customErrors.BadRequest,
			expJSONReturn: `{"ErrorType":"BadRequest","ErrorMessage":"the server cannot process the request"}`,
			expStatusCode: 400,
		},
	}
	for i, tc := range testCases {
		testName := fmt.Sprintf("Test %d", i+1)
		t.Run(testName, func(t *testing.T) {

			// mock request
			r := httptest.NewRequest("POST", "/api/events", bytes.NewBuffer([]byte(tc.jsonStr)))
			w := httptest.NewRecorder()

			// mock business logic
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockBL := mocks.NewMockBusinessLogicInterface(mockCtrl)
			mockBL.EXPECT().AddEvent(r.Context(), tc.eventToMock).Return(tc.mockErrReturn)

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

func TestAddEventHandlerWithDecodeError(t *testing.T) {
	testCases := []struct {
		jsonStr       string
		expJSONReturn string
		expStatusCode int
	}{
		{
			jsonStr:       `{json string}`,
			expJSONReturn: `{"ErrorType":"BadRequest","ErrorMessage":"the server cannot process the request"}`,
			expStatusCode: 400,
		},
		{
			jsonStr:       `{"hello": world}`,
			expJSONReturn: `{"ErrorType":"BadRequest","ErrorMessage":"the server cannot process the request"}`,
			expStatusCode: 400,
		},
		{
			jsonStr:       `{"name": false}`,
			expJSONReturn: `{"ErrorType":"BadRequest","ErrorMessage":"the server cannot process the request"}`,
			expStatusCode: 400,
		},
	}
	for i, tc := range testCases {
		testName := fmt.Sprintf("Test %d", i+1)
		t.Run(testName, func(t *testing.T) {

			// mock request
			r := httptest.NewRequest("POST", "/api/events", bytes.NewBuffer([]byte(tc.jsonStr)))
			w := httptest.NewRecorder()

			// mock business logic
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockBL := mocks.NewMockBusinessLogicInterface(mockCtrl)

			// create handler with mocks
			handler := InitHandler(mockBL)
			handler.AddEventHandler(w, r)

			resp := w.Result()

			data, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Error(err)
			}

			require.JSONEqf(t, tc.expJSONReturn, string(data), "JSON data should to be equal")

			require.Equal(t, tc.expStatusCode, resp.StatusCode, "Wrong status code returned")

		})
	}
}

func TestDeleteEventHandler(t *testing.T) {
	testCases := []struct {
		URLParamValue string
		expID         int64
		mockErrReturn error
		expJSONReturn string
		expStatusCode int
	}{
		{
			URLParamValue: "5",
			expID:         5,
			expStatusCode: 204,
		},
		{
			URLParamValue: "450",
			expID:         450,
			expStatusCode: 204,
		},
		{
			URLParamValue: "450",
			expID:         450,
			mockErrReturn: customErrors.ErrUnexpected,
			expJSONReturn: `{"ErrorType":"Unexpected","ErrorMessage":"an unexpected error occurred"}`,
			expStatusCode: 404,
		},
	}
	for i, tc := range testCases {
		testName := fmt.Sprintf("Test %d", i+1)
		t.Run(testName, func(t *testing.T) {

			r := httptest.NewRequest("DELETE", "/api/events", nil)
			w := httptest.NewRecorder()

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tc.URLParamValue)

			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			// mock business logic
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockBL := mocks.NewMockBusinessLogicInterface(mockCtrl)
			mockBL.EXPECT().DeleteEvent(r.Context(), tc.expID).Return(tc.mockErrReturn)

			// create handler with mocks
			handler := InitHandler(mockBL)
			handler.DeleteEventHandler(w, r)

			resp := w.Result()

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
		jsonStr       string
		URLParamValue string
		expID         int64
		eventToMock   types.Event
		mockErrReturn error
		expJSONReturn string
		expStatusCode int
	}{
		{
			jsonStr:       `{"id":1,"name":"Onboarding Meeting","startTime":"2022-09-14T09:00:00Z","endTime":"2022-09-14T09:00:00Z", "alertTime":"0001-01-01T00:00:00Z"}`,
			URLParamValue: "5",
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
			jsonStr:       `{"id":1,"name":"Supposedly too long Meeting Name","startTime":"2022-09-14T09:00:00Z","endTime":"2022-09-14T09:00:00Z", "alertTime":"0001-01-01T00:00:00Z"}`,
			URLParamValue: "100",
			eventToMock: types.Event{
				ID:        1,
				Name:      "Supposedly too long Meeting Name",
				StartTime: time.Date(2022, 9, 14, 9, 0, 0, 0, time.UTC),
				EndTime:   time.Date(2022, 9, 14, 9, 0, 0, 0, time.UTC),
			},
			expID:         100,
			mockErrReturn: customErrors.ErrUnexpected,
			expJSONReturn: `{"ErrorType":"Unexpected","ErrorMessage":"an unexpected error occurred"}`,
			expStatusCode: 404,
		},
	}
	for i, tc := range testCases {
		testName := fmt.Sprintf("Test %d", i+1)
		t.Run(testName, func(t *testing.T) {

			r := httptest.NewRequest("UPDATE", "/api/events", bytes.NewBuffer([]byte(tc.jsonStr)))
			w := httptest.NewRecorder()

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tc.URLParamValue)

			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			// mock business logic
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockBL := mocks.NewMockBusinessLogicInterface(mockCtrl)
			mockBL.EXPECT().UpdateEvent(r.Context(), tc.eventToMock, tc.expID).Return(tc.mockErrReturn)

			// create handler with mocks
			handler := InitHandler(mockBL)
			handler.UpdateEventHandler(w, r)

			resp := w.Result()

			data, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Error(err)
			}

			require.JSONEq(t, tc.expJSONReturn, string(data), "JSON data should be equal")

			require.Equal(t, tc.expStatusCode, resp.StatusCode, "Wrong status code returned")

		})
	}
}

func TestUpdateEventHandlerWithDecodeError(t *testing.T) {
	testCases := []struct {
		URLParamValue string
		jsonStr       string
		expJSONReturn string
		expStatusCode int
	}{
		{
			URLParamValue: "5",
			jsonStr:       `{json string}`,
			expJSONReturn: `{"ErrorType":"BadRequest","ErrorMessage":"the server cannot process the request"}`,
			expStatusCode: 400,
		},
		{
			URLParamValue: "10",
			jsonStr:       `{"hello": world}`,
			expJSONReturn: `{"ErrorType":"BadRequest","ErrorMessage":"the server cannot process the request"}`,
			expStatusCode: 400,
		},
		{
			URLParamValue: "15",
			jsonStr:       `{"name": false}`,
			expJSONReturn: `{"ErrorType":"BadRequest","ErrorMessage":"the server cannot process the request"}`,
			expStatusCode: 400,
		},
		{
			URLParamValue: "1.5",
			expStatusCode: 400,
		},
	}
	for i, tc := range testCases {
		testName := fmt.Sprintf("Test %d", i+1)
		t.Run(testName, func(t *testing.T) {

			r := httptest.NewRequest("UPDATE", "/api/events", bytes.NewBuffer([]byte(tc.jsonStr)))
			w := httptest.NewRecorder()

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tc.URLParamValue)

			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			// mock business logic
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockBL := mocks.NewMockBusinessLogicInterface(mockCtrl)

			// create handler with mocks
			handler := InitHandler(mockBL)
			handler.UpdateEventHandler(w, r)

			resp := w.Result()

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
