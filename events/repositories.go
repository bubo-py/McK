package events

type DatabaseRepository interface {
	CheckEvent(id int) (bool, int)
	GetStorage() []Event
	GetStoragePosition(id int) Event
	AppendEvent(e Event)
	DeleteEvent(id int) bool
	UpdateEvent(e Event, id int) bool
}
