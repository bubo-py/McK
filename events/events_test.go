package events

import (
	"testing"
)

func TestAppendEvent(t *testing.T) {
	event := Event{
		ID:          100,
		Name:        "Daily meeting",
		StartTime:   "2022-09-14T09:00:00.000Z",
		EndTime:     "2022-09-14T09:30:00.000Z",
		Description: "A daily meeting for backend team",
		AlertTime:   "2022-09-14T08:45:00.000Z",
	}

	event2 := Event{
		ID:          200,
		Name:        "Weekly meeting",
		StartTime:   "2022-09-16T19:00:00.000Z",
		EndTime:     "2022-09-16T19:30:00.000Z",
		Description: "A Weekly meeting for frontend team",
		AlertTime:   "2022-09-16T18:45:00.000Z",
	}
	AppendEvent(event)
	AppendEvent(event2)

	if len(Db) < 2 {
		t.Error("Failed to add an event")
	}
}

func TestDeleteEvent(t *testing.T) {
	event := Event{
		ID:          300,
		Name:        "Daily meeting",
		StartTime:   "2022-09-14T09:00:00.000Z",
		EndTime:     "2022-09-14T09:30:00.000Z",
		Description: "A daily meeting for backend team",
		AlertTime:   "2022-09-14T08:45:00.000Z",
	}

	event2 := Event{
		ID:          400,
		Name:        "Weekly meeting",
		StartTime:   "2022-09-16T19:00:00.000Z",
		EndTime:     "2022-09-16T19:30:00.000Z",
		Description: "A Weekly meeting for frontend team",
		AlertTime:   "2022-09-16T18:45:00.000Z",
	}
	AppendEvent(event)
	AppendEvent(event2)

	DeleteEvent(2)

	if len(Db)%2 == 0 {
		t.Error("Failed to delete an event")
	}
}

func TestUpdateEvent(t *testing.T) {
	event := Event{
		ID:          500,
		Name:        "Updated event",
		StartTime:   "2022-09-14T09:00:00.000",
		EndTime:     "2022-09-14T09:00:00.000",
		Description: "An event that has just been updated",
		AlertTime:   "2022-09-14T09:00:00.000",
	}

	AppendEvent(event)
	UpdateEvent(event, 1)

	if Db[0].Name != event.Name {
		t.Error("Failed to update an event")
	}
}
