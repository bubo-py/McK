package service

import (
	"errors"
	"strconv"

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
	GetEventsByDay(day string) ([]types.Event, error)
	GetEventsByMonth(month string) ([]types.Event, error)
	GetEventsByYear(year string) ([]types.Event, error)
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

func (bl BusinessLogic) GetEventsByDay(day string) ([]types.Event, error) {
	d, err := strconv.Atoi(day)
	if err != nil {
		return []types.Event{}, err
	}

	if d > 0 && d < 32 {
		return db.GetEventsByDay(day), nil
	}

	return []types.Event{}, errors.New("invalid day value")
}

func (bl BusinessLogic) GetEventsByMonth(month string) ([]types.Event, error) {
	m, err := strconv.Atoi(month)
	if err != nil {
		return []types.Event{}, err
	}

	if m > 0 && m < 13 {
		return db.GetEventsByMonth(month), nil
	}

	return []types.Event{}, errors.New("invalid month value")
}

func (bl BusinessLogic) GetEventsByYear(year string) ([]types.Event, error) {
	y, err := strconv.Atoi(year)
	if err != nil {
		return []types.Event{}, err
	}

	if y > 0 {
		return db.GetEventsByYear(year), nil
	}

	return []types.Event{}, errors.New("invalid year value")
}
