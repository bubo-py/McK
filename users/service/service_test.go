package service

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/bubo-py/McK/types"
	"github.com/bubo-py/McK/users/repositories/serviceDb"
)

var loginErr = errors.New("login should be at least 3 and contain up to 30 characters")
var passwordErr = errors.New("password should be at least 5 characters")

var ctx = context.Background()
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

			_, err := bl.AddUser(ctx, tc.user)
			if err != nil {
				if err.Error() != tc.expError.Error() {
					t.Errorf("Failed to add user: got error: %v, expected: %v", err, tc.expError)
				}
			}
		})
	}
}
