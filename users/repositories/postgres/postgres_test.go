package postgres

import (
	"context"
	"errors"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/bubo-py/McK/types"
)

var (
	deleteErr = errors.New("user with specified id not found")
	authErr   = errors.New("incorrect credentials")
)

func TestMain(m *testing.M) {
	// Setup
	ctx := context.Background()

	db, err := Init(ctx, os.Getenv("PGURL"))
	if err != nil {
		log.Fatalf("Could not initialize database: %v", err)
	}

	_, _ = db.pool.Exec(ctx, "DROP TABLE users")
	_, _ = db.pool.Exec(ctx, "DROP TABLE users_migration")
	_ = RunMigration(ctx, db)

	code := m.Run()

	// Tear down

	os.Exit(code)
}

func TestAddUser(t *testing.T) {
	ctx := context.Background()

	db, err := Init(ctx, os.Getenv("PGURL"))
	if err != nil {
		t.Error(err)
	}

	deleteAllUsers(ctx, db)

	user := types.User{
		ID:       158,
		Login:    "Hello",
		Password: "Hello",
		Timezone: "Asia/Tokyo",
	}

	user2 := types.User{
		ID:       2,
		Login:    "Second User",
		Password: "Hello",
		Timezone: "Europe/London",
	}

	u, err := db.AddUser(ctx, user)
	if err != nil {
		t.Error(err)
	}

	if u.ID != 1 {
		t.Errorf("Failed to add an user properly: got ID: %v, expected: %v", u.ID, 1)
	}

	u, err = db.AddUser(ctx, user2)
	if err != nil {
		t.Error(err)
	}

	if u.ID != 2 || u.Login != user2.Login {
		t.Errorf("Failed to add an user properly: got ID: %v, expected: %v "+
			"got login: %v, expected: %v", u.ID, 2, u.Login, user2.Login)
	}

	err = db.DeleteUser(ctx, 1)
	if err != nil {
		t.Error(err)
	}

	err = db.DeleteUser(ctx, 1)
	if err == nil {
		t.Errorf("Should return an error: %v", deleteErr)
	} else {
		if err.Error() != deleteErr.Error() {
			t.Errorf("Should return a different error: got %v, expected: %v", err, deleteErr)
		}
	}
}

func TestUpdateUser(t *testing.T) {
	ctx := context.Background()

	db, err := Init(ctx, os.Getenv("PGURL"))
	if err != nil {
		t.Error(err)
	}

	deleteAllUsers(ctx, db)

	user := types.User{
		ID:       158,
		Login:    "Hello",
		Password: "Hello",
		Timezone: "Asia/Tokyo",
	}

	user2 := types.User{
		ID:       23,
		Login:    "Hey",
		Password: "Hello",
		Timezone: "Europe/London",
	}

	fullUserUpdate := types.User{
		ID:       255,
		Login:    "Updated User",
		Password: "Hello",
		Timezone: "Africa/Ouagadougou",
	}

	partialUserUpdate := types.User{
		ID:       235,
		Login:    "",
		Password: "Hello",
		Timezone: "",
	}

	_, _ = db.AddUser(ctx, user)
	_, _ = db.AddUser(ctx, user2)

	u, err := db.UpdateUser(ctx, fullUserUpdate, 1)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(fullUserUpdate, u) {
		t.Errorf("Failed to update user: got: %v, expected: %v", u, fullUserUpdate)
	}

	_, err = db.UpdateUser(ctx, partialUserUpdate, 2)
	if err != nil {
		t.Error(err)
	}

	u, _ = db.GetUserByLogin(ctx, user2.Login)

	if u.Login != user2.Login {
		t.Errorf("Failed to partialy update user: got: %v, expected: %v", u.Login, user2.Login)
	}
}

func TestGetUserByLogin(t *testing.T) {
	ctx := context.Background()

	db, err := Init(ctx, os.Getenv("PGURL"))
	if err != nil {
		t.Error(err)
	}

	deleteAllUsers(ctx, db)

	user := types.User{
		ID:       158,
		Login:    "Hello",
		Password: "Hello",
		Timezone: "Asia/Tokyo",
	}

	user2 := types.User{
		ID:       2,
		Login:    "Check The Login",
		Password: "Hello",
		Timezone: "Europe/London",
	}

	_, _ = db.AddUser(ctx, user)
	_, _ = db.AddUser(ctx, user2)

	u, err := db.GetUserByLogin(ctx, "Check The Login")
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(user2, u) {
		t.Errorf("Failed to retrive user by login: got: %v, expected: %v", u, user2)
	}

	u, err = db.GetUserByLogin(ctx, "Hello---World")
	if err == nil {
		t.Errorf("Should return an error: %v", authErr)
	} else {
		if err.Error() != authErr.Error() {
			t.Errorf("Should return a different error: got %v, expected: %v", err, authErr)
		}
	}
}

func deleteAllUsers(ctx context.Context, pg Db) {
	query := "TRUNCATE users RESTART IDENTITY"

	_, err := pg.pool.Exec(ctx, query)
	if err != nil {
		log.Fatal(err)
	}
}
