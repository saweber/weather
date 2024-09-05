package main

import (
	"context"
	"database/sql"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type Repository struct {
	db *bun.DB
}

func (r *Repository) Init() {
	dsn := "postgres://user:password@localhost:5432?sslmode=disable"
	pgdb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	r.db = bun.NewDB(pgdb, pgdialect.New())
}

func (r *Repository) InsertRecordIntoPostgres(report StormReport) error {
	_, err := r.db.NewInsert().
		Model(&report).
		On("CONFLICT (report_date, storm_type, location, county, state) DO UPDATE").
		Exec(context.Background())
	if err != nil {
		return err
	}
	return nil
}
