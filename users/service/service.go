package service

import (
	"context"
	"errors"

	"github.com/bubo-py/McK/types"
	"github.com/bubo-py/McK/users/repositories"
)

type BusinessLogicInterface interface {
	AddUser(ctx context.Context, u types.User) error
	UpdateUser(ctx context.Context, u types.User, id int64) error
	DeleteUser(ctx context.Context, id int64) error
}

type BusinessLogic struct {
	db repositories.UserRepository
}

func InitBusinessLogic(db repositories.UserRepository) BusinessLogic {
	var bl BusinessLogic
	bl.db = db
	return bl
}

func (bl BusinessLogic) AddUser(ctx context.Context, u types.User) error {
	err := validateName(u.Login)
	if err != nil {
		return err
	}

	return bl.db.AddUser(ctx, u)
}

func (bl BusinessLogic) UpdateUser(ctx context.Context, u types.User, id int64) error {
	err := validateName(u.Login)
	if err != nil {
		return err
	}

	return bl.db.UpdateUser(ctx, u, id)
}

func (bl BusinessLogic) DeleteUser(ctx context.Context, id int64) error {
	return bl.db.DeleteUser(ctx, id)
}

func validateName(s string) error {
	if len(s) > 30 {
		return errors.New("login should contain up to 30 characters")
	}
	return nil
}
