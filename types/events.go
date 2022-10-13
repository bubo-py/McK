package types

import "time"

type Event struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	StartTime   time.Time `json:"startTime"` // format: 2022-09-14T09:00:00.000Z
	EndTime     time.Time `json:"endTime"`   // RFC 3339, section 5.6
	Description string    `json:"description,omitempty"`
	AlertTime   time.Time `json:"alertTime,omitempty"`
}
