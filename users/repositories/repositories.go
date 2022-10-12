package repositories

import (
	"context"

	"github.com/bubo-py/McK/types"
)

type UserRepository interface {
	AddUser(ctx context.Context, u types.User) error
	UpdateUser(ctx context.Context, u types.User, id int64) error
	DeleteUser(ctx context.Context, id int64) error
}
