package handlers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bubo-py/McK/customErrors"
	"github.com/bubo-py/McK/events/repositories/mocks"
	"github.com/bubo-py/McK/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetEventsHandler(t *testing.T) {
	testCases := []struct {
		r              *http.Request
		w              *httptest.ResponseRecorder
		mockDataReturn []types.Event
		mockErrReturn  error
		expFilters     types.Filters
		expJSONReturn  string
		expErrString   string
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
			expErrString:  "null\n",
			expStatusCode: 200,
		},
		{
			r:             httptest.NewRequest("GET", "/api/events?year=2222", nil),
			w:             httptest.NewRecorder(),
			mockErrReturn: nil,
			expFilters:    types.Filters{Day: 0, Month: 0, Year: 2222},
			expErrString:  "null\n",
			expStatusCode: 200,
		},
		{
			r:             httptest.NewRequest("GET", "/api/events?day=100", nil),
			w:             httptest.NewRecorder(),
			mockErrReturn: customErrors.BadRequest,
			expFilters:    types.Filters{Day: 100, Month: 0, Year: 0},
			expErrString:  `"` + customErrors.BadRequest.Error() + `"` + "\n",
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

			if tc.expJSONReturn != "" {
				assert.JSONEqf(t, tc.expJSONReturn, string(data), "JSON data should be equal")
			} else {
				if tc.expErrString != string(data) {
					assert.Equal(t, tc.expErrString, string(data), "Error strings should to be equal")
				}
			}

			assert.Equal(t, tc.expStatusCode, resp.StatusCode, "wrong status code returned")
		})
	}
}

func TestGetEventsHandlerWithStrconvError(t *testing.T) {
	testCases := []struct {
		r             *http.Request
		w             *httptest.ResponseRecorder
		expErrString  string
		expStatusCode int
	}{
		{
			r:             httptest.NewRequest("GET", "/api/events?day=1.5", nil),
			w:             httptest.NewRecorder(),
			expErrString:  `"` + customErrors.BadRequest.Error() + `"` + "\n",
			expStatusCode: 400,
		},
		{
			r:             httptest.NewRequest("GET", "/api/events?month=3.14", nil),
			w:             httptest.NewRecorder(),
			expErrString:  `"` + customErrors.BadRequest.Error() + `"` + "\n",
			expStatusCode: 400,
		},
		{
			r:             httptest.NewRequest("GET", "/api/events?year=20.22", nil),
			w:             httptest.NewRecorder(),
			expErrString:  `"` + customErrors.BadRequest.Error() + `"` + "\n",
			expStatusCode: 400,
		},
		{
			r:             httptest.NewRequest("GET", "/api/events?day=1.5&year=20.22&month=3.14", nil),
			w:             httptest.NewRecorder(),
			expErrString:  `"` + customErrors.BadRequest.Error() + `"` + "\n",
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

			assert.Equal(t, tc.expErrString, string(data), "Error strings should to be equal")

			assert.Equal(t, tc.expStatusCode, resp.StatusCode, "wrong status code returned")
		})
	}
}

func TestAddEventHandler(t *testing.T) {
	testCases := []struct {
		jsonStr       string
		eventToMock   types.Event
		mockErrReturn error
		expJSONReturn string
		expErrString  string
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
			expErrString:  `"` + customErrors.BadRequest.Error() + `"` + "\n",
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

			if tc.expJSONReturn != "" {
				assert.JSONEqf(t, tc.expJSONReturn, string(data), "JSON data should be equal")
			} else {
				if tc.expErrString != string(data) {
					assert.Equal(t, tc.expErrString, string(data), "Error strings should to be equal")
				}
			}

		})
	}
}
