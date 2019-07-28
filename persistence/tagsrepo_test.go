package persistence

import (
	"fmt"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func TestGetAllTags(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	dbx := sqlx.NewDb(db, "mysql")
	r := repository{dbx}

	// success
	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "Tag1").
		AddRow(2, "Tag2")
	mock.ExpectQuery("SELECT t.id, t.name FROM tags t ORDER BY name ASC").WillReturnRows(rows)

	tags, err := r.GetAllTags()
	if err != nil {
		t.Errorf("could not get all tags: %v", err)
	}
	if len(tags) != 2 {
		t.Errorf("expected 2 items, got %d", len(tags))
	}

	// no results
	rows = sqlmock.NewRows([]string{"id", "name"})
	mock.ExpectQuery("SELECT t.id, t.name FROM tags t ORDER BY name ASC").WillReturnRows(rows)
	tags, err = r.GetAllTags()
	if err != nil {
		t.Errorf("could not get all tags: %v", err)
	}
	if len(tags) != 0 {
		t.Errorf("expected 0 items, got %d", len(tags))
	}

	// error
	mock.ExpectQuery("SELECT t.id, t.name FROM tags t ORDER BY name ASC").WillReturnError(fmt.Errorf("no rows"))
	tags, err = r.GetAllTags()
	if err == nil {
		t.Errorf("error during SQL expected")
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSearchForTags(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	dbx := sqlx.NewDb(db, "mysql")
	r := repository{dbx}

	// excact match
	search := "Tag1"
	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "Tag1")
	mock.ExpectQuery("SELECT t.id, t.name FROM tags t").WithArgs("%" + strings.ToLower(search) + "%").WillReturnRows(rows)
	tags, err := r.SearchTags(search)
	if err != nil {
		t.Errorf("could not search for tags by '%s': %v", "Tag1", err)
	}
	if len(tags) != 1 {
		t.Errorf("expected 1 items, got %d", len(tags))
	}

	// no match
	search = "_no_tag_"
	rows = sqlmock.NewRows([]string{"id", "name"})
	mock.ExpectQuery("SELECT t.id, t.name FROM tags t").WithArgs("%" + strings.ToLower(search) + "%").WillReturnRows(rows)
	tags, err = r.SearchTags(search)
	if err != nil {
		t.Errorf("could not search for tags by '%s': %v", "Tag1", err)
	}
	if len(tags) != 0 {
		t.Errorf("expected 0 items, got %d", len(tags))
	}

	// error
	search = "foo"
	mock.ExpectQuery("SELECT t.id, t.name FROM tags t").WithArgs("%" + strings.ToLower(search) + "%").WillReturnError(fmt.Errorf("no rows"))
	tags, err = r.SearchTags(search)
	if err == nil {
		t.Errorf("error during SQL expected")
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
