package persistence

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

// TagReader search for tags in the store
type TagReader interface {
	GetAllTags() ([]Tag, error)
	SearchTags(s string) ([]Tag, error)
}

type dbTagReader struct {
	db *sqlx.DB
}

// NewTagReader creates a new instance using an existing
// connection to a repository
func NewTagReader(conn RepositoryConnection) (TagReader, error) {
	if conn.h == nil {
		return nil, fmt.Errorf("no repository connection available")
	}
	return dbTagReader{db: conn.h}, nil
}

// GetAllTags returns all available tags in alphabetical order
func (r dbTagReader) GetAllTags() ([]Tag, error) {
	var tags []Tag
	if err := r.db.Select(&tags, "SELECT t.id, t.name FROM TAGS t ORDER BY name ASC"); err != nil {
		return nil, err
	}
	return tags, nil
}

// SearchTags returns tags matching the given search string
// the search string is matched independent of case and works in a wildcard fashion
func (r dbTagReader) SearchTags(s string) ([]Tag, error) {
	var tags []Tag
	search := strings.ToLower(s)
	search = "%" + search + "%"
	if err := r.db.Select(&tags, "SELECT t.id, t.name FROM TAGS t WHERE t.name LIKE ? ORDER BY name ASC", search); err != nil {
		return nil, err
	}
	return tags, nil
}
