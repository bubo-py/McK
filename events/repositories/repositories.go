package repositories

import (
	"context"

	"github.com/bubo-py/McK/types"
)

type DatabaseRepository interface {
	GetEvents(ctx context.Context) ([]types.Event, error)
	GetEventsFiltered(ctx context.Context, f types.Filters) ([]types.Event, error)
	GetEvent(ctx context.Context, id int64) (types.Event, error)
	AddEvent(ctx context.Context, e types.Event) error
	DeleteEvent(ctx context.Context, id int64) error
	UpdateEvent(ctx context.Context, e types.Event, id int64) error
}
