package service

import (
	"context"
	"errors"
	"log"

	"github.com/bubo-py/McK/contextHelpers"
	"github.com/bubo-py/McK/types"
	"github.com/bubo-py/McK/users/repositories"
	"golang.org/x/crypto/bcrypt"
)

type BusinessLogicInterface interface {
	AddUser(ctx context.Context, u types.User) (types.User, error)
	UpdateUser(ctx context.Context, u types.User, id int64) (types.User, error)
	DeleteUser(ctx context.Context, id int64) error
	LoginUser(ctx context.Context, login, password string) error
	GetUserByLogin(ctx context.Context, login string) (types.User, error)
}

type BusinessLogic struct {
	db repositories.UserRepository
}

func InitBusinessLogic(db repositories.UserRepository) BusinessLogic {
	var bl BusinessLogic
	bl.db = db
	return bl
}

func (bl BusinessLogic) AddUser(ctx context.Context, u types.User) (types.User, error) {
	err := validateLogin(u.Login)
	if err != nil {
		return u, err
	}

	u.Password, err = hashPassword(u.Password)
	if err != nil {
		return u, err
	}

	return bl.db.AddUser(ctx, u)
}

func (bl BusinessLogic) UpdateUser(ctx context.Context, u types.User, id int64) (types.User, error) {
	if u.Login != "" {
		err := validateLogin(u.Login)
		if err != nil {
			return u, err
		}
	}

	if u.Password != "" {
		hashedPwd, err := hashPassword(u.Password)
		if err != nil {
			return u, err
		}
		u.Password = hashedPwd
	}

	currentUserLogin := contextHelpers.RetrieveLoginFromContext(ctx)

	currentUser, _ := bl.db.GetUserByLogin(ctx, currentUserLogin)
	if currentUser.ID != id {
		return u, errors.New("cannot modify another user's account")
	}

	return bl.db.UpdateUser(ctx, u, id)
}

func (bl BusinessLogic) DeleteUser(ctx context.Context, id int64) error {
	currentUserLogin := contextHelpers.RetrieveLoginFromContext(ctx)

	currentUser, _ := bl.db.GetUserByLogin(ctx, currentUserLogin)
	if currentUser.ID != id {
		return errors.New("cannot delete another user's account")
	}

	return bl.db.DeleteUser(ctx, id)
}

func (bl BusinessLogic) GetUserByLogin(ctx context.Context, login string) (types.User, error) {
	return bl.db.GetUserByLogin(ctx, login)
}

func (bl BusinessLogic) LoginUser(ctx context.Context, login, password string) error {
	return bl.checkPassword(ctx, login, password)
}

func hashPassword(s string) (string, error) {
	if len(s) < 5 {
		return s, errors.New("password should be at least 5 characters")
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.DefaultCost)
	if err != nil {
		return s, err
	}

	s = string(bytes)
	if len(s) > 60 {
		return s, errors.New("failed to hash password")
	}

	return s, nil
}

func (bl BusinessLogic) checkPassword(ctx context.Context, login, password string) error {
	u, err := bl.db.GetUserByLogin(ctx, login)
	if err != nil {
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		log.Println("incorrect password")
		return errors.New("incorrect credentials")
	}

	return nil
}

func validateLogin(s string) error {
	if len([]rune(s)) > 30 || len([]rune(s)) < 3 {
		return errors.New("login should be at least 3 and contain up to 30 characters")
	}

	return nil
}
