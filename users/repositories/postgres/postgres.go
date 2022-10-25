package postgres

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"log"

	"github.com/bubo-py/McK/customErrors"
	"github.com/bubo-py/McK/types"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/tern/migrate"
)

var errUnexpected = customErrors.ErrUnexpected

//go:embed migrations
var f embed.FS

type Db struct {
	pool *pgxpool.Pool
}

func Init(ctx context.Context, connString string) (Db, error) {
	var pg Db

	dbPool, err := pgxpool.Connect(ctx, connString)
	if err != nil {
		errUnexpected.Err = err
		return pg, fmt.Errorf("database initialization error: %w", errUnexpected)
	}

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
	err := db.migrate(ctx, f, "migrations", "users_migration")
	if err != nil {
		errUnexpected.Err = err
		return fmt.Errorf("database migration error: %w", errUnexpected)
	}

	log.Println("Migrations from users domain run correctly")
	return nil
}

func (pg Db) AddUser(ctx context.Context, u types.User) (types.User, error) {
	ib := sqlbuilder.PostgreSQL.NewInsertBuilder()

	ib.InsertInto("users")
	ib.Cols("login", "password", "timezone")
	ib.Values(u.Login, u.Password, u.Timezone)

	ib.SQL("RETURNING *")

	q, args := ib.Build()

	err := pgxscan.Get(ctx, pg.pool, &u, q, args...)
	if err != nil {
		errUnexpected.Err = err
		return u, fmt.Errorf("database error: SQL query error: %w", errUnexpected)
	}

	return u, nil
}

func (pg Db) UpdateUser(ctx context.Context, u types.User, id int64) (types.User, error) {
	ub := sqlbuilder.PostgreSQL.NewUpdateBuilder()

	exists, err := pg.exists(ctx, id)
	if err != nil {
		return u, err
	}

	if exists == false {
		return u, errors.New("user with specified id not found")
	}

	ub.Update("users")

	if u.Login != "" {
		ub.SetMore(ub.Assign("login", u.Login))
	}

	if u.Password != "" {
		ub.SetMore(ub.Assign("password", u.Password))
	}

	if u.Timezone != "" {
		ub.SetMore(ub.Assign("timezone", u.Timezone))
	}

	ub.Where(ub.Equal("id", id))
	ub.SQL("RETURNING *")

	q, args := ub.Build()

	_, err = pg.pool.Exec(ctx, q, args...)
	if err != nil {
		errUnexpected.Err = err
		return u, fmt.Errorf("database error: SQL query error: %w", errUnexpected)
	}

	return u, nil
}

func (pg Db) DeleteUser(ctx context.Context, id int64) error {
	db := sqlbuilder.PostgreSQL.NewDeleteBuilder()

	exists, err := pg.exists(ctx, id)
	if err != nil {
		return err
	}

	if exists == false {
		return errors.New("user with specified id not found")
	}

	db.DeleteFrom("users")
	db.Where(db.Equal("id", id))

	q, args := db.Build()

	_, err = pg.pool.Exec(ctx, q, args...)
	if err != nil {
		errUnexpected.Err = err
		return fmt.Errorf("database error: SQL query error: %w", errUnexpected)
	}

	return nil
}

func (pg Db) GetUserByLogin(ctx context.Context, login string) (types.User, error) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	var u types.User

	sb.Select("id", "login", "password", "timezone")
	sb.From("users")
	sb.Where(sb.Equal("login", login))

	q, args := sb.Build()

	err := pgxscan.Get(ctx, pg.pool, &u, q, args...)
	if err != nil {
		log.Println("user with provided login not found")
		return u, customErrors.IncorrectCredentials
	}

	return u, nil
}

func (pg Db) exists(ctx context.Context, id int64) (bool, error) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	var exists bool

	sb.Select("EXISTS(select 1 from users)")
	sb.From("users")
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
