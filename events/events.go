package events

type Event struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	StartTime   string `json:"startTime"` // format: 2022-09-14T09:00:00.000Z
	EndTime     string `json:"endTime"`   // RFC 3339, section 5.6
	Description string `json:"description,omitempty"`
	AlertTime   string `json:"alertTime,omitempty"`
}

type Database struct {
	ID      int
	Storage []Event
}

func InitDatabase() Database {
	return Database{}
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

func (db *Database) AppendEvent(e Event) {
	db.ID += 1
	e.ID = db.ID
	db.Storage = append(db.Storage, e)
}

func (db *Database) DeleteEvent(id int) bool {
	present := false

	for i, event := range db.Storage {
		if event.ID == id {
			copy(db.Storage[i:], db.Storage[i+1:])
			db.Storage[len(db.Storage)-1] = Event{}
			db.Storage = db.Storage[:len(db.Storage)-1]
			present = true
		}
	}
	return present
}

func (db *Database) UpdateEvent(e Event, id int) bool {
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
