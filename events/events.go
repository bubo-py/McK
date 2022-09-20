package events

type Event struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	StartTime   string `json:"startTime"` // format: 2022-09-14T09:00:00.000Z
	EndTime     string `json:"endTime"`   // RFC 3339, section 5.6
	Description string `json:"description,omitempty"`
	AlertTime   string `json:"alertTime,omitempty"`
}

var ID int

var Db []Event

func CheckEvent(id int) (bool, int) {
	var index int
	present := false

	for i, event := range Db {
		if event.ID == id {
			index = i
			present = true
		}
	}
	return present, index
}

func AppendEvent(e Event) {
	ID += 1
	e.ID = ID
	Db = append(Db, e)
}

func DeleteEvent(id int) bool {
	present := false

	for i, event := range Db {
		if event.ID == id {
			copy(Db[i:], Db[i+1:])
			Db[len(Db)-1] = Event{}
			Db = Db[:len(Db)-1]
			present = true
		}
	}
	return present
}

func UpdateEvent(e Event, id int) bool {
	present := false

	for i, event := range Db {
		if event.ID == id {
			Db[i].Name = e.Name
			Db[i].StartTime = e.StartTime
			Db[i].EndTime = e.EndTime
			Db[i].Description = e.Description
			Db[i].AlertTime = e.AlertTime
			present = true
		}
	}
	return present
}
