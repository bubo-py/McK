package repositories

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/bubo-py/McK/types"
)

var ctx context.Context

func TestAppendEvent(t *testing.T) {
	db := InitDatabase()
	ti := time.Date(2022, 9, 16, 20, 30, 0, 0, time.Local)

	event := types.Event{
		ID:          100,
		Name:        "Daily meeting",
		StartTime:   ti,
		EndTime:     ti,
		Description: "A daily meeting for backend team",
		AlertTime:   ti,
	}

	event2 := types.Event{
		ID:          200,
		Name:        "Weekly meeting",
		StartTime:   ti,
		EndTime:     ti,
		Description: "A Weekly meeting for frontend team",
		AlertTime:   ti,
	}
	_ = db.AddEvent(ctx, event)
	_ = db.AddEvent(ctx, event2)

	e, _ := db.GetEvents(ctx)

	if len(e) < 2 {
		t.Error("Failed to add an event")
	}
}

func TestDeleteEvent(t *testing.T) {
	ti := time.Date(2022, 9, 16, 20, 30, 0, 0, time.Local)

	testCases := []struct {
		id        int64
		expLength int
		expError  error
	}{
		{1, 1, nil},
		{2, 1, nil},
		{8, 2, errors.New("event with specified id not found")},
	}
	for _, tc := range testCases {
		testName := fmt.Sprintf("Delete id %d", tc.id)
		t.Run(testName, func(t *testing.T) {
			db := InitDatabase()

			event := types.Event{
				ID:          300,
				Name:        "Daily meeting",
				StartTime:   ti,
				EndTime:     ti,
				Description: "A daily meeting for backend team",
				AlertTime:   ti,
			}

			_ = db.AddEvent(ctx, event)
			_ = db.AddEvent(ctx, event)

			err := db.DeleteEvent(ctx, tc.id)

			e, _ := db.GetEvents(ctx)
			if len(e) != tc.expLength {
				t.Errorf("Failed to delete an event: got length: %v, expected: %v", len(e), tc.expLength)
			}

			if err != nil {
				if err.Error() != tc.expError.Error() {
					t.Errorf("Should return different error: got: %v, expected: %v", err, tc.expError)
				}
			}
		})
	}
}

func TestUpdateEvent(t *testing.T) {
	ti := time.Date(2022, 9, 16, 20, 30, 0, 0, time.Local)

	testCases := []struct {
		id       int64
		expName  string
		expError error
	}{
		{1, "Updated event", nil},
		{8, "", errors.New("event with specified id not found")},
	}
	for _, tc := range testCases {
		testName := fmt.Sprintf("Delete id %d", tc.id)
		t.Run(testName, func(t *testing.T) {
			db := InitDatabase()

			event := types.Event{
				ID:          300,
				Name:        "Daily meeting",
				StartTime:   ti,
				EndTime:     ti,
				Description: "A daily meeting for backend team",
				AlertTime:   ti,
			}

			uEvent := types.Event{
				ID:          300,
				Name:        "Updated event",
				StartTime:   ti,
				EndTime:     ti,
				Description: "A daily meeting for backend team",
				AlertTime:   ti,
			}

			_ = db.AddEvent(ctx, event)
			_ = db.AddEvent(ctx, event)

			err := db.UpdateEvent(ctx, uEvent, tc.id)
			if err != nil {
				if err.Error() != tc.expError.Error() {
					t.Errorf("Should return different error: got: %v, expected: %v", err, tc.expError)
				}
			}

			e, _ := db.GetEvent(ctx, tc.id)
			if e.Name != tc.expName {
				t.Errorf("Failed to update an event: got name: %v, expected: %v", e.Name, tc.expName)
			}

		})
	}
}
