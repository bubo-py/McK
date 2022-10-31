package handlers

import (
	"bytes"
	"context"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/bubo-py/McK/customErrors"
	"github.com/bubo-py/McK/types"
	"github.com/bubo-py/McK/users/repositories/mocks"
	"github.com/go-chi/chi"
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

			mockBL := mocks.NewMockBusinessLogicInterface(mockCtrl)

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
		URLParamValue string
		expID         int64
		mockErrReturn error
		expJSONReturn string
		expStatusCode int
	}{
		{
			testName:      "DeleteUser_positive_paramValue_5",
			URLParamValue: "5",
			expID:         5,
			expStatusCode: 204,
		},
		{
			testName:      "DeleteUser_positive_paramValue_450",
			URLParamValue: "450",
			expID:         450,
			expStatusCode: 204,
		},
		{
			testName:      "DeleteUser_BadRequest",
			URLParamValue: "450",
			expID:         450,
			mockErrReturn: customErrors.ErrBadRequest,
			expJSONReturn: `{"ErrorType":"BadRequest","ErrorMessage":"the server cannot process the request"}`,
			expStatusCode: 400,
		},
		{
			testName:      "DeleteUser_Unauthorized",
			URLParamValue: "450",
			expID:         450,
			mockErrReturn: customErrors.ErrUnauthorized,
			expJSONReturn: `{"ErrorType":"Unauthorized","ErrorMessage":"the server cannot process the request due to lack of client's access rights"}`,
			expStatusCode: 403,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {

			r := httptest.NewRequest("DELETE", "/api/users", nil)
			w := httptest.NewRecorder()

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tc.URLParamValue)

			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			// mock business logic
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockBL := mocks.NewMockBusinessLogicInterface(mockCtrl)
			mockBL.EXPECT().DeleteUser(r.Context(), tc.expID).Return(tc.mockErrReturn)

			// create handler with mocks
			handler := InitHandler(mockBL)
			handler.DeleteUserHandler(w, r)

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

func TestUpdateUserHandler(t *testing.T) {
	testCases := []struct {
		testName         string
		decodeErrPresent bool
		jsonStr          string
		URLParamValue    string
		expID            int64
		expIDReturn      string
		userToMock       types.User
		mockErrReturn    error
		expJSONReturn    string
		expStatusCode    int
	}{
		{
			testName:      "UpdateUser_positive_return",
			jsonStr:       `{"id":1,"login":"hello","password":"hello","timezone":"Europe/London"}`,
			URLParamValue: "1",
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
			testName:      "UpdateUser_BadRequest",
			jsonStr:       `{"id":2,"login":"TooLongLogin","password":"hello","timezone":"Europe/London"}`,
			URLParamValue: "2",
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
			testName:      "UpdateUser_Unauthorized",
			jsonStr:       `{"id":2,"login":"TooLongLogin","password":"hello","timezone":"Europe/London"}`,
			URLParamValue: "2",
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
			decodeErrPresent: true,
			URLParamValue:    "5",
			jsonStr:          `{json string}`,
			expJSONReturn:    `{"ErrorType":"BadRequest","ErrorMessage":"the server cannot process the request"}`,
			expStatusCode:    400,
		},
		{
			testName:         "UpdateUser_DecodeErr2",
			decodeErrPresent: true,
			URLParamValue:    "10",
			jsonStr:          `{"hello": world}`,
			expJSONReturn:    `{"ErrorType":"BadRequest","ErrorMessage":"the server cannot process the request"}`,
			expStatusCode:    400,
		},
		{
			testName:         "UpdateUser_DecodeErr3",
			decodeErrPresent: true,
			URLParamValue:    "15",
			jsonStr:          `{"name: false"}`,
			expJSONReturn:    `{"ErrorType":"BadRequest","ErrorMessage":"the server cannot process the request"}`,
			expStatusCode:    400,
		},
		{
			testName:         "UpdateUser_StrconvErr",
			decodeErrPresent: true,
			URLParamValue:    "1.5",
			expStatusCode:    400,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {

			r := httptest.NewRequest("UPDATE", "/api/users", bytes.NewBuffer([]byte(tc.jsonStr)))
			w := httptest.NewRecorder()

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tc.URLParamValue)

			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

			// mock business logic
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockBL := mocks.NewMockBusinessLogicInterface(mockCtrl)

			if !tc.decodeErrPresent {
				mockBL.EXPECT().UpdateUser(r.Context(), tc.userToMock, tc.expID).Return(tc.userToMock, tc.mockErrReturn)
			}

			// create handler with mocks
			handler := InitHandler(mockBL)
			handler.UpdateUserHandler(w, r)

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
