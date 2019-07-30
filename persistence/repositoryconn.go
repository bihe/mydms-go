package persistence

import "github.com/jmoiron/sqlx"

// RepositoryConnection represents an arbitrary connection to a store/repository
type RepositoryConnection struct {
	// handle to database connection
	h *sqlx.DB
}

// NewConnection creates a connection to a store/repository
func NewConnection(dbdialect, connstr string) RepositoryConnection {
	db := sqlx.MustConnect(dbdialect, connstr)
	return RepositoryConnection{h: db}
}
