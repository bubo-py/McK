package serviceDb

import (
	"context"

	"github.com/bubo-py/McK/types"
)

type Db struct{}

func (db Db) AddUser(ctx context.Context, u types.User) (types.User, error) {
	return u, nil
}

func (db Db) UpdateUser(ctx context.Context, u types.User, id int64) (types.User, error) {
	return u, nil
}

func (db Db) DeleteUser(ctx context.Context, id int64) error {
	return nil
}

func (db Db) GetUserByLogin(ctx context.Context, login string) (types.User, error) {
	var u types.User
	return u, nil
}
