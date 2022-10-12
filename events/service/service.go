package service

import (
	"context"
	"errors"

	"github.com/bubo-py/McK/events/repositories"
	"github.com/bubo-py/McK/types"
)

type BusinessLogicInterface interface {
	GetEvents(ctx context.Context, f types.Filters) ([]types.Event, error)
	GetEvent(ctx context.Context, id int64) (types.Event, error)
	AddEvent(ctx context.Context, e types.Event) error
	DeleteEvent(ctx context.Context, id int64) error
	UpdateEvent(ctx context.Context, e types.Event, id int64) error
}

type BusinessLogic struct {
	db repositories.DatabaseRepository
}

func InitBusinessLogic(db repositories.DatabaseRepository) BusinessLogic {
	var bl BusinessLogic
	bl.db = db
	return bl
}

func (bl BusinessLogic) GetEvents(ctx context.Context, f types.Filters) ([]types.Event, error) {
	var s []types.Event

	if f.Day == 0 && f.Month == 0 && f.Year == 0 {
		e, err := bl.db.GetEvents(ctx)
		if err != nil {
			return s, err
		}

		s = append(s, e...)
		return s, nil
	}

	if f.Day != 0 {
		if f.Day <= 0 || f.Day >= 32 {
			return s, errors.New("invalid day value")
		}
	}

	if f.Month != 0 {
		if f.Month <= 0 || f.Month >= 13 {
			return s, errors.New("invalid month value")
		}
	}

	if f.Year != 0 {
		if f.Year <= 0 {
			return s, errors.New("invalid year value")
		}
	}

	e, err := bl.db.GetEventsFiltered(ctx, f)
	if err != nil {
		return s, err
	}

	s = append(s, e...)

	return s, nil
}

func (bl BusinessLogic) GetEvent(ctx context.Context, id int64) (types.Event, error) {
	return bl.db.GetEvent(ctx, id)
}

func (bl BusinessLogic) AddEvent(ctx context.Context, e types.Event) error {
	err := validatePostRequest(e)
	if err != nil {
		return err
	}

	err = validateLength(e.Name)
	if err != nil {
		return err
	}

	err = bl.db.AddEvent(ctx, e)
	if err != nil {
		return err
	}

	return nil
}

func (bl BusinessLogic) DeleteEvent(ctx context.Context, id int64) error {
	return bl.db.DeleteEvent(ctx, id)
}

func (bl BusinessLogic) UpdateEvent(ctx context.Context, e types.Event, id int64) error {
	err := validateLength(e.Name)
	if err != nil {
		return err
	}

	return bl.db.UpdateEvent(ctx, e, id)
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
