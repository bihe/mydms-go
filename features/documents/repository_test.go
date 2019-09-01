package documents

import (
	"database/sql"
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
const deleteExpErr = "error was not expected while delete item: %v"
const existsErr = "error was not expected while checking for existence of item: %v"
const expected = "error expected"

func TestAtomic(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(fatalErr, err)
	}
	defer db.Close()

	dbx := sqlx.NewDb(db, "mysql")
	repo, err := NewRepository(persistence.NewFromDB(dbx))
	if err != nil {
		t.Errorf("could not get a repository: %v", err)
	}

	mock.ExpectBegin()

	_, err = repo.CreateAtomic()
	if err != nil {
		t.Errorf("could not ceate a new atomic object: %v", err)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf(expectations, err)
	}
}

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

func TestSave(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(fatalErr, err)
	}
	defer db.Close()

	var now = time.Now().UTC()

	dbx := sqlx.NewDb(db, "mysql")
	c := persistence.NewFromDB(dbx)
	rw := dbRepository{c}

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
	rw := dbRepository{c}

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
	rw := dbRepository{c}
	columns := []string{"id", "title", "filename", "alternativeid", "previewlink", "amount", "taglist", "senderlist", "created", "modified"}
	q := "SELECT id,title,filename,alternativeid,previewlink,amount,taglist,senderlist,created,modified FROM DOCUMENTS"
	id := "id"

	expected := DocumentEntity{
		ID:          "id",
		Title:       "title",
		FileName:    "filename",
		AltID:       "altid",
		PreviewLink: sql.NullString{String: "previewlink", Valid: true},
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

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(fatalErr, err)
	}
	defer db.Close()

	dbx := sqlx.NewDb(db, "mysql")
	c := persistence.NewFromDB(dbx)
	rw := dbRepository{c}
	stmt := "DELETE FROM DOCUMENTS"

	item := DocumentEntity{
		ID: "id",
	}

	mock.ExpectBegin()
	mock.ExpectExec(stmt).WithArgs(item.ID).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// now we execute our method
	if err = rw.Delete(item.ID, persistence.Atomic{}); err != nil {
		t.Errorf(deleteExpErr, err)
	}

	// externally supplied tx
	mock.ExpectBegin()
	mock.ExpectExec(stmt).WithArgs(item.ID).WillReturnResult(sqlmock.NewResult(1, 1))

	a, err := c.CreateAtomic()
	if err = rw.Delete(item.ID, a); err != nil {
		t.Errorf("error was not expected while delete item: %v", err)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf(expectations, err)
	}
}

func TestExists(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(fatalErr, err)
	}
	defer db.Close()

	dbx := sqlx.NewDb(db, "mysql")
	c := persistence.NewFromDB(dbx)
	rw := dbRepository{c}
	q := "SELECT count\\(id\\) FROM DOCUMENTS"
	id := "id"
	rows := []string{"count(id)"}

	mock.ExpectBegin()
	mock.ExpectQuery(q).WithArgs(id).WillReturnRows(sqlmock.NewRows(rows).AddRow(1))
	mock.ExpectCommit()

	// now we execute our method
	if err = rw.Exists(id, persistence.Atomic{}); err != nil {
		t.Errorf(existsErr, err)
	}

	// externally supplied tx
	mock.ExpectBegin()
	mock.ExpectQuery(q).WithArgs(id).WillReturnRows(sqlmock.NewRows(rows).AddRow(1))

	a, err := c.CreateAtomic()
	if err = rw.Exists(id, a); err != nil {
		t.Errorf(existsErr, err)
	}

	// error
	mock.ExpectBegin()
	mock.ExpectQuery(q).WithArgs(id).WillReturnError(fmt.Errorf("error"))
	mock.ExpectRollback()

	if err = rw.Exists(id, persistence.Atomic{}); err == nil {
		t.Errorf(expected)
	}

	// notfound
	mock.ExpectBegin()
	mock.ExpectQuery(q).WithArgs(id).WillReturnRows(sqlmock.NewRows(rows).AddRow(0))
	mock.ExpectRollback()

	if err = rw.Exists(id, persistence.Atomic{}); err == nil {
		t.Errorf(expected)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf(expectations, err)
	}
}

func TestDeleteError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(fatalErr, err)
	}
	defer db.Close()

	dbx := sqlx.NewDb(db, "mysql")
	c := persistence.NewFromDB(dbx)
	rw := dbRepository{c}

	item := DocumentEntity{
		ID: "id",
	}

	mock.ExpectBegin()
	mock.ExpectRollback()

	// now we execute our method
	if err = rw.Delete(item.ID, persistence.Atomic{}); err == nil {
		t.Errorf("error was expected for insert item")
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf(expectations, err)
	}
}

func TestSearch(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(fatalErr, err)
	}
	defer db.Close()

	dbx := sqlx.NewDb(db, "mysql")
	c := persistence.NewFromDB(dbx)
	rw := dbRepository{c}
	columns := []string{"id", "title", "filename", "alternativeid", "previewlink", "amount", "taglist", "senderlist", "created", "modified"}

	q := "SELECT id,title,filename,alternativeid,previewlink,amount,taglist,senderlist,created,modified FROM DOCUMENTS"
	qc := "SELECT count\\(id\\) FROM DOCUMENTS"

	expected := DocumentEntity{
		ID:          "id",
		Title:       "title",
		FileName:    "filename",
		AltID:       "altid",
		PreviewLink: sql.NullString{String: "previewlink", Valid: true},
		Amount:      1.0,
		Created:     time.Now().UTC(),
		Modified:    mysql.NullTime{},
		TagList:     "tags",
		SenderList:  "senders",
	}

	// success
	cr := sqlmock.NewRows([]string{"count(id)"}).AddRow(1)
	mock.ExpectQuery(qc).WillReturnRows(cr)

	dr := sqlmock.NewRows(columns).
		AddRow(expected.ID, expected.Title, expected.FileName, expected.AltID, expected.PreviewLink, expected.Amount, expected.TagList, expected.SenderList, expected.Created, expected.Modified)
	mock.ExpectQuery(q).WillReturnRows(dr)

	ts := time.Now().UTC()
	from := ts.Add(-time.Hour)
	until := ts.Add(time.Hour)
	search := DocSearch{
		Skip:   1,
		Limit:  1,
		Title:  "title",
		Tag:    "tags",
		Sender: "senders",
		From:   from,
		Until:  until,
	}
	order := []OrderBy{
		OrderBy{Order: DESC, Field: "modified"},
		OrderBy{Order: ASC, Field: "title"},
	}

	doc, err := rw.Search(search, order)
	if err != nil {
		t.Errorf("could not query documents: %v", err)
	}
	if len(doc.Documents) == 0 {
		t.Errorf("document list empty, nothing returned")
	}

	item := doc.Documents[0]

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

	// failure1
	mock.ExpectQuery(qc).WillReturnError(fmt.Errorf("could not get count"))
	_, err = rw.Search(search, order)
	if err == nil {
		t.Errorf("error expected")
	}

	// failure2
	cr = sqlmock.NewRows([]string{"count(id)"}).AddRow(1)
	mock.ExpectQuery(qc).WillReturnRows(cr)
	mock.ExpectQuery(q).WillReturnError(fmt.Errorf("could not get documents"))
	_, err = rw.Search(search, order)
	if err == nil {
		t.Errorf("error expected")
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf(expectations, err)
	}
}