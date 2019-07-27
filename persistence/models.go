package persistence

// Tag represents an entry in the database
type Tag struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}
