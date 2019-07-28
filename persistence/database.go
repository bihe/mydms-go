package persistence

import "github.com/jmoiron/sqlx"

// Repository gives access to all persistence methods and interacts with the given store
type Repository interface {
	GetAllTags() ([]Tag, error)
	SearchTags(s string) ([]Tag, error)
}

// DBRepository wraps the underlying database implementation
type repository struct {
	db *sqlx.DB
}

// New create a new instance of the database interaction logic
// by setting up the datbase
func New(dbdialect, connstr string) Repository {
	db := sqlx.MustConnect(dbdialect, connstr)
	return repository{db: db}
}
