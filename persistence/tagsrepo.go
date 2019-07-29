package persistence

import "strings"

// TagReader search for tags in the store
type TagReader interface {
	GetAllTags() ([]Tag, error)
	SearchTags(s string) ([]Tag, error)
}

// GetAllTags returns all available tags in alphabetical order
func (d dbstore) GetAllTags() ([]Tag, error) {
	var tags []Tag
	if err := d.db.Select(&tags, "SELECT t.id, t.name FROM TAGS t ORDER BY name ASC"); err != nil {
		return nil, err
	}
	return tags, nil
}

// SearchTags returns tags matching the given search string
// the searchstring is matched independent of case and works in a wildcard fashion
func (d dbstore) SearchTags(s string) ([]Tag, error) {
	var tags []Tag
	search := strings.ToLower(s)
	search = "%" + search + "%"
	if err := d.db.Select(&tags, "SELECT t.id, t.name FROM TAGS t WHERE t.name LIKE ? ORDER BY name ASC", search); err != nil {
		return nil, err
	}
	return tags, nil
}
