package persistence

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

// SenderReader search for tags in the store
type SenderReader interface {
	GetAllSenders() ([]Sender, error)
	SearchSenders(s string) ([]Sender, error)
}

type dbSenderReader struct {
	db *sqlx.DB
}

// NewSenderReader creates a new instance using an existing
// connection to a repository
func NewSenderReader(conn RepositoryConnection) (SenderReader, error) {
	if conn.h == nil {
		return nil, fmt.Errorf("no repository connection available")
	}
	return dbSenderReader{db: conn.h}, nil
}

// GetAllSenders returns all available senders in alphabetical order
func (r dbSenderReader) GetAllSenders() ([]Sender, error) {
	var senders []Sender
	if err := r.db.Select(&senders, "SELECT t.id, t.name FROM SENDERS t ORDER BY name ASC"); err != nil {
		return nil, err
	}
	return senders, nil
}

// SearchSenders returns senders matching the given search string
// the search string is matched independent of case and works in a wildcard fashion
func (r dbSenderReader) SearchSenders(s string) ([]Sender, error) {
	var senders []Sender
	search := strings.ToLower(s)
	search = "%" + search + "%"
	if err := r.db.Select(&senders, "SELECT t.id, t.name FROM SENDERS t WHERE t.name LIKE ? ORDER BY name ASC", search); err != nil {
		return nil, err
	}
	return senders, nil
}
