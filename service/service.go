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
}

type BusinessLogic struct{}

var db repositories.DatabaseRepository = repositories.InitDatabase()

func (bl BusinessLogic) GetEvents(f types.Filters) ([]types.Event, error) {
	var s []types.Event

	if f.Day == 0 && f.Month == 0 && f.Year == 0 {
		s = append(s, db.GetEvents()...)
		return s, nil
	}

	if f.Day != 0 {
		if f.Day <= 0 || f.Day >= 32 {
			return s, errors.New("invalid day value")
		}

		s = append(s, db.GetEventsByDay(f.Day)...)
	}

	if f.Month != 0 {
		if f.Month <= 0 || f.Month >= 13 {
			return s, errors.New("invalid month value")
		}

		s = append(s, db.GetEventsByMonth(f.Month)...)
	}

	if f.Year != 0 {
		if f.Year <= 0 {
			return s, errors.New("invalid year value")
		}

		s = append(s, db.GetEventsByYear(f.Year)...)
	}

	return s, nil
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
