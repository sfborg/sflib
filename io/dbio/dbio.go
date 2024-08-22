package dbio

import (
	"database/sql"

	"github.com/sfborg/sflib/ent/sfga"
)

type dbio struct {
	dir  string
	file string
	db   *sql.DB
}

func New(dbDir string) sfga.DB {
	return &dbio{dir: dbDir}
}

func (d *dbio) Connect() (*sql.DB, error) {
	return nil, nil
}

func (d *dbio) Close() error {
	return nil
}

func (d *dbio) FileDB() string {
	return d.file
}
