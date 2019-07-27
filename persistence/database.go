package persistence

import "github.com/jmoiron/sqlx"

// Repository wraps the underlying database implementation
type Repository struct {
	db *sqlx.DB
}

// New create a new instance of the database interaction logic
// by setting up the datbase
func New(dbdialect, connstr string) *Repository {
	db := sqlx.MustConnect(dbdialect, connstr)
	return &Repository{db: db}
}
