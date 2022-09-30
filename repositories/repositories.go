package repositories

import (
	"github.com/bubo-py/McK/types"
)

type DatabaseRepository interface {
	GetEvents() []types.Event
	GetEvent(id int) (types.Event, error)
	AddEvent(e types.Event)
	DeleteEvent(id int) error
	UpdateEvent(e types.Event, id int) error
	GetEventsByDay(day int) []types.Event
	GetEventsByMonth(month int) []types.Event
	GetEventsByYear(year int) []types.Event
}
