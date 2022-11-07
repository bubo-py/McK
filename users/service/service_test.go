package service

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/bubo-py/McK/contextHelpers"
	"github.com/bubo-py/McK/types"
	"github.com/bubo-py/McK/users/repositories/serviceDb"
)

var (
	loginErr    = errors.New("the server cannot process the request: login should be at least 3 and contain up to 30 characters")
	passwordErr = errors.New("the server cannot process the request: password should be at least 5 characters")
	authErr     = errors.New("the server cannot process the request due to lack of client's access rights: cannot modify another user's account")
)

var db = serviceDb.Db{}

func TestAddUser(t *testing.T) {
	testCases := []struct {
		user     types.User
		expError error
	}{
		{
			user: types.User{
				Login:    "Hello",
				Password: "Hello",
				Timezone: "",
			},
			expError: nil,
		},
		{
			user: types.User{
				Login:    "pQhTEzNNmAVXg3Yy18sIIJ0JUs59vPu", // length 31
				Password: "Hello",
				Timezone: "",
			},
			expError: loginErr,
		},
		{
			user: types.User{
				Login:    "地的银行将跟进下调首套", // length 33, rune 11
				Password: "Hello",
				Timezone: "",
			},
			expError: nil,
		},
		{
			user: types.User{
				Login:    "ab", // length 2
				Password: "Hello",
				Timezone: "",
			},
			expError: loginErr,
		},
		{
			user: types.User{
				Login:    "Hello",
				Password: "Hey",
				Timezone: "",
			},
			expError: passwordErr,
		},
	}
	for i, tc := range testCases {
		testName := fmt.Sprintf("Test %d", i+1)
		t.Run(testName, func(t *testing.T) {
			bl := InitBusinessLogic(db)
			ctx := context.Background()

			_, err := bl.AddUser(ctx, tc.user)
			if err != nil {
				if err.Error() != tc.expError.Error() {
					t.Errorf("Failed to add user: got error: %v, expected: %v", err, tc.expError)
				}
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	// Init context values
	ctx := context.Background()
	ctx = contextHelpers.WriteLoginToContext(ctx, "hello")
	ctx = contextHelpers.WriteTimezoneToContext(ctx, "Europe/London")

	testCases := []struct {
		user     types.User
		expError error
	}{
		{
			user: types.User{
				Login:    "",
				Password: "Hello",
				Timezone: "Asia/Tokyo",
			},
			expError: authErr,
		},
		{
			user: types.User{
				Login:    "Hello",
				Password: "",
				Timezone: "",
			},
			expError: authErr,
		},
		{
			user: types.User{
				Login:    "x",
				Password: "",
				Timezone: "",
			},
			expError: loginErr,
		},
		{
			user: types.User{
				Login:    "",
				Password: "up",
				Timezone: "",
			},
			expError: passwordErr,
		},
	}
	for i, tc := range testCases {
		testName := fmt.Sprintf("Test %d", i+1)
		t.Run(testName, func(t *testing.T) {
			bl := InitBusinessLogic(db)

			_, err := bl.UpdateUser(ctx, tc.user, 1)
			if err != nil {
				if err.Error() != tc.expError.Error() {
					t.Errorf("Failed to update user: got error: %v, expected: %v", err, tc.expError)
				}
			}
		})
	}
}

func TestDeleteUser(t *testing.T) {
	// Init context values
	ctx := context.Background()
	ctx = contextHelpers.WriteLoginToContext(ctx, "hello")
	ctx = contextHelpers.WriteTimezoneToContext(ctx, "Europe/London")

	bl := InitBusinessLogic(db)

	err := bl.DeleteUser(ctx, 1)
	expErr := authErr
	if err.Error() != expErr.Error() {
		t.Errorf("Failed to delete user: got error: %v, expected: %v", err, expErr)
	}

}
