package repositories

import (
	"context"

	"github.com/bubo-py/McK/types"
)

type UserRepository interface {
	AddUser(ctx context.Context, u types.User) (types.User, error)
	UpdateUser(ctx context.Context, u types.User, id int64) (types.User, error)
	DeleteUser(ctx context.Context, id int64) error
	GetUserByLogin(ctx context.Context, login string) (types.User, error)
}
