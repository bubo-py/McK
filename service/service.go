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
	GetEventsByDay(day string) ([]types.Event, error)
	GetEventsByMonth(month string) ([]types.Event, error)
	GetEventsByYear(year string) ([]types.Event, error)
}

type BusinessLogic struct{}

var db repositories.DatabaseRepository = repositories.InitDatabase()

func (bl BusinessLogic) GetEvents(f types.Filters) []types.Event {
	var s []types.Event

	if f.Day == 0 && f.Month == 0 && f.Year == 0 {
		s = append(s, db.GetEvents()...)
		return s
	}

	s = append(s, db.GetEventsByDay(f.Day)...)
	s = append(s, db.GetEventsByMonth(f.Month)...)
	s = append(s, db.GetEventsByYear(f.Year)...)

	return s
}

func (bl BusinessLogic) GetEvent(id int) (types.Event, error) {
	return db.GetEvent(id)
}

func (bl BusinessLogic) AddEvent(e types.Event) error {
	err := validatePostRequest(e)

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

func validatePostRequest(e types.Event) error {
	if e.Name == "" || e.StartTime.IsZero() || e.EndTime.IsZero() {
		return errors.New("invalid post request")
	}

	return nil
}

func (bl BusinessLogic) GetEventsByDay(day int) ([]types.Event, error) {
	if day <= 0 && day >= 32 {
		return []types.Event{}, errors.New("invalid day value")
	}

	return db.GetEventsByDay(day), nil
}

func (bl BusinessLogic) GetEventsByMonth(month int) ([]types.Event, error) {
	if month <= 0 && month >= 13 {
		return []types.Event{}, errors.New("invalid month value")
	}

	return db.GetEventsByMonth(month), nil
}

func (bl BusinessLogic) GetEventsByYear(year int) ([]types.Event, error) {
	if year <= 0 {
		return []types.Event{}, errors.New("invalid year value")
	}

	return db.GetEventsByYear(year), nil
}
