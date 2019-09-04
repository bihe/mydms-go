package documents

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/bihe/mydms/features/filestore"
	"github.com/bihe/mydms/persistence"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

const invalidJSON = "could not get valid json: %v"
const errExp = "error expected"
const ID = "id"
const notExists = "!exists"
const noDelete = "!delete"
const noResult = "!result"
const noFileDelete = "!fileDelete"

/* MOCK
type Repository interface {
	Get(id string) (d DocumentEntity, err error)
	Exists(id string, a persistence.Atomic) (filePath string, err error)
	Save(doc DocumentEntity, a persistence.Atomic) (d DocumentEntity, err error)
	Delete(id string, a persistence.Atomic) (err error)
	Search(s DocSearch, order []OrderBy) (PagedDocuments, error)
}
*/

type mockRepository struct {
	c    persistence.Connection
	fail bool
}

func (m *mockRepository) Get(id string) (d DocumentEntity, err error) {
	if id == "" {
		return DocumentEntity{}, fmt.Errorf("no document")
	}
	return DocumentEntity{
		Modified:    sql.NullTime{Time: time.Now().UTC(), Valid: true},
		PreviewLink: sql.NullString{String: "string", Valid: true},
	}, nil
}

func (m *mockRepository) Save(doc DocumentEntity, a persistence.Atomic) (d DocumentEntity, err error) {
	return doc, nil
}

func (m *mockRepository) Delete(id string, a persistence.Atomic) (err error) {
	if id == noDelete {
		return fmt.Errorf("delete error")
	}
	return nil
}

func (m *mockRepository) Search(s DocSearch, order []OrderBy) (PagedDocuments, error) {
	if s.Title == noResult {
		return PagedDocuments{}, fmt.Errorf("search error")
	}

	return PagedDocuments{
		Count: 2,
		Documents: []DocumentEntity{
			DocumentEntity{
				Title:       "title1",
				FileName:    "filename1",
				Amount:      1,
				TagList:     "taglist1",
				SenderList:  "senderlist1",
				PreviewLink: sql.NullString{String: "previewlink", Valid: true},
				Created:     time.Now().UTC(),
				Modified:    sql.NullTime{Time: time.Now().UTC(), Valid: true},
			},
			DocumentEntity{
				Title:      "title2",
				FileName:   "filename2",
				Amount:     2,
				TagList:    "taglist2",
				SenderList: "senderlist2",
				Created:    time.Now().UTC(),
			},
		},
	}, nil
}

func (m *mockRepository) Exists(id string, a persistence.Atomic) (filePath string, err error) {
	if id == notExists {
		return "", fmt.Errorf("exists error")
	}
	if id == noFileDelete {
		return noFileDelete, nil
	}
	return "file", nil
}

func (m *mockRepository) CreateAtomic() (persistence.Atomic, error) {
	if m.fail {
		return persistence.Atomic{}, fmt.Errorf("start transaction failed")
	}
	return m.c.CreateAtomic()
}

/* MOCK
type FileService interface {
	SaveFile(file FileItem) error
	GetFile(filePath string) (FileItem, error)
	DeleteFile(filePath string) error
}
*/

// rather small PDF payload
// https://stackoverflow.com/questions/17279712/what-is-the-smallest-possible-valid-pdf
const pdfPayload = `%PDF-1.0
1 0 obj<</Type/Catalog/Pages 2 0 R>>endobj 2 0 obj<</Type/Pages/Kids[3 0 R]/Count 1>>endobj 3 0 obj<</Type/Page/MediaBox[0 0 3 3]>>endobj
xref
0 4
0000000000 65535 f
0000000010 00000 n
0000000053 00000 n
0000000102 00000 n
trailer<</Size 4/Root 1 0 R>>
startxref
149
%EOF
`

type mockFileService struct{}

func (m *mockFileService) SaveFile(file filestore.FileItem) error {
	return nil
}

func (m *mockFileService) GetFile(filePath string) (filestore.FileItem, error) {
	return filestore.FileItem{
		FileName:   "test.pdf",
		FolderName: "PATH",
		MimeType:   "application/pdf",
		Payload:    []byte(pdfPayload),
	}, nil
}

func (m *mockFileService) DeleteFile(filePath string) error {
	if filePath == noFileDelete {
		return fmt.Errorf("no file delete")
	}
	return nil
}

func TestGetDocumentByID(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mdr := &mockRepository{}
	svc := &mockFileService{}

	repos := Repositories{
		DocRepo: mdr,
	}

	h := NewHandler(repos, svc)
	c.SetParamNames(ID)
	c.SetParamValues(ID)

	err := h.GetDocumentByID(c)
	if err != nil {
		t.Errorf("cannot get document by id: %v", err)
	}

	assert.Equal(t, http.StatusOK, rec.Code)
	var doc Document
	err = json.Unmarshal(rec.Body.Bytes(), &doc)
	if err != nil {
		t.Errorf(invalidJSON, err)
	}

	// error
	c = e.NewContext(req, rec)

	mdr = &mockRepository{}
	repos = Repositories{
		DocRepo: mdr,
	}
	h = NewHandler(repos, svc)

	err = h.GetDocumentByID(c)
	if err == nil {
		t.Errorf(errExp)
	}
}

func TestDeleteDocumentByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf(fatalErr, err)
	}
	defer db.Close()

	dbx := sqlx.NewDb(db, "mysql")
	con := persistence.NewFromDB(dbx)

	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mdr := &mockRepository{con, false}
	svc := &mockFileService{}

	repos := Repositories{
		DocRepo: mdr,
	}

	h := NewHandler(repos, svc)
	c.SetParamNames(ID)
	c.SetParamValues(ID)

	// straight
	mock.ExpectBegin()
	mock.ExpectCommit()

	err = h.DeleteDocumentByID(c)
	if err != nil {
		t.Errorf("cannot delete document by id: %v", err)
	}

	// start transaction failes
	c = e.NewContext(req, rec)
	failmdr := &mockRepository{con, true}
	faileRepo := Repositories{
		DocRepo: failmdr,
	}
	failH := NewHandler(faileRepo, svc)
	err = failH.DeleteDocumentByID(c)
	if err == nil {
		t.Errorf(errExp)
	}

	// error exists
	mock.ExpectBegin()
	mock.ExpectRollback()

	c = e.NewContext(req, rec)
	c.SetParamNames(ID)
	c.SetParamValues(notExists)
	err = h.DeleteDocumentByID(c)
	if err == nil {
		t.Errorf(errExp)
	}

	// error delete
	mock.ExpectBegin()
	mock.ExpectRollback()

	c = e.NewContext(req, rec)
	c.SetParamNames(ID)
	c.SetParamValues(noDelete)
	err = h.DeleteDocumentByID(c)
	if err == nil {
		t.Errorf(errExp)
	}

	// error no file delete
	mock.ExpectBegin()
	mock.ExpectRollback()

	c = e.NewContext(req, rec)
	c.SetParamNames(ID)
	c.SetParamValues(noFileDelete)
	err = h.DeleteDocumentByID(c)
	if err == nil {
		t.Errorf(errExp)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf(expectations, err)
	}
}

func TestSearchDocuments(t *testing.T) {
	// Setup
	e := echo.New()

	q := make(url.Values)
	q.Set("skip", "0")
	q.Set("limit", "20")
	q.Set("from", time.Now().Format(time.RFC3339))

	req := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mdr := &mockRepository{}
	svc := &mockFileService{}

	repos := Repositories{
		DocRepo: mdr,
	}
	h := NewHandler(repos, svc)

	// success
	err := h.SearchDocuments(c)
	if err != nil {
		t.Errorf("cannot search for documents: %v", err)
		assert.Equal(t, http.StatusOK, rec.Code)
		var pd PagedDcoument
		err = json.Unmarshal(rec.Body.Bytes(), &pd)
		if err != nil {
			t.Errorf("could not unmarshal json: %v", err)
		}
		assert.Equal(t, 2, pd.TotalEntries)
		assert.Equal(t, 2, len(pd.Documents))
	}

	// error
	q = make(url.Values)
	q.Set("skip", "-")
	q.Set("limit", "-")
	q.Set("title", noResult)

	req = httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	err = h.SearchDocuments(c)
	if err == nil {
		t.Errorf(errExp)
	}
}
