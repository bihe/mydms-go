package persistence

// GetAllTags returns all available tags in alphabetical order
func (r *Repository) GetAllTags() ([]Tag, error) {
	var tags []Tag
	if err := r.db.Select(&tags, "SELECT t.id, t.name FROM tags t ORDER BY name ASC"); err != nil {
		return nil, err
	}
	return tags, nil
}
