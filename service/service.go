package service

import (
	"errors"

	"github.com/bubo-py/McK/repositories"
	"github.com/bubo-py/McK/types"
)

type BusinessLogicInterface interface {
	GetEvents(f types.Filters) ([]types.Event, error)
	GetEvent(id int) (types.Event, error)
	AddEvent(e types.Event) error
	DeleteEvent(id int) error
	UpdateEvent(e types.Event, id int) error
}

type BusinessLogic struct {
	db repositories.DatabaseRepository
}

func InitBusinessLogic(db repositories.DatabaseRepository) BusinessLogic {
	var bl BusinessLogic
	bl.db = db
	return bl
}

func (bl BusinessLogic) GetEvents(f types.Filters) ([]types.Event, error) {
	var s []types.Event

	if f.Day == 0 && f.Month == 0 && f.Year == 0 {
		s = append(s, bl.db.GetEvents()...)
		return s, nil
	}

	if f.Day != 0 {
		if f.Day <= 0 || f.Day >= 32 {
			return s, errors.New("invalid day value")
		}

		s = append(s, bl.db.GetEventsByDay(f.Day)...)
	}

	if f.Month != 0 {
		if f.Month <= 0 || f.Month >= 13 {
			return s, errors.New("invalid month value")
		}

		s = append(s, bl.db.GetEventsByMonth(f.Month)...)
	}

	if f.Year != 0 {
		if f.Year <= 0 {
			return s, errors.New("invalid year value")
		}

		s = append(s, bl.db.GetEventsByYear(f.Year)...)
	}

	return s, nil
}

func (bl BusinessLogic) GetEvent(id int) (types.Event, error) {
	return bl.db.GetEvent(id)
}

func (bl BusinessLogic) AddEvent(e types.Event) error {
	err := validatePostRequest(e)
	if err != nil {
		return err
	}

	err = validateLength(e.Name)
	if err != nil {
		return err
	}

	err = bl.db.AddEvent(e)
	if err != nil {
		return err
	}

	return nil
}

func (bl BusinessLogic) DeleteEvent(id int) error {
	return bl.db.DeleteEvent(id)
}

func (bl BusinessLogic) UpdateEvent(e types.Event, id int) error {
	err := validateLength(e.Name)
	if err != nil {
		return err
	}

	return bl.db.UpdateEvent(e, id)
}

func validatePostRequest(e types.Event) error {
	if e.Name == "" || e.StartTime.IsZero() || e.EndTime.IsZero() {
		return errors.New("invalid post request")
	}

	return nil
}

func validateLength(s string) error {
	if len([]rune(s)) > 255 {
		return errors.New("length should be less than 255 characters")
	}

	return nil
}
