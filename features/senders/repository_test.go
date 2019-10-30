package senders

import (
	"fmt"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bihe/mydms/internal/persistence"
	"github.com/jmoiron/sqlx"
)

const fatalErr = "an error '%s' was not expected when opening a stub database connection"
const expectations = "there were unfulfilled expectations: %s"
const errExpected = "error during SQL expected"

var errNoRows = fmt.Errorf("no rows")
var rowDef = []string{"id", "name"}

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

func TestGetAllEntitySenders(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(fatalErr, err)
	}
	defer db.Close()

	var errCould = "could not get all senders: %v"

	dbx := sqlx.NewDb(db, "mysql")
	c := persistence.NewFromDB(dbx)
	r := dbRepository{c}
	q := "SELECT t.id, t.name FROM SENDERS t ORDER BY name ASC"

	// success
	rows := sqlmock.NewRows(rowDef).
		AddRow(1, "Sender1").
		AddRow(2, "Sender2")
	mock.ExpectQuery(q).WillReturnRows(rows)

	senders, err := r.GetAllSenders()
	if err != nil {
		t.Errorf(errCould, err)
	}
	if len(senders) != 2 {
		t.Errorf("expected 2 items, got %d", len(senders))
	}

	// no results
	rows = sqlmock.NewRows(rowDef)
	mock.ExpectQuery(q).WillReturnRows(rows)
	senders, err = r.GetAllSenders()
	if err != nil {
		t.Errorf(errCould, err)
	}
	if len(senders) != 0 {
		t.Errorf("expected 0 items, got %d", len(senders))
	}

	// error
	mock.ExpectQuery(q).WillReturnError(errNoRows)
	senders, err = r.GetAllSenders()
	if err == nil {
		t.Errorf(errExpected)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf(expectations, err)
	}
}

func TestSearchForEntitySenders(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(fatalErr, err)
	}
	defer db.Close()

	var errCouldSearch = "could not search for senders by '%s': %v"

	dbx := sqlx.NewDb(db, "mysql")
	c := persistence.NewFromDB(dbx)
	r := dbRepository{c}
	q := "SELECT t.id, t.name FROM SENDERS t"

	// excact match
	search := "Sender1"
	rows := sqlmock.NewRows(rowDef).
		AddRow(1, search)
	mock.ExpectQuery(q).WithArgs("%" + strings.ToLower(search) + "%").WillReturnRows(rows)
	senders, err := r.SearchSenders(search)
	if err != nil {
		t.Errorf(errCouldSearch, search, err)
	}
	if len(senders) != 1 {
		t.Errorf("expected 1 items, got %d", len(senders))
	}

	// no match
	search = "_no_sender_"
	rows = sqlmock.NewRows(rowDef)
	mock.ExpectQuery(q).WithArgs("%" + strings.ToLower(search) + "%").WillReturnRows(rows)
	senders, err = r.SearchSenders(search)
	if err != nil {
		t.Errorf(errCouldSearch, search, err)
	}
	if len(senders) != 0 {
		t.Errorf("expected 0 items, got %d", len(senders))
	}

	// single search
	rows = sqlmock.NewRows(rowDef).
		AddRow(1, search)
	mock.ExpectQuery(q).WithArgs(strings.ToLower(search)).WillReturnRows(rows)
	sender, err := r.GetSenderByName(search)
	if err != nil {
		t.Errorf(errCouldSearch, search, err)
	}
	if sender.Name != search {
		t.Errorf("could not get single result for search '%s'", search)
	}

	// error
	search = "foo"
	mock.ExpectQuery(q).WithArgs("%" + strings.ToLower(search) + "%").WillReturnError(errNoRows)
	senders, err = r.SearchSenders(search)
	if err == nil {
		t.Errorf(errExpected)
	}

	// single-search error
	search = "foo"
	mock.ExpectQuery(q).WithArgs(strings.ToLower(search)).WillReturnError(errNoRows)
	sender, err = r.GetSenderByName(search)
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

	q := "SELECT count\\(s.id\\) FROM SENDERS"
	f := "count\\(s.id\\)"
	stmt := "INSERT INTO SENDERS \\(name\\)"

	senders := []string{"sender1", "sender2"}

	mock.ExpectBegin()
	mock.ExpectQuery(q).WillReturnRows(sqlmock.NewRows([]string{f}).AddRow(1))
	mock.ExpectQuery(q).WillReturnRows(sqlmock.NewRows([]string{f}).AddRow(0))
	mock.ExpectExec(stmt).WithArgs(senders[1]).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err = r.SaveSenders(senders, persistence.Atomic{})
	if err != nil {
		t.Errorf("could not save all senders: %v", err)
	}

	// error
	mock.ExpectBegin()
	mock.ExpectQuery(q).WillReturnError(errNoRows)
	mock.ExpectRollback()

	err = r.SaveSenders(senders, persistence.Atomic{})
	if err == nil {
		t.Errorf(errExpected)
	}

	// error
	mock.ExpectBegin()
	mock.ExpectQuery(q).WillReturnRows(sqlmock.NewRows([]string{f}).AddRow(1))
	mock.ExpectQuery(q).WillReturnRows(sqlmock.NewRows([]string{f}).AddRow(0))
	mock.ExpectExec(stmt).WillReturnError(errNoRows)
	mock.ExpectRollback()

	err = r.SaveSenders(senders, persistence.Atomic{})
	if err == nil {
		t.Errorf(errExpected)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf(expectations, err)
	}
}

func TestCreateSender(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(fatalErr, err)
	}
	defer db.Close()

	dbx := sqlx.NewDb(db, "mysql")
	c := persistence.NewFromDB(dbx)
	r := dbRepository{c}

	q := "SELECT id,name FROM SENDERS"
	stmt := "INSERT INTO SENDERS \\(name\\)"
	senderName := "Sender1"

	errSave := "could not save sender: %v"
	errID := "the ID of the created sender is 0"
	Err := fmt.Errorf("error")

	// Existing Sender
	mock.ExpectBegin()
	rows := sqlmock.NewRows(rowDef).
		AddRow(1, senderName)
	mock.ExpectQuery(q).WillReturnRows(rows)
	mock.ExpectCommit()

	tag, err := r.CreateSender(senderName, persistence.Atomic{})
	if err != nil {
		t.Errorf(errSave, err)
	}
	if tag.ID == 0 {
		t.Errorf(errID)
	}

	// New Sender
	mock.ExpectBegin()
	mock.ExpectQuery(q).WillReturnError(errNoRows)
	mock.ExpectExec(stmt).WithArgs(senderName).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	tag, err = r.CreateSender(senderName, persistence.Atomic{})
	if err != nil {
		t.Errorf(errSave, err)
	}
	if tag.ID == 0 {
		t.Errorf(errID)
	}

	// error1
	mock.ExpectBegin()
	mock.ExpectQuery(q).WillReturnError(errNoRows)
	mock.ExpectExec(stmt).WithArgs(senderName).WillReturnError(Err)
	mock.ExpectRollback()

	tag, err = r.CreateSender(senderName, persistence.Atomic{})
	if err == nil {
		t.Errorf(errExpected)
	}

	// error2
	mock.ExpectBegin()
	mock.ExpectQuery(q).WillReturnError(errNoRows)
	mock.ExpectExec(stmt).WithArgs(senderName).WillReturnResult(sqlmock.NewErrorResult(Err))
	mock.ExpectRollback()

	tag, err = r.CreateSender(senderName, persistence.Atomic{})
	if err == nil {
		t.Errorf(errExpected)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf(expectations, err)
	}

}
