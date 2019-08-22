package documents

import (
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bihe/mydms/persistence"
	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

const fatalErr = "an error '%s' was not expected when opening a stub database connection"
const expectations = "there were unfulfilled expectations: %s"

func TestSave(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(fatalErr, err)
	}
	defer db.Close()

	var now = time.Now().UTC()

	dbx := sqlx.NewDb(db, "mysql")
	c := persistence.NewFromDB(dbx)
	rw := dbDocumentReaderWriter{c}

	item := DocumentEntity{
		Title:      "title",
		FileName:   "filename",
		Amount:     10,
		TagList:    "taglist",
		SenderList: "senderlist",
	}

	// INSERT
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO DOCUMENTS").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	var d DocumentEntity
	if d, err = rw.Save(item, persistence.Atomic{}); err != nil {
		t.Errorf("error was not expected while insert item: %v", err)
	}
	assert.Equal(t, item.Title, d.Title)
	assert.Equal(t, item.FileName, d.FileName)
	assert.Equal(t, item.Amount, d.Amount)
	assert.Equal(t, item.TagList, d.TagList)
	assert.Equal(t, item.SenderList, d.SenderList)
	assert.True(t, d.ID != "")
	assert.True(t, d.AltID != "")

	// UPDATE
	mock.ExpectBegin()
	item.ID = uuid.New().String()
	item.AltID = d.AltID

	rows := sqlmock.NewRows([]string{"id", "title", "filename", "alternativeid", "previewlink", "amount", "taglist", "senderlist", "created", "modified"}).
		AddRow(item.ID, item.Title, item.FileName, item.AltID, item.PreviewLink, item.Amount, item.TagList, item.SenderList, d.Created, nil)
	mock.ExpectQuery("SELECT id,title,filename,alternativeid,previewlink,amount,taglist,senderlist,created,modified FROM DOCUMENTS").WillReturnRows(rows)
	mock.ExpectExec("UPDATE DOCUMENTS").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	var up DocumentEntity
	if up, err = rw.Save(item, persistence.Atomic{}); err != nil {
		t.Errorf("error was not expected while insert item: %v", err)
	}
	assert.Equal(t, item.ID, up.ID)
	assert.Equal(t, item.AltID, up.AltID)
	assert.Equal(t, item.Title, up.Title)
	assert.Equal(t, item.FileName, up.FileName)
	assert.Equal(t, item.Amount, up.Amount)
	assert.Equal(t, item.TagList, up.TagList)
	assert.Equal(t, item.SenderList, up.SenderList)
	assert.Equal(t, d.Created, up.Created)
	assert.True(t, up.Modified.Time.After(now))

	// UPDATE with wrong ID
	mock.ExpectBegin()
	item.ID = uuid.New().String()
	item.AltID = d.AltID

	mock.ExpectQuery("SELECT id,title,filename,alternativeid,previewlink,amount,taglist,senderlist,created,modified FROM DOCUMENTS").WillReturnError(fmt.Errorf("no rows"))
	mock.ExpectExec("INSERT INTO DOCUMENTS").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	if up, err = rw.Save(item, persistence.Atomic{}); err != nil {
		t.Errorf("error was not expected while insert item: %v", err)
	}

	assert.NotEqual(t, item.AltID, up.AltID)
	assert.Equal(t, item.Title, up.Title)
	assert.Equal(t, item.FileName, up.FileName)
	assert.Equal(t, item.Amount, up.Amount)
	assert.Equal(t, item.TagList, up.TagList)
	assert.Equal(t, item.SenderList, up.SenderList)
	assert.True(t, up.Created.After(now))

	// externally supplied tx
	item.ID = ""
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO DOCUMENTS").WillReturnResult(sqlmock.NewResult(1, 1))
	a, _ := c.CreateAtomic()
	if d, err = rw.Save(item, a); err != nil {
		t.Errorf("error was not expected while insert item: %v", err)
	}
	assert.Equal(t, item.Title, d.Title)
	assert.Equal(t, item.FileName, d.FileName)
	assert.Equal(t, item.Amount, d.Amount)
	assert.Equal(t, item.TagList, d.TagList)
	assert.Equal(t, item.SenderList, d.SenderList)
	assert.True(t, d.ID != "")
	assert.True(t, d.AltID != "")

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf(expectations, err)
	}
}

func TestSaveError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(fatalErr, err)
	}
	defer db.Close()

	dbx := sqlx.NewDb(db, "mysql")
	c := persistence.NewFromDB(dbx)
	rw := dbDocumentReaderWriter{c}

	item := DocumentEntity{
		Title:      "title",
		FileName:   "filename",
		Amount:     10,
		TagList:    "taglist",
		SenderList: "senderlist",
	}

	// INSERT Error
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO DOCUMENTS").WillReturnError(fmt.Errorf("does not work"))
	mock.ExpectRollback()

	if _, err = rw.Save(item, persistence.Atomic{}); err == nil {
		t.Errorf("error was expected while insert item")
	}

	// Rows affected Error
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO DOCUMENTS").WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("result error")))
	mock.ExpectRollback()

	if _, err = rw.Save(item, persistence.Atomic{}); err == nil {
		t.Errorf("error was expected while insert item")
	}

	// Rows affected number mismatch
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO DOCUMENTS").WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectRollback()

	if _, err = rw.Save(item, persistence.Atomic{}); err == nil {
		t.Errorf("error was expected while insert item")
	}
}

func TestRead(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(fatalErr, err)
	}
	defer db.Close()

	dbx := sqlx.NewDb(db, "mysql")
	c := persistence.NewFromDB(dbx)
	rw := dbDocumentReaderWriter{c}
	columns := []string{"id", "title", "filename", "alternativeid", "previewlink", "amount", "taglist", "senderlist", "created", "modified"}
	q := "SELECT id,title,filename,alternativeid,previewlink,amount,taglist,senderlist,created,modified FROM DOCUMENTS"
	id := "id"

	expected := DocumentEntity{
		ID:          "id",
		Title:       "title",
		FileName:    "filename",
		AltID:       "altid",
		PreviewLink: "previewlink",
		Amount:      1.0,
		Created:     time.Now().UTC(),
		Modified:    mysql.NullTime{},
		TagList:     "tags",
		SenderList:  "senders",
	}

	// success
	rows := sqlmock.NewRows(columns).
		AddRow(expected.ID, expected.Title, expected.FileName, expected.AltID, expected.PreviewLink, expected.Amount, expected.TagList, expected.SenderList, expected.Created, expected.Modified)
	mock.ExpectQuery(q).WithArgs(id).WillReturnRows(rows)

	item, err := rw.Get(id)
	if err != nil {
		t.Errorf("could not get item: %v", err)
	}

	assert.Equal(t, expected.ID, item.ID)
	assert.Equal(t, expected.Title, item.Title)
	assert.Equal(t, expected.FileName, item.FileName)
	assert.Equal(t, expected.AltID, item.AltID)
	assert.Equal(t, expected.PreviewLink, item.PreviewLink)
	assert.Equal(t, expected.Amount, item.Amount)
	assert.Equal(t, expected.TagList, item.TagList)
	assert.Equal(t, expected.SenderList, item.SenderList)
	assert.Equal(t, expected.Created, item.Created)
	assert.Equal(t, expected.Modified, item.Modified)

	// no result
	rows = sqlmock.NewRows(columns)
	mock.ExpectQuery(q).WithArgs(id).WillReturnRows(rows)

	item, err = rw.Get(id)
	if err == nil {
		t.Errorf("should have returned an error")
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf(expectations, err)
	}
}
