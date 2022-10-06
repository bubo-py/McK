package repositories

import (
	"context"
	"errors"

	"github.com/bubo-py/McK/types"
)

type Database struct {
	ID      int64
	Storage []types.Event
}

func InitDatabase() *Database {
	return &Database{}
}

func (db *Database) GetEvents(ctx context.Context) ([]types.Event, error) {
	return db.Storage, nil
}

func (db *Database) GetEvent(ctx context.Context, id int64) (types.Event, error) {
	for i, event := range db.Storage {
		if event.ID == id {
			return db.Storage[i], nil
		}
	}
	return types.Event{}, errors.New("event with specified id not found")
}

func (db *Database) AddEvent(ctx context.Context, e types.Event) error {
	db.ID += 1
	e.ID = db.ID
	db.Storage = append(db.Storage, e)

	return nil
}

func (db *Database) DeleteEvent(ctx context.Context, id int64) error {
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

func (db *Database) UpdateEvent(ctx context.Context, e types.Event, id int64) error {
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

func (db *Database) GetEventsByDay(ctx context.Context, day int) ([]types.Event, error) {
	filtered := make([]types.Event, 0)

	for _, event := range db.Storage {
		if event.StartTime.Day() == day {
			filtered = append(filtered, event)
		}
	}
	return filtered, nil
}

func (db *Database) GetEventsByMonth(ctx context.Context, month int) ([]types.Event, error) {
	filtered := make([]types.Event, 0)

	for _, event := range db.Storage {
		if int(event.StartTime.Month()) == month {
			filtered = append(filtered, event)
		}
	}
	return filtered, nil
}

func (db *Database) GetEventsByYear(ctx context.Context, year int) ([]types.Event, error) {
	filtered := make([]types.Event, 0)

	for _, event := range db.Storage {
		if event.StartTime.Year() == year {
			filtered = append(filtered, event)
		}
	}
	return filtered, nil
}
