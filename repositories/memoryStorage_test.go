package repositories

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/bubo-py/McK/types"
)

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
	db.AddEvent(event)
	db.AddEvent(event2)

	if len(db.GetEvents()) < 2 {
		t.Error("Failed to add an event")
	}
}

func TestDeleteEvent(t *testing.T) {
	db := InitDatabase()
	ti := time.Date(2022, 9, 16, 20, 30, 0, 0, time.Local)

	testCases := []struct {
		id        int
		expLength int
		expError  error
	}{
		{1, 1, nil},
		{1, 3, errors.New("event with specified id not found")},
		{8, 5, errors.New("event with specified id not found")},
	}
	for _, tc := range testCases {
		testName := fmt.Sprintf("Delete id %d", tc.id)
		t.Run(testName, func(t *testing.T) {
			event := types.Event{
				ID:          300,
				Name:        "Daily meeting",
				StartTime:   ti,
				EndTime:     ti,
				Description: "A daily meeting for backend team",
				AlertTime:   ti,
			}

			db.AddEvent(event)
			db.AddEvent(event)

			err := db.DeleteEvent(tc.id)
			if err != nil {
				if err.Error() != tc.expError.Error() {
					t.Errorf("Should return different error: got: %v, expected: %v", err, tc.expError)
				}
			}

			if len(db.GetEvents()) != tc.expLength {
				t.Errorf("Failed to delete an event: got length: %v, expected: %v", len(db.GetEvents()), tc.expLength)
			}
		})
	}
}

func TestUpdateEvent(t *testing.T) {
	ti := time.Date(2022, 9, 16, 20, 30, 0, 0, time.Local)

	testCases := []struct {
		id       int
		expName  string
		expError error
	}{
		{1, "Updated event", nil},
		{8, "Hello", errors.New("event with specified id not found")},
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

			db.AddEvent(event)
			db.AddEvent(event)

			err := db.UpdateEvent(event, 1)
			if err != nil {
				if err.Error() != tc.expError.Error() {
					t.Errorf("Should return different error: got: %v, expected: %v", err, tc.expError)
				}
			}

			e, err := db.GetEvent(1)
			if err != nil {
				t.Error(err)
			}
			if e.Name != tc.expName {
				t.Errorf("Failed to update an event: got name: %v, expected: %v", len(db.GetEvents()), tc.expName)
			}

		})
	}
}
