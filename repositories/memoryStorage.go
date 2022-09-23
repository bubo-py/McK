package repositories

import (
	"github.com/bubo-py/McK/types"
)

type Database struct {
	ID      int
	Storage []types.Event
}

func InitDatabase() *Database {
	return &Database{}
}

func (db *Database) CheckEvent(id int) (bool, int) {
	var index int
	present := false

	for i, event := range db.Storage {
		if event.ID == id {
			index = i
			present = true
		}
	}
	return present, index
}

func (db *Database) GetEvents() []types.Event {
	return db.Storage
}

func (db *Database) GetEventsPosition(id int) types.Event {
	return db.Storage[id]
}

func (db *Database) AppendEvent(e types.Event) {
	db.ID += 1
	e.ID = db.ID
	db.Storage = append(db.Storage, e)
}

func (db *Database) DeleteEvent(id int) bool {
	present := false

	for i, event := range db.Storage {
		if event.ID == id {
			copy(db.Storage[i:], db.Storage[i+1:])
			db.Storage[len(db.Storage)-1] = types.Event{}
			db.Storage = db.Storage[:len(db.Storage)-1]
			present = true
		}
	}
	return present
}

func (db *Database) UpdateEvent(e types.Event, id int) bool {
	present := false

	for i, event := range db.Storage {
		if event.ID == id {
			db.Storage[i].Name = e.Name
			db.Storage[i].StartTime = e.StartTime
			db.Storage[i].EndTime = e.EndTime
			db.Storage[i].Description = e.Description
			db.Storage[i].AlertTime = e.AlertTime
			present = true
		}
	}
	return present
}
