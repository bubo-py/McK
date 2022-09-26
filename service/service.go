package service

import (
	"errors"
	"github.com/bubo-py/McK/repositories"
	"github.com/bubo-py/McK/types"
)

type BusinessLogicInterface interface {
	GetEvents(id int)
	GetEvent(id int)
	AddEvent(e types.Event)
	DeleteEvent(id int)
	UpdateEvent(e types.Event, id int)
	ValidatePostRequest(e types.Event) error
}

type BusinessLogic struct{}

var db repositories.DatabaseRepository = repositories.InitDatabase()

func (bl BusinessLogic) GetEvents() []types.Event {
	return db.GetEvents()
}

func (bl BusinessLogic) GetEvent(id int) (types.Event, error) {
	return db.GetEvent(id)
}

func (bl BusinessLogic) AddEvent(e types.Event) error {
	err := bl.ValidatePostRequest(e)
	if err != nil {
		return err
	}

	db.AddEvent(e)
	return nil
}

func (bl BusinessLogic) DeleteEvent(id int) error {
	return db.DeleteEvent(id)
}

func (bl BusinessLogic) UpdateEvent(e types.Event, id int) error {
	return db.UpdateEvent(e, id)
}

func (bl BusinessLogic) ValidatePostRequest(e types.Event) error {
	if e.Name != "" && e.StartTime != "" && e.EndTime != "" {
		return nil
	}
	return errors.New("invalid post request")
}
