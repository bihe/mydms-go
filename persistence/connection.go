package persistence

import "github.com/jmoiron/sqlx"

// Connection defines a storage/database/... connection
type Connection struct {
	*sqlx.DB
	Active bool
}

// NewConn creates a connection to a store/repository
func NewConn(connstr string) Connection {
	// dialect is specifically set to "mysql"
	db := sqlx.MustConnect("mysql", connstr)
	return Connection{DB: db, Active: true}
}

// NewFromDB creates a new connection based on existing DB handle
func NewFromDB(db *sqlx.DB) Connection {
	return Connection{DB: db, Active: true}
}

// Atomic defines a transactional operation - like the A in ACID https://en.wikipedia.org/wiki/ACID
type Atomic struct {
	*sqlx.Tx
	Active bool
}

// CreateAtomic starts a new transaction
func (c Connection) CreateAtomic() (Atomic, error) {
	tx, err := c.Beginx()
	if err != nil {
		return Atomic{}, err
	}
	return Atomic{Tx: tx, Active: true}, nil
}
