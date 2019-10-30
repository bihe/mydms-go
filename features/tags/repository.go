package tags

import (
	"fmt"
	"strings"

	"github.com/bihe/mydms/internal/persistence"
)

// TagEntity is used to categorize a document
type TagEntity struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

// Repository search for tags and saves tags in the store
type Repository interface {
	// GetAlltags retrieves all available tags from the storage
	GetAllTags() ([]TagEntity, error)
	// SearchTags returns a tag based on the supplied search term
	SearchTags(s string) ([]TagEntity, error)
	// SaveTags processes the list of tag-names and stores the entries if not available
	SaveTags(tags []string, a persistence.Atomic) (err error)
	// CreateTag creates a tag with the given name or returns an existing one
	CreateTag(name string, a persistence.Atomic) (tag TagEntity, err error)
	// GetTagByName returns the entity defined by the given name
	GetTagByName(name string) (TagEntity, error)
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

// GetAllTags returns all available tags in alphabetical order
func (r *dbRepository) GetAllTags() ([]TagEntity, error) {
	var tags []TagEntity
	if err := r.c.Select(&tags, "SELECT t.id, t.name FROM TAGS t ORDER BY name ASC"); err != nil {
		return nil, err
	}
	return tags, nil
}

// SearchTags returns tags matching the given search string
// the search string is matched independent of case and works in a wildcard fashion
func (r *dbRepository) SearchTags(s string) ([]TagEntity, error) {
	var tags []TagEntity
	search := strings.ToLower(s)
	search = "%" + search + "%"
	if err := r.c.Select(&tags, "SELECT t.id, t.name FROM TAGS t WHERE lower(t.name) LIKE ? ORDER BY name ASC", search); err != nil {
		return nil, err
	}
	return tags, nil
}

// GetTagByName returns the tag by given name
// the search is performed ignoring case sensitivity
func (r *dbRepository) GetTagByName(name string) (TagEntity, error) {
	var tag TagEntity
	search := strings.ToLower(name)
	if err := r.c.Get(&tag, "SELECT t.id, t.name FROM TAGS t WHERE lower(t.name) = ?", search); err != nil {
		return TagEntity{}, err
	}
	return tag, nil
}

// SaveTags takes a slice of strings and saves tag entries if they do not exist
// the existance-check is done by comparing the tag-name
func (r *dbRepository) SaveTags(tags []string, a persistence.Atomic) (err error) {
	var atomic *persistence.Atomic

	defer func() {
		err = persistence.HandleTX(!a.Active, atomic, err)
	}()

	if atomic, err = persistence.CheckTX(r.c, &a); err != nil {
		return
	}

	var c int
	for _, t := range tags {
		t = strings.ToLower(t)
		err = atomic.Get(&c, "SELECT count(t.id) FROM TAGS t WHERE lower(t.name) = ?", t)
		if err != nil {
			err = fmt.Errorf("could not search for a tag: %v", err)
			return
		}
		if c > 0 {
			continue
		}

		_, err = atomic.Exec("INSERT INTO TAGS (name) VALUES (?)", t)
		if err != nil {
			err = fmt.Errorf("cannot save tag item: %v", err)
			return
		}
	}
	return
}

// CreateTag creates a tag with the given name or returns an existing one
func (r *dbRepository) CreateTag(name string, a persistence.Atomic) (tag TagEntity, err error) {
	var atomic *persistence.Atomic

	defer func() {
		err = persistence.HandleTX(!a.Active, atomic, err)
	}()

	if atomic, err = persistence.CheckTX(r.c, &a); err != nil {
		return
	}

	err = atomic.Get(&tag, "SELECT id,name FROM TAGS t WHERE lower(t.name) = ?", strings.ToLower(name))
	if err == nil {
		return tag, nil
	}

	res, e := atomic.Exec("INSERT INTO TAGS (name) VALUES (?)", name)
	if e != nil {
		err = fmt.Errorf("cannot save tag item, %v", e)
		return
	}
	id, e := res.LastInsertId()
	if e != nil {
		err = fmt.Errorf("could not get last inserted id, %v", e)
		return
	}
	return TagEntity{ID: int(id), Name: name}, nil
}
