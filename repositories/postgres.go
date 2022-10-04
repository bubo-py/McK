package repositories

import (
	"context"
	"embed"
	"fmt"
	"log"

	"github.com/bubo-py/McK/types"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/tern/migrate"
)

// go:embed ../cmd/001_initMigration.sql
// var F embed.FS
var Path string = "../repositories/migrations/001_initMigration.sql"

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

func (pg PostgresDb) Migrate(ctx context.Context, mFS embed.FS, rootDir, table string) error {
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

func (pg PostgresDb) GetEvents() []types.Event {
	q := sqlbuilder.Select("*").From("events")

	fmt.Println(q)
	return nil
}
