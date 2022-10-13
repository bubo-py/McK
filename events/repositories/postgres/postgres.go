package postgres

import (
	"context"
	"embed"
	"errors"
	"log"
	"time"

	"github.com/bubo-py/McK/types"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/tern/migrate"
)

//go:embed migrations
var f embed.FS

type eventDb struct {
	ID          int64     `db:"id"`
	Name        string    `db:"name"`
	StartTime   time.Time `db:"starttime"` // format: 2022-09-14T09:00:00.000Z
	EndTime     time.Time `db:"endtime"`   // RFC 3339, section 5.6
	Description string    `db:"description,omitempty"`
	AlertTime   time.Time `db:"alerttime,omitempty"`
}

type Db struct {
	pool *pgxpool.Pool
}

func Init(ctx context.Context, connString string) (Db, error) {
	var pg Db

	dbPool, err := pgxpool.Connect(ctx, connString)
	if err != nil {
		return pg, err
	}

	//defer dbPool.Close()
	pg.pool = dbPool

	return pg, nil
}

func (pg Db) migrate(ctx context.Context, mFS embed.FS, rootDir, table string) error {
	c, err := pg.pool.Acquire(ctx)
	if err != nil {
		return err
	}
	defer c.Release()

	opts := &migrate.MigratorOptions{
		MigratorFS: adapterFS{
			FS:      mFS,
			rootDir: rootDir,
		},
	}
	migrator, err := migrate.NewMigratorEx(ctx, c.Conn(), table, opts)
	if err != nil {
		return err
	}
	if err := migrator.LoadMigrations(rootDir); err != nil {
		return err
	}
	if err := migrator.Migrate(ctx); err != nil {
		return err
	}

	_, err = migrator.GetCurrentVersion(ctx)
	if err != nil {
		return err
	}

	return nil
}

func RunMigration(ctx context.Context, db Db) error {
	err := db.migrate(ctx, f, "migrations", "events_migration")
	if err != nil {
		return err
	}

	log.Println("Migrations from events domain run correctly")
	return nil
}

func (pg Db) GetEvents(ctx context.Context) ([]types.Event, error) {
	var s []types.Event
	var events []*eventDb

	q := sqlbuilder.Select("*").From("events")

	err := pgxscan.Select(ctx, pg.pool, &events, q.String())
	if err != nil {
		return s, err
	}

	for _, event := range events {
		s = append(s, types.Event(*event))
	}

	return s, nil
}

func (pg Db) GetEvent(ctx context.Context, id int64) (types.Event, error) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	var e eventDb

	exists, err := pg.exists(ctx, id)
	if err != nil {
		return types.Event(e), err
	}

	if exists == true {
		sb.Select("id", "name", "startTime", "endTime", "description", "alertTime")
		sb.From("events")
		sb.Where(sb.Equal("id", id))

		q, args := sb.Build()

		err := pgxscan.Get(ctx, pg.pool, &e, q, args...)
		if err != nil {
			return types.Event(e), err
		}

		return types.Event(e), nil
	}

	return types.Event(e), errors.New("event with specified id not found")
}

func (pg Db) AddEvent(ctx context.Context, e types.Event) error {
	ib := sqlbuilder.PostgreSQL.NewInsertBuilder()

	ib.InsertInto("events")
	ib.Cols("name", "startTime", "endTime", "description", "alertTime")
	ib.Values(e.Name, e.StartTime, e.EndTime, e.Description, e.AlertTime)

	q, args := ib.Build()

	_, err := pg.pool.Exec(ctx, q, args...)
	if err != nil {
		return err
	}

	return nil
}

func (pg Db) DeleteEvent(ctx context.Context, id int64) error {
	db := sqlbuilder.PostgreSQL.NewDeleteBuilder()

	exists, err := pg.exists(ctx, id)
	if err != nil {
		return err
	}

	if exists == true {
		db.DeleteFrom("events")
		db.Where(db.Equal("id", id))

		q, args := db.Build()

		_, err = pg.pool.Exec(ctx, q, args...)
		if err != nil {
			return err
		}

		return nil
	}

	return errors.New("event with specified id not found")
}

func (pg Db) UpdateEvent(ctx context.Context, e types.Event, id int64) error {
	ub := sqlbuilder.PostgreSQL.NewUpdateBuilder()

	exists, err := pg.exists(ctx, id)
	if err != nil {
		return err
	}

	if exists == true {
		ub.Update("events")
		ub.Set(
			ub.Assign("name", e.Name),
			ub.Assign("startTime", e.StartTime),
			ub.Assign("endTime", e.EndTime),
			ub.Assign("description", e.Description),
			ub.Assign("alertTime", e.AlertTime),
		)

		ub.Where(ub.Equal("id", id))

		q, args := ub.Build()

		_, err := pg.pool.Exec(ctx, q, args...)
		if err != nil {
			return err
		}

		return nil
	}

	return errors.New("event with specified id not found")
}

func (pg Db) GetEventsFiltered(ctx context.Context, f types.Filters) ([]types.Event, error) {
	var filtered []types.Event
	var events []*eventDb

	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()

	sb.Select("id", "name", "startTime", "endTime", "description", "alertTime")
	sb.From("events")

	if f.Day != 0 {
		sb.Where(sb.Equal("EXTRACT(day FROM startTime)", f.Day))
	}

	if f.Month != 0 {
		sb.Where(sb.Equal("EXTRACT(month FROM startTime)", f.Month))
	}

	if f.Year != 0 {
		sb.Where(sb.Equal("EXTRACT(year FROM startTime)", f.Year))
	}

	q, args := sb.Build()

	err := pgxscan.Select(ctx, pg.pool, &events, q, args...)
	if err != nil {
		return filtered, err
	}

	for _, event := range events {
		filtered = append(filtered, types.Event(*event))
	}

	return filtered, nil
}

func (pg Db) exists(ctx context.Context, id int64) (bool, error) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	var exists bool

	sb.Select("EXISTS(select 1 from events)")
	sb.From("events")
	sb.Where(sb.Equal("id", id))

	q, args := sb.Build()

	rows, err := pg.pool.Query(ctx, q, args...)

	if err != nil {
		return exists, err
	}

	for rows.Next() {

		values, err := rows.Values()
		if err != nil {
			return exists, err
		}
		exists = values[0].(bool)
	}

	return exists, nil
}
