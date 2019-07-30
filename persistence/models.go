package persistence

// Tag is used to categorize a document
type Tag struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}

// Sender is the originator of a document
type Sender struct {
	ID   int64  `db:"id"`
	Name string `db:"name"`
}
