package persistence

import "strings"

// GetAllTags returns all available tags in alphabetical order
func (r *Repository) GetAllTags() ([]Tag, error) {
	var tags []Tag
	if err := r.db.Select(&tags, "SELECT t.id, t.name FROM tags t ORDER BY name ASC"); err != nil {
		return nil, err
	}
	return tags, nil
}

// SearchTags returns tags matching the given search string
// the searchstring is matched independent of case and works in a wildcard fashion
func (r *Repository) SearchTags(s string) ([]Tag, error) {
	var tags []Tag
	search := strings.ToLower(s)
	search = "%" + search + "%"
	if err := r.db.Select(&tags, "SELECT t.id, t.name FROM tags t WHERE t.name LIKE ? ORDER BY name ASC", search); err != nil {
		return nil, err
	}
	return tags, nil
}
