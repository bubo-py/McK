package events

import (
	"testing"

	"github.com/bubo-py/McK/events"
)

func TestDeleteEvent(t *testing.T) {
	event := events.Event{
		ID:          300,
		Name:        "Daily meeting",
		StartTime:   "2022-09-14T09:00:00.000Z",
		EndTime:     "2022-09-14T09:30:00.000Z",
		Description: "A daily meeting for backend team",
		AlertTime:   "2022-09-14T08:45:00.000Z",
	}

	event2 := events.Event{
		ID:          400,
		Name:        "Weekly meeting",
		StartTime:   "2022-09-16T19:00:00.000Z",
		EndTime:     "2022-09-16T19:30:00.000Z",
		Description: "A Weekly meeting for frontend team",
		AlertTime:   "2022-09-16T18:45:00.000Z",
	}
	events.AppendEvent(event)
	events.AppendEvent(event2)

	events.DeleteEvent(2)

	if len(events.Db)%2 == 0 {
		t.Error("Failed to delete an event")
	}
}
