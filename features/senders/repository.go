package senders

import (
	"fmt"
	"strings"

	"github.com/bihe/mydms/persistence"
)

// SenderEntity is the originator of a document
type SenderEntity struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

// Reader search for tags in the store
type Reader interface {
	GetAllSenders() ([]SenderEntity, error)
	SearchSenders(s string) ([]SenderEntity, error)
}

type dbSenderReader struct {
	c persistence.Connection
}

// NewReader creates a new instance using an existing
// connection to a repository
func NewReader(c persistence.Connection) (Reader, error) {
	if !c.Active {
		return nil, fmt.Errorf("no repository connection available")
	}
	return dbSenderReader{c}, nil
}

// GetAllSenders returns all available senders in alphabetical order
func (r dbSenderReader) GetAllSenders() ([]SenderEntity, error) {
	var senders []SenderEntity
	if err := r.c.Select(&senders, "SELECT t.id, t.name FROM SENDERS t ORDER BY name ASC"); err != nil {
		return nil, err
	}
	return senders, nil
}

// SearchSenders returns senders matching the given search string
// the search string is matched independent of case and works in a wildcard fashion
func (r dbSenderReader) SearchSenders(s string) ([]SenderEntity, error) {
	var senders []SenderEntity
	search := strings.ToLower(s)
	search = "%" + search + "%"
	if err := r.c.Select(&senders, "SELECT t.id, t.name FROM SENDERS t WHERE t.name LIKE ? ORDER BY name ASC", search); err != nil {
		return nil, err
	}
	return senders, nil
}
