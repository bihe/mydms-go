package senders

import (
	"fmt"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bihe/mydms/persistence"
	"github.com/jmoiron/sqlx"
)

func TestNewSenderReader(t *testing.T) {
	_, err := NewReader(persistence.Connection{})
	if err == nil {
		t.Errorf("no reader without connection possible")
	}

	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	dbx := sqlx.NewDb(db, "mysql")
	_, err = NewReader(persistence.NewFromDB(dbx))
	if err != nil {
		t.Errorf("could not get a reader: %v", err)
	}
}

func TestGetAllEntitySenders(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	dbx := sqlx.NewDb(db, "mysql")
	c := persistence.NewFromDB(dbx)
	r := dbSenderReader{c}

	// success
	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "Sender1").
		AddRow(2, "Sender2")
	mock.ExpectQuery("SELECT t.id, t.name FROM SENDERS t ORDER BY name ASC").WillReturnRows(rows)

	senders, err := r.GetAllSenders()
	if err != nil {
		t.Errorf("could not get all senders: %v", err)
	}
	if len(senders) != 2 {
		t.Errorf("expected 2 items, got %d", len(senders))
	}

	// no results
	rows = sqlmock.NewRows([]string{"id", "name"})
	mock.ExpectQuery("SELECT t.id, t.name FROM SENDERS t ORDER BY name ASC").WillReturnRows(rows)
	senders, err = r.GetAllSenders()
	if err != nil {
		t.Errorf("could not get all senders: %v", err)
	}
	if len(senders) != 0 {
		t.Errorf("expected 0 items, got %d", len(senders))
	}

	// error
	mock.ExpectQuery("SELECT t.id, t.name FROM SENDERS t ORDER BY name ASC").WillReturnError(fmt.Errorf("no rows"))
	senders, err = r.GetAllSenders()
	if err == nil {
		t.Errorf("error during SQL expected")
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSearchForEntitySenders(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	dbx := sqlx.NewDb(db, "mysql")
	c := persistence.NewFromDB(dbx)
	r := dbSenderReader{c}

	// excact match
	search := "Sender1"
	rows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(1, "Sender1")
	mock.ExpectQuery("SELECT t.id, t.name FROM SENDERS t").WithArgs("%" + strings.ToLower(search) + "%").WillReturnRows(rows)
	senders, err := r.SearchSenders(search)
	if err != nil {
		t.Errorf("could not search for senders by '%s': %v", "Sender1", err)
	}
	if len(senders) != 1 {
		t.Errorf("expected 1 items, got %d", len(senders))
	}

	// no match
	search = "_no_sender_"
	rows = sqlmock.NewRows([]string{"id", "name"})
	mock.ExpectQuery("SELECT t.id, t.name FROM SENDERS t").WithArgs("%" + strings.ToLower(search) + "%").WillReturnRows(rows)
	senders, err = r.SearchSenders(search)
	if err != nil {
		t.Errorf("could not search for senders by '%s': %v", "Sender1", err)
	}
	if len(senders) != 0 {
		t.Errorf("expected 0 items, got %d", len(senders))
	}

	// error
	search = "foo"
	mock.ExpectQuery("SELECT t.id, t.name FROM SENDERS t").WithArgs("%" + strings.ToLower(search) + "%").WillReturnError(fmt.Errorf("no rows"))
	senders, err = r.SearchSenders(search)
	if err == nil {
		t.Errorf("error during SQL expected")
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
