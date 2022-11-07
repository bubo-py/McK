package handlers

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bubo-py/McK/customErrors"
	"github.com/bubo-py/McK/types"
	"github.com/bubo-py/McK/users"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestAddUserHandler(t *testing.T) {
	testCases := []struct {
		testName         string
		decodeErrPresent bool
		jsonStr          string
		userToMock       types.User
		mockErrReturn    error
		expJSONReturn    string
		expIDReturn      string
		expStatusCode    int
	}{
		{
			testName: "AddUser_positive_return",
			jsonStr:  `{"id":1,"login":"hello","password":"hello","timezone":"Europe/London"}`,
			userToMock: types.User{
				ID:       1,
				Login:    "hello",
				Password: "hello",
				Timezone: "Europe/London",
			},
			expIDReturn:   "1\n",
			expStatusCode: 200,
		},
		{
			testName: "AddUser_BadRequest",
			jsonStr:  `{"id":2,"login":"TooLongLogin","password":"hello","timezone":"Europe/London"}`,
			userToMock: types.User{
				ID:       2,
				Login:    "TooLongLogin",
				Password: "hello",
				Timezone: "Europe/London",
			},
			mockErrReturn: customErrors.ErrBadRequest,
			expJSONReturn: `{"ErrorType":"BadRequest","ErrorMessage":"the server cannot process the request"}`,
			expStatusCode: 400,
		},
		{
			testName:         "AddUser_DecodeErr1",
			decodeErrPresent: true,
			jsonStr:          `{json string}`,
			expJSONReturn:    `{"ErrorType":"BadRequest","ErrorMessage":"the server cannot process the request"}`,
			expStatusCode:    400,
		},
		{
			testName:         "AddUser_DecodeErr2",
			decodeErrPresent: true,
			jsonStr:          `{"hello": world}`,
			expJSONReturn:    `{"ErrorType":"BadRequest","ErrorMessage":"the server cannot process the request"}`,
			expStatusCode:    400,
		},
		{
			testName:         "AddUser_DecodeErr3",
			decodeErrPresent: true,
			jsonStr:          `{"name": """"false""""}`,
			expJSONReturn:    `{"ErrorType":"BadRequest","ErrorMessage":"the server cannot process the request"}`,
			expStatusCode:    400,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {

			// mock request
			r := httptest.NewRequest("POST", "/api/users", bytes.NewBuffer([]byte(tc.jsonStr)))
			w := httptest.NewRecorder()

			// mock business logic
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockBL := users.NewMockBusinessLogicInterface(mockCtrl)

			if !tc.decodeErrPresent {
				mockBL.EXPECT().AddUser(r.Context(), tc.userToMock).Return(tc.userToMock, tc.mockErrReturn)
			}

			// create handler with mocks
			handler := InitHandler(mockBL)
			handler.AddUserHandler(w, r)

			resp := w.Result()

			data, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Error(err)
			}

			if tc.expJSONReturn != "" {
				require.JSONEqf(t, tc.expJSONReturn, string(data), "JSON data should be equal")
			}

			if tc.expIDReturn != "" {
				require.Equal(t, tc.expIDReturn, string(data))
			}

			require.Equal(t, tc.expStatusCode, resp.StatusCode, "Wrong status code returned")

		})
	}
}

func TestDeleteUserHandler(t *testing.T) {
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
			testName:      "DeleteUser_positive_paramValue_5",
			r:             httptest.NewRequest("DELETE", "/5", nil),
			w:             httptest.NewRecorder(),
			expID:         5,
			expStatusCode: 204,
		},
		{
			testName:      "DeleteUser_positive_paramValue_450",
			r:             httptest.NewRequest("DELETE", "/450", nil),
			w:             httptest.NewRecorder(),
			expID:         450,
			expStatusCode: 204,
		},
		{
			testName:      "DeleteUser_BadRequest",
			r:             httptest.NewRequest("DELETE", "/450", nil),
			w:             httptest.NewRecorder(),
			expID:         450,
			mockErrReturn: customErrors.ErrBadRequest,
			expJSONReturn: `{"ErrorType":"BadRequest","ErrorMessage":"the server cannot process the request"}`,
			expStatusCode: 400,
		},
		{
			testName:      "DeleteUser_Unauthorized",
			r:             httptest.NewRequest("DELETE", "/450", nil),
			w:             httptest.NewRecorder(),
			expID:         450,
			mockErrReturn: customErrors.ErrUnauthorized,
			expJSONReturn: `{"ErrorType":"Unauthorized","ErrorMessage":"the server cannot process the request due to lack of client's access rights"}`,
			expStatusCode: 403,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {

			// mock business logic
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockBL := users.NewMockBusinessLogicInterface(mockCtrl)
			mockBL.EXPECT().DeleteUser(gomock.Any(), tc.expID).Return(tc.mockErrReturn)

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

func TestUpdateUserHandler(t *testing.T) {
	testCases := []struct {
		testName         string
		r                *http.Request
		w                *httptest.ResponseRecorder
		decodeErrPresent bool
		expID            int64
		expIDReturn      string
		userToMock       types.User
		mockErrReturn    error
		expJSONReturn    string
		expStatusCode    int
	}{
		{
			testName: "UpdateUser_positive_return",
			r:        httptest.NewRequest("PUT", "/1", bytes.NewBuffer([]byte(`{"id":1,"login":"hello","password":"hello","timezone":"Europe/London"}`))),
			w:        httptest.NewRecorder(),
			userToMock: types.User{
				ID:       1,
				Login:    "hello",
				Password: "hello",
				Timezone: "Europe/London",
			},
			expID:         1,
			expIDReturn:   "1\n",
			expStatusCode: 200,
		},
		{
			testName: "UpdateUser_BadRequest",
			r:        httptest.NewRequest("PUT", "/2", bytes.NewBuffer([]byte(`{"id":2,"login":"TooLongLogin","password":"hello","timezone":"Europe/London"}`))),
			w:        httptest.NewRecorder(),
			userToMock: types.User{
				ID:       2,
				Login:    "TooLongLogin",
				Password: "hello",
				Timezone: "Europe/London",
			},
			expID:         2,
			mockErrReturn: customErrors.ErrBadRequest,
			expJSONReturn: `{"ErrorType":"BadRequest","ErrorMessage":"the server cannot process the request"}`,
			expStatusCode: 400,
		},
		{
			testName: "UpdateUser_Unauthorized",
			r:        httptest.NewRequest("PUT", "/2", bytes.NewBuffer([]byte(`{"id":2,"login":"TooLongLogin","password":"hello","timezone":"Europe/London"}`))),
			w:        httptest.NewRecorder(),
			userToMock: types.User{
				ID:       2,
				Login:    "TooLongLogin",
				Password: "hello",
				Timezone: "Europe/London",
			},
			expID:         2,
			mockErrReturn: customErrors.ErrUnauthorized,
			expJSONReturn: `{"ErrorType":"Unauthorized","ErrorMessage":"the server cannot process the request due to lack of client's access rights"}`,
			expStatusCode: 403,
		},
		{
			testName:         "UpdateUser_DecodeErr1",
			r:                httptest.NewRequest("PUT", "/5", bytes.NewBuffer([]byte(`{json string}`))),
			w:                httptest.NewRecorder(),
			decodeErrPresent: true,
			expJSONReturn:    `{"ErrorType":"BadRequest","ErrorMessage":"the server cannot process the request"}`,
			expStatusCode:    400,
		},
		{
			testName:         "UpdateUser_DecodeErr2",
			r:                httptest.NewRequest("PUT", "/10", bytes.NewBuffer([]byte(`{"hello": world}`))),
			w:                httptest.NewRecorder(),
			decodeErrPresent: true,
			expJSONReturn:    `{"ErrorType":"BadRequest","ErrorMessage":"the server cannot process the request"}`,
			expStatusCode:    400,
		},
		{
			testName:         "UpdateUser_DecodeErr3",
			r:                httptest.NewRequest("PUT", "/15", bytes.NewBuffer([]byte(`{"name: false"}`))),
			w:                httptest.NewRecorder(),
			decodeErrPresent: true,
			expJSONReturn:    `{"ErrorType":"BadRequest","ErrorMessage":"the server cannot process the request"}`,
			expStatusCode:    400,
		},
		{
			testName:         "UpdateUser_StrconvErr",
			r:                httptest.NewRequest("PUT", "/1.5", nil),
			w:                httptest.NewRecorder(),
			decodeErrPresent: true,
			expStatusCode:    400,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {

			// mock business logic
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockBL := users.NewMockBusinessLogicInterface(mockCtrl)

			if !tc.decodeErrPresent {
				mockBL.EXPECT().UpdateUser(gomock.Any(), tc.userToMock, tc.expID).Return(tc.userToMock, tc.mockErrReturn)
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
				require.JSONEqf(t, tc.expJSONReturn, string(data), "JSON data should be equal")
			}

			if tc.expIDReturn != "" {
				require.Equal(t, tc.expIDReturn, string(data))
			}

			require.Equal(t, tc.expStatusCode, resp.StatusCode, "Wrong status code returned")

		})
	}
}
