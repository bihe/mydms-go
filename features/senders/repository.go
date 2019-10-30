package senders

import (
	"fmt"
	"strings"

	"github.com/bihe/mydms/internal/persistence"
)

// SenderEntity is the originator of a document
type SenderEntity struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

// Repository search for senders and saves senders in the store
type Repository interface {
	// GetAllSenders returns all available sender entries
	GetAllSenders() ([]SenderEntity, error)
	// SearchSenders finds all entries using the supplied search-term
	SearchSenders(s string) ([]SenderEntity, error)
	// SaveSenders processes the list of sender-names and stores the entries if not available
	SaveSenders(senders []string, a persistence.Atomic) (err error)
	// CreateSender creates a sender with the given name or returns an existing one
	CreateSender(name string, a persistence.Atomic) (sender SenderEntity, err error)
	// GetSenderByName returns the entity defined by the given name
	GetSenderByName(name string) (SenderEntity, error)
}

type dbRepository struct {
	c persistence.Connection
}

// NewRepository creates a new instance using an existing
// connection to a repository
func NewRepository(c persistence.Connection) (Repository, error) {
	if !c.Active {
		return nil, fmt.Errorf("no repository connection available")
	}
	return &dbRepository{c}, nil
}

// GetAllSenders returns all available senders in alphabetical order
func (r *dbRepository) GetAllSenders() ([]SenderEntity, error) {
	var senders []SenderEntity
	if err := r.c.Select(&senders, "SELECT t.id, t.name FROM SENDERS t ORDER BY name ASC"); err != nil {
		return nil, err
	}
	return senders, nil
}

// SearchSenders returns senders matching the given search string
// the search string is matched independent of case and works in a wildcard fashion
func (r *dbRepository) SearchSenders(s string) ([]SenderEntity, error) {
	var senders []SenderEntity
	search := strings.ToLower(s)
	search = "%" + search + "%"
	if err := r.c.Select(&senders, "SELECT t.id, t.name FROM SENDERS t WHERE lower(t.name) LIKE ? ORDER BY name ASC", search); err != nil {
		return nil, err
	}
	return senders, nil
}

// GetSenderByName returns the sender by given name
// the search is performed ignoring case sensitivity
func (r *dbRepository) GetSenderByName(name string) (SenderEntity, error) {
	var sender SenderEntity
	search := strings.ToLower(name)
	if err := r.c.Get(&sender, "SELECT t.id, t.name FROM SENDERS t WHERE lower(t.name) = ?", search); err != nil {
		return SenderEntity{}, err
	}
	return sender, nil
}

// SaveSenders takes a slice of strings and saves sender entries if they do not exist
// the existance-check is done by comparing the sender-name
func (r *dbRepository) SaveSenders(senders []string, a persistence.Atomic) (err error) {
	var atomic *persistence.Atomic

	defer func() {
		err = persistence.HandleTX(!a.Active, atomic, err)
	}()

	if atomic, err = persistence.CheckTX(r.c, &a); err != nil {
		return
	}

	var c int
	for _, s := range senders {
		s = strings.ToLower(s)
		err = atomic.Get(&c, "SELECT count(s.id) FROM SENDERS s WHERE lower(s.name) = ?", s)
		if err != nil {
			err = fmt.Errorf("could not search for a sender: %v", err)
			return
		}
		if c > 0 {
			continue
		}

		_, err = atomic.Exec("INSERT INTO SENDERS (name) VALUES (?)", s)
		if err != nil {
			err = fmt.Errorf("cannot save sender item: %v", err)
			return
		}
	}
	return
}

// CreateSender creates a sender with the given name or returns an existing one
func (r *dbRepository) CreateSender(name string, a persistence.Atomic) (sender SenderEntity, err error) {
	var atomic *persistence.Atomic

	defer func() {
		err = persistence.HandleTX(!a.Active, atomic, err)
	}()

	if atomic, err = persistence.CheckTX(r.c, &a); err != nil {
		return
	}

	err = atomic.Get(&sender, "SELECT id,name FROM SENDERS t WHERE lower(t.name) = ?", strings.ToLower(name))
	if err == nil {
		return sender, nil
	}

	res, err := atomic.Exec("INSERT INTO SENDERS (name) VALUES (?)", name)
	if err != nil {
		err = fmt.Errorf("cannot save sender item, %v", err)
		return
	}
	id, err := res.LastInsertId()
	if err != nil {
		err = fmt.Errorf("could not get last inserted id, %v", err)
		return
	}
	return SenderEntity{ID: int(id), Name: name}, nil
}
