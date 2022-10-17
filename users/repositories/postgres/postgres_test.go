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

var deleteErr = errors.New("user with specified id not found")

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

func TestPostgresDb(t *testing.T) {
	user := types.User{
		ID:       158,
		Login:    "Hello",
		Password: "Hello",
		Timezone: "CET",
	}

	user2 := types.User{
		ID:       2,
		Login:    "Second User",
		Password: "Hello",
		Timezone: "Europe/London",
	}

	ctx := context.Background()
	db, err := Init(ctx, os.Getenv("PGURL"))
	if err != nil {
		t.Error(err)
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

	u, err = db.GetUserByLogin(ctx, "Second User")
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(user2, u) {
		t.Errorf("Failed to retrive user by login: got: %v, expected: %v", u, user2)
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
			t.Errorf("Should return an error: got %v, expected: %v", err, deleteErr)
		}
	}
}
