package repositories

import (
	"context"
	"embed"
	"log"
	"time"

	"github.com/bubo-py/McK/types"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/tern/migrate"
)

//go:embed migrations
var f embed.FS

type PostgresDb struct {
	pool *pgxpool.Pool
}

func PostgresInit(ctx context.Context) PostgresDb {
	var pg PostgresDb
	dbUrl := "postgres://postgres:12345@localhost:5432/postgres"

	dbPool, err := pgxpool.Connect(ctx, dbUrl)
	if err != nil {
		log.Println(err)
	}

	//defer dbPool.Close()
	pg.pool = dbPool

	return pg
}

func (pg PostgresDb) migrate(ctx context.Context, mFS embed.FS, rootDir, table string) error {
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

func RunMigration(ctx context.Context, db PostgresDb) error {
	err := db.migrate(ctx, f, "migrations", "migration")
	if err != nil {
		return err
	}

	return nil
}

func (pg PostgresDb) GetEvents(ctx context.Context) ([]types.Event, error) {
	var s []types.Event

	q := sqlbuilder.Select("*").From("events")

	rows, err := pg.pool.Query(ctx, q.String())
	if err != nil {
		return s, err
	}

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return s, err
		}
		id := values[0].(int64)
		name := values[1].(string)
		startTime := values[2].(time.Time)
		endTime := values[3].(time.Time)
		description := values[4].(string)
		alertTime := values[5].(time.Time)

		e := types.Event{
			ID:          id,
			Name:        name,
			StartTime:   startTime,
			EndTime:     endTime,
			Description: description,
			AlertTime:   alertTime,
		}
		s = append(s, e)
	}

	return s, nil
}

func (pg PostgresDb) GetEvent(ctx context.Context, id int64) (types.Event, error) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	var e types.Event

	sb.Select("id", "name", "startTime", "endTime", "description", "alertTime")
	sb.From("events")
	sb.Where(sb.Equal("id", id))

	q, args := sb.Build()

	rows, err := pg.pool.Query(ctx, q, args[0])
	if err != nil {
		return e, err
	}

	for rows.Next() {

		values, err := rows.Values()
		if err != nil {
			return e, err
		}

		id := values[0].(int64)
		name := values[1].(string)
		startTime := values[2].(time.Time)
		endTime := values[3].(time.Time)
		description := values[4].(string)
		alertTime := values[5].(time.Time)

		e.ID = id
		e.Name = name
		e.StartTime = startTime
		e.EndTime = endTime
		e.Description = description
		e.AlertTime = alertTime

		return e, nil
	}

	return e, nil
}

func (pg PostgresDb) AddEvent(ctx context.Context, e types.Event) error {
	ib := sqlbuilder.PostgreSQL.NewInsertBuilder()

	ib.InsertInto("events")
	ib.Cols("name", "startTime", "endTime", "description", "alertTime")
	ib.Values(e.Name, e.StartTime, e.EndTime, e.Description, e.AlertTime)

	q, args := ib.Build()

	_, err := pg.pool.Exec(ctx, q, args[0], args[1], args[2], args[3], args[4])
	if err != nil {
		return err
	}

	return nil
}

func (pg PostgresDb) DeleteEvent(ctx context.Context, id int64) error {
	db := sqlbuilder.PostgreSQL.NewDeleteBuilder()

	db.DeleteFrom("events")
	db.Where(db.Equal("id", id))

	q, args := db.Build()

	_, err := pg.pool.Exec(ctx, q, args[0])
	if err != nil {
		return err
	}

	return nil
}

func (pg PostgresDb) UpdateEvent(ctx context.Context, e types.Event, id int64) error {
	ub := sqlbuilder.PostgreSQL.NewUpdateBuilder()

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

	_, err := pg.pool.Exec(ctx, q, args[0], args[1], args[2], args[3], args[4], args[5])
	if err != nil {
		return err
	}

	return nil
}

func (pg PostgresDb) GetEventsByDay(ctx context.Context, day int) ([]types.Event, error) {
	filtered := make([]types.Event, 0)

	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()

	sb.Select("id", "name", "startTime", "endTime", "description", "alertTime")
	sb.From("events")
	sb.Where(sb.Equal("EXTRACT(day FROM startTime)", day))

	q, args := sb.Build()

	rows, err := pg.pool.Query(ctx, q, args[0])
	if err != nil {
		return filtered, err
	}

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return filtered, err
		}
		id := values[0].(int64)
		name := values[1].(string)
		startTime := values[2].(time.Time)
		endTime := values[3].(time.Time)
		description := values[4].(string)
		alertTime := values[5].(time.Time)

		e := types.Event{
			ID:          id,
			Name:        name,
			StartTime:   startTime,
			EndTime:     endTime,
			Description: description,
			AlertTime:   alertTime,
		}
		filtered = append(filtered, e)
	}

	return filtered, nil
}

func (pg PostgresDb) GetEventsByMonth(ctx context.Context, month int) ([]types.Event, error) {
	filtered := make([]types.Event, 0)

	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()

	sb.Select("id", "name", "startTime", "endTime", "description", "alertTime")
	sb.From("events")
	sb.Where(sb.Equal("EXTRACT(month FROM startTime)", month))

	q, args := sb.Build()

	rows, err := pg.pool.Query(ctx, q, args[0])
	if err != nil {
		return filtered, err
	}

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return filtered, err
		}
		id := values[0].(int64)
		name := values[1].(string)
		startTime := values[2].(time.Time)
		endTime := values[3].(time.Time)
		description := values[4].(string)
		alertTime := values[5].(time.Time)

		e := types.Event{
			ID:          id,
			Name:        name,
			StartTime:   startTime,
			EndTime:     endTime,
			Description: description,
			AlertTime:   alertTime,
		}
		filtered = append(filtered, e)
	}

	return filtered, nil
}

func (pg PostgresDb) GetEventsByYear(ctx context.Context, year int) ([]types.Event, error) {
	filtered := make([]types.Event, 0)

	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()

	sb.Select("id", "name", "startTime", "endTime", "description", "alertTime")
	sb.From("events")
	sb.Where(sb.Equal("EXTRACT(year FROM startTime)", year))

	q, args := sb.Build()

	rows, err := pg.pool.Query(ctx, q, args[0])
	if err != nil {
		return filtered, err
	}

	for rows.Next() {
		values, err := rows.Values()
		if err != nil {
			return filtered, err
		}
		id := values[0].(int64)
		name := values[1].(string)
		startTime := values[2].(time.Time)
		endTime := values[3].(time.Time)
		description := values[4].(string)
		alertTime := values[5].(time.Time)

		e := types.Event{
			ID:          id,
			Name:        name,
			StartTime:   startTime,
			EndTime:     endTime,
			Description: description,
			AlertTime:   alertTime,
		}
		filtered = append(filtered, e)
	}

	return filtered, nil
}
