package postgres

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/bubo-py/McK/types"
)

func TestMain(m *testing.M) {
	// Setup
	ctx := context.Background()

	db, _ := PostgresInit(ctx)

	_, _ = db.pool.Exec(ctx, "DROP TABLE events")
	_, _ = db.pool.Exec(ctx, "DROP TABLE migration")
	_ = RunMigration(ctx, db)

	code := m.Run()

	// Tear down

	os.Exit(code)
}

func TestPostgresDb_GetEvent(t *testing.T) {
	ti := time.Date(2010, 9, 16, 20, 30, 0, 0, time.Local)

	ctx := context.Background()
	db, err := PostgresInit(ctx)
	if err != nil {
		t.Error(err)
	}

	event := types.Event{
		ID:          100,
		Name:        "Initial meeting",
		StartTime:   ti,
		EndTime:     ti,
		Description: "A meeting",
		AlertTime:   ti,
	}

	event2 := types.Event{
		ID:          200,
		Name:        "Second meeting",
		StartTime:   ti,
		EndTime:     ti,
		Description: "A meeting",
		AlertTime:   ti,
	}

	err = db.AddEvent(ctx, event)
	if err != nil {
		t.Error(err)
	}

	err = db.AddEvent(ctx, event2)
	if err != nil {
		t.Error(err)
	}

	var id int64 = 2
	e, err := db.GetEvent(ctx, id)
	if err != nil {
		t.Error(err)
	}

	if e.ID != id || e.StartTime.Year() != 2010 {
		t.Error("Failed to fetch an event with given id")
	}

	err = db.UpdateEvent(ctx, event, 20)
	if err == nil {
		t.Errorf("Error is nil, should have: %s", "event with specified id not found")
	}
}

func TestPostgresDb_DeleteEvent(t *testing.T) {
	ti := time.Date(2022, 9, 16, 20, 30, 0, 0, time.Local)
	ti2 := time.Date(2010, 9, 16, 20, 30, 0, 0, time.Local)

	ctx := context.Background()
	db, err := PostgresInit(ctx)
	if err != nil {
		t.Error(err)
	}

	event := types.Event{
		ID:          100,
		Name:        "Initial meeting",
		StartTime:   ti,
		EndTime:     ti,
		Description: "A meeting",
		AlertTime:   ti,
	}

	event2 := types.Event{
		ID:          100,
		Name:        "Second meeting2",
		StartTime:   ti2,
		EndTime:     ti2,
		Description: "A meeting2",
		AlertTime:   ti2,
	}

	err = db.AddEvent(ctx, event)
	if err != nil {
		t.Error(err)
	}

	err = db.AddEvent(ctx, event2)
	if err != nil {
		t.Error(err)
	}

	err = db.DeleteEvent(ctx, 2)
	if err != nil {
		t.Error(err)
	}

	err = db.DeleteEvent(ctx, 2)
	if err == nil {
		t.Errorf("Error is nil, should have: %s", "event with specified id not found")
	}

	e, err := db.GetEvents(ctx)
	if err != nil {
		t.Error(err)
	}

	if len(e)%2 == 0 {
		t.Errorf("Failed to delete an event")
	}

	if len(e) > 4 {
		t.Errorf("Events added incorrectly, should have less than 5, got: %d", len(e))
	}
}
