package tags

import (
	"fmt"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bihe/mydms/persistence"
	"github.com/jmoiron/sqlx"
)

const fatalErr = "an error '%s' was not expected when opening a stub database connection"
const expectations = "there were unfulfilled expectations: %s"
const errExpected = "error during SQL expected"

var errNoRows = fmt.Errorf("no rows")
var rowDef = []string{"id", "name"}
var tags = []string{"tag1", "tag2"}

func TestNewRepository(t *testing.T) {
	_, err := NewRepository(persistence.Connection{})
	if err == nil {
		t.Errorf("no reader without connection possible")
	}

	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf(fatalErr, err)
	}
	defer db.Close()

	dbx := sqlx.NewDb(db, "mysql")
	_, err = NewRepository(persistence.NewFromDB(dbx))
	if err != nil {
		t.Errorf("could not get a repository: %v", err)
	}
}

func TestGetAllEntityTags(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(fatalErr, err)
	}
	defer db.Close()

	var errCould = "could not get all tags: %v"

	dbx := sqlx.NewDb(db, "mysql")
	c := persistence.NewFromDB(dbx)
	r := dbRepository{c}
	q := "SELECT t.id, t.name FROM TAGS t ORDER BY name ASC"

	// success
	rows := sqlmock.NewRows(rowDef).
		AddRow(1, "Tag1").
		AddRow(2, "Tag2")
	mock.ExpectQuery(q).WillReturnRows(rows)

	tags, err := r.GetAllTags()
	if err != nil {
		t.Errorf(errCould, err)
	}
	if len(tags) != 2 {
		t.Errorf("expected 2 items, got %d", len(tags))
	}

	// no results
	rows = sqlmock.NewRows(rowDef)
	mock.ExpectQuery(q).WillReturnRows(rows)
	tags, err = r.GetAllTags()
	if err != nil {
		t.Errorf(errCould, err)
	}
	if len(tags) != 0 {
		t.Errorf("expected 0 items, got %d", len(tags))
	}

	// error
	mock.ExpectQuery(q).WillReturnError(errNoRows)
	tags, err = r.GetAllTags()
	if err == nil {
		t.Errorf(errExpected)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf(expectations, err)
	}
}

func TestSearchForEntityTags(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(fatalErr, err)
	}
	defer db.Close()

	var errCouldSearch = "could not search for tags by '%s': %v"

	dbx := sqlx.NewDb(db, "mysql")
	c := persistence.NewFromDB(dbx)
	r := dbRepository{c}
	q := "SELECT t.id, t.name FROM TAGS t"

	// excact match
	search := "Tag1"
	rows := sqlmock.NewRows(rowDef).
		AddRow(1, search)
	mock.ExpectQuery(q).WithArgs("%" + strings.ToLower(search) + "%").WillReturnRows(rows)
	tags, err := r.SearchTags(search)
	if err != nil {
		t.Errorf(errCouldSearch, search, err)
	}
	if len(tags) != 1 {
		t.Errorf("expected 1 items, got %d", len(tags))
	}

	// no match
	search = "_no_tag_"
	rows = sqlmock.NewRows(rowDef)
	mock.ExpectQuery(q).WithArgs("%" + strings.ToLower(search) + "%").WillReturnRows(rows)
	tags, err = r.SearchTags(search)
	if err != nil {
		t.Errorf(errCouldSearch, search, err)
	}
	if len(tags) != 0 {
		t.Errorf("expected 0 items, got %d", len(tags))
	}

	// single search
	rows = sqlmock.NewRows(rowDef).
		AddRow(1, search)
	mock.ExpectQuery(q).WithArgs(strings.ToLower(search)).WillReturnRows(rows)
	tag, err := r.GetTagByName(search)
	if err != nil {
		t.Errorf(errCouldSearch, search, err)
	}
	if tag.Name != search {
		t.Errorf("could not get single result for search '%s'", search)
	}

	// error
	search = "foo"
	mock.ExpectQuery(q).WithArgs("%" + strings.ToLower(search) + "%").WillReturnError(errNoRows)
	tags, err = r.SearchTags(search)
	if err == nil {
		t.Errorf(errExpected)
	}

	// single-search error
	search = "foo"
	mock.ExpectQuery(q).WithArgs(strings.ToLower(search)).WillReturnError(errNoRows)
	tag, err = r.GetTagByName(search)
	if err == nil {
		t.Errorf(errExpected)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf(expectations, err)
	}
}

func TestSaveTags(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(fatalErr, err)
	}
	defer db.Close()

	dbx := sqlx.NewDb(db, "mysql")
	c := persistence.NewFromDB(dbx)
	r := dbRepository{c}

	q := "SELECT count\\(t.id\\) FROM TAGS"
	f := "count\\(t.id\\)"
	stmt := "INSERT INTO TAGS \\(name\\)"

	tags := []string{"tag1", "tag2"}

	mock.ExpectBegin()
	mock.ExpectQuery(q).WillReturnRows(sqlmock.NewRows([]string{f}).AddRow(1))
	mock.ExpectQuery(q).WillReturnRows(sqlmock.NewRows([]string{f}).AddRow(0))
	mock.ExpectExec(stmt).WithArgs(tags[1]).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = r.SaveTags(tags, persistence.Atomic{})
	if err != nil {
		t.Errorf("could not save all tags: %v", err)
	}

	// error
	mock.ExpectBegin()
	mock.ExpectQuery(q).WillReturnError(errNoRows)
	mock.ExpectRollback()

	err = r.SaveTags(tags, persistence.Atomic{})
	if err == nil {
		t.Errorf(errExpected)
	}

	// error
	mock.ExpectBegin()
	mock.ExpectQuery(q).WillReturnRows(sqlmock.NewRows([]string{f}).AddRow(1))
	mock.ExpectQuery(q).WillReturnRows(sqlmock.NewRows([]string{f}).AddRow(0))
	mock.ExpectExec(stmt).WillReturnError(errNoRows)
	mock.ExpectRollback()

	err = r.SaveTags(tags, persistence.Atomic{})
	if err == nil {
		t.Errorf(errExpected)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf(expectations, err)
	}

}

func TestCreateTag(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(fatalErr, err)
	}
	defer db.Close()

	dbx := sqlx.NewDb(db, "mysql")
	c := persistence.NewFromDB(dbx)
	r := dbRepository{c}

	q := "SELECT id,name FROM TAGS"
	stmt := "INSERT INTO TAGS \\(name\\)"
	tagName := "Tag1"

	errSave := "could not save tag: %v"
	errID := "the ID of the created tag is 0"
	Err := fmt.Errorf("error")

	// Existing Tag
	mock.ExpectBegin()
	rows := sqlmock.NewRows(rowDef).
		AddRow(1, tagName)
	mock.ExpectQuery(q).WillReturnRows(rows)
	mock.ExpectCommit()

	tag, err := r.CreateTag(tagName, persistence.Atomic{})
	if err != nil {
		t.Errorf(errSave, err)
	}
	if tag.ID == 0 {
		t.Errorf(errID)
	}

	// New Tag
	mock.ExpectBegin()
	mock.ExpectQuery(q).WillReturnError(errNoRows)
	mock.ExpectExec(stmt).WithArgs(tagName).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	tag, err = r.CreateTag(tagName, persistence.Atomic{})
	if err != nil {
		t.Errorf(errSave, err)
	}
	if tag.ID == 0 {
		t.Errorf(errID)
	}

	// error1
	mock.ExpectBegin()
	mock.ExpectQuery(q).WillReturnError(errNoRows)
	mock.ExpectExec(stmt).WithArgs(tagName).WillReturnError(Err)
	mock.ExpectRollback()

	tag, err = r.CreateTag(tagName, persistence.Atomic{})
	if err == nil {
		t.Errorf(errExpected)
	}

	// error2
	mock.ExpectBegin()
	mock.ExpectQuery(q).WillReturnError(errNoRows)
	mock.ExpectExec(stmt).WithArgs(tagName).WillReturnResult(sqlmock.NewErrorResult(Err))
	mock.ExpectRollback()

	tag, err = r.CreateTag(tagName, persistence.Atomic{})
	if err == nil {
		t.Errorf(errExpected)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf(expectations, err)
	}

}
