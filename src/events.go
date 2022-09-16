package main

type Event struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	StartTime   string `json:"startTime"` // format: 2022-09-14T09:00:00.000Z
	EndTime     string `json:"endTime"`   // RFC 3339, section 5.6
	Description string `json:"description,omitempty"`
	AlertTime   string `json:"alertTime,omitempty"`
}

func AppendEvent(e Event) {
	db = append(db, e)
}
