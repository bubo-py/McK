package events

import (
	"testing"

	"github.com/bubo-py/McK/events"
)

func TestUpdateEvent(t *testing.T) {
	event := events.Event{
		ID:          500,
		Name:        "Updated event",
		StartTime:   "2022-09-14T09:00:00.000",
		EndTime:     "2022-09-14T09:00:00.000",
		Description: "An event that has just been updated",
		AlertTime:   "2022-09-14T09:00:00.000",
	}

	events.UpdateEvent(event, 1)

	if events.Db[0].Name != event.Name {
		t.Error("Failed to update an event")
	}
}
