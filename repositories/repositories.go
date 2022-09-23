package repositories

import (
	"github.com/bubo-py/McK/types"
)

type DatabaseRepository interface {
	CheckEvent(id int) (bool, int)
	GetEvents() []types.Event
	GetEventsPosition(id int) types.Event
	AppendEvent(e types.Event)
	DeleteEvent(id int) bool
	UpdateEvent(e types.Event, id int) bool
}

type DbHandler struct {
	DatabaseRepository
}
