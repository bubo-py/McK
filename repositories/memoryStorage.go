package repositories

import (
	"errors"

	"github.com/bubo-py/McK/types"
)

type Database struct {
	ID      int
	Storage []types.Event
}

func InitDatabase() *Database {
	return &Database{}
}

func (db *Database) GetEvents() []types.Event {
	return db.Storage
}

func (db *Database) GetEvent(id int) (types.Event, error) {
	for i, event := range db.Storage {
		if event.ID == id {
			return db.Storage[i], nil
		}
	}
	return types.Event{}, errors.New("event with specified id not found")
}

func (db *Database) AddEvent(e types.Event) {
	db.ID += 1
	e.ID = db.ID
	db.Storage = append(db.Storage, e)
}

func (db *Database) DeleteEvent(id int) error {
	for i, event := range db.Storage {
		if event.ID == id {
			copy(db.Storage[i:], db.Storage[i+1:])
			db.Storage[len(db.Storage)-1] = types.Event{}
			db.Storage = db.Storage[:len(db.Storage)-1]
			return nil
		}
	}
	return errors.New("event with specified id not found")
}

func (db *Database) UpdateEvent(e types.Event, id int) error {
	for i, event := range db.Storage {
		if event.ID == id {
			db.Storage[i].Name = e.Name
			db.Storage[i].StartTime = e.StartTime
			db.Storage[i].EndTime = e.EndTime
			db.Storage[i].Description = e.Description
			db.Storage[i].AlertTime = e.AlertTime
			return nil
		}
	}
	return errors.New("event with specified id not found")
}

func (db *Database) GetEventsByDay(day string) []types.Event {
	filtered := make([]types.Event, 0)

	for _, event := range db.Storage {
		if event.StartTime[8:10] == day {
			filtered = append(filtered, event)
		}
	}
	return filtered
}

func (db *Database) GetEventsByMonth(month string) []types.Event {
	filtered := make([]types.Event, 0)

	for _, event := range db.Storage {
		if event.StartTime[5:7] == month {
			filtered = append(filtered, event)
		}
	}
	return filtered
}

func (db *Database) GetEventsByYear(year string) []types.Event {
	filtered := make([]types.Event, 0)

	for _, event := range db.Storage {
		if event.StartTime[8:10] == year {
			filtered = append(filtered, event)
		}
	}
	return filtered
}
