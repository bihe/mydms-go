package tags

import (
	"fmt"
	"strings"

	"github.com/bihe/mydms/persistence"
)

// TagEntity is used to categorize a document
type TagEntity struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

// Reader search for tags in the store
type Reader interface {
	GetAllTags() ([]TagEntity, error)
	SearchTags(s string) ([]TagEntity, error)
}

type dbTagReader struct {
	c persistence.Connection
}

// NewReader creates a new instance using an existing
// connection to a repository
func NewReader(c persistence.Connection) (Reader, error) {
	if !c.Active {
		return nil, fmt.Errorf("no repository connection available")
	}
	return dbTagReader{c}, nil
}

// GetAllTags returns all available tags in alphabetical order
func (r dbTagReader) GetAllTags() ([]TagEntity, error) {
	var tags []TagEntity
	if err := r.c.Select(&tags, "SELECT t.id, t.name FROM TAGS t ORDER BY name ASC"); err != nil {
		return nil, err
	}
	return tags, nil
}

// SearchTags returns tags matching the given search string
// the search string is matched independent of case and works in a wildcard fashion
func (r dbTagReader) SearchTags(s string) ([]TagEntity, error) {
	var tags []TagEntity
	search := strings.ToLower(s)
	search = "%" + search + "%"
	if err := r.c.Select(&tags, "SELECT t.id, t.name FROM TAGS t WHERE t.name LIKE ? ORDER BY name ASC", search); err != nil {
		return nil, err
	}
	return tags, nil
}
