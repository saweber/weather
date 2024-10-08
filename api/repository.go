package main

import (
	"context"
	"database/sql"
	"time"

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

type StormReport struct {
	bun.BaseModel `bun:"table:storm_reports,alias:sr"`
	ReportDate    string    `json:"reportDate,omitempty" bun:"report_date"`
	StormType     string    `json:"stormType,omitempty" bun:"storm_type"`
	Latitude      string    `json:"latitude,omitempty"`
	Longitude     string    `json:"longitude,omitempty"`
	Location      string    `json:"location,omitempty"`
	County        string    `json:"county,omitempty"`
	State         string    `json:"state,omitempty"`
	Comments      string    `json:"comments,omitempty"`
	Speed         int       `json:"speed,omitempty"`
	Size          int       `json:"size,omitempty"`
	FScale        int       `json:"fScale,omitempty" bun:"f_scale"`
	Time          time.Time `json:"time,omitempty" bun:"time"`
}

type GetStormReportsOptions struct {
	ReportDateStart *time.Time `json:"reportDateStart"`
	ReportDateEnd   *time.Time `json:"reportDateEnd"`
	Location        string     `json:"location,omitempty"`
}

func (r *Repository) GetStormReports(opts GetStormReportsOptions) *[]StormReport {
	query := r.db.NewSelect().TableExpr("storm_reports")

	if opts.Location != "" {
		query = query.Where("location LIKE ?", "%"+opts.Location+"%")
	}
	if opts.ReportDateStart != nil {
		query = query.Where("time >= ?::timestamp", opts.ReportDateStart.Format(time.RFC3339))
	}
	if opts.ReportDateEnd != nil {
		query = query.Where("time < ?::timestamp", opts.ReportDateEnd.Format(time.RFC3339))
	}

	reports := []StormReport{}
	err := query.Scan(context.Background(), &reports)
	if err != nil {
		return nil
	}
	return &reports
}
