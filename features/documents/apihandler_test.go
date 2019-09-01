package documents

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bihe/mydms/persistence"
	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

const invalidJSON = "could not get valid json: %v"

/*
type Repository interface {
	Get(id string) (d DocumentEntity, err error)
	Save(doc DocumentEntity, a persistence.Atomic) (d DocumentEntity, err error)
	Delete(id string, a persistence.Atomic) (err error)
	Search(s DocSearch, order []OrderBy) (PagedDocuments, error)
}
*/

type mockRepository struct{}

func (m *mockRepository) Get(id string) (d DocumentEntity, err error) {
	if id == "" {
		return DocumentEntity{}, fmt.Errorf("no document")
	}
	return DocumentEntity{
		Modified:    mysql.NullTime{Time: time.Now().UTC(), Valid: true},
		PreviewLink: sql.NullString{String: "string", Valid: true},
	}, nil
}

func (m *mockRepository) Save(doc DocumentEntity, a persistence.Atomic) (d DocumentEntity, err error) {
	return doc, nil
}

func (m *mockRepository) Delete(id string, a persistence.Atomic) (err error) {
	return nil
}

func (m *mockRepository) Search(s DocSearch, order []OrderBy) (PagedDocuments, error) {
	return PagedDocuments{}, nil
}

func (m *mockRepository) Exists(id string, a persistence.Atomic) (err error) {
	return nil
}

func (m *mockRepository) CreateAtomic() (persistence.Atomic, error) {
	return persistence.Atomic{}, nil
}

func TestGetDocumentByID(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	mdr := &mockRepository{}
	repos := Repositories{
		DocRepo: mdr,
	}

	h := NewHandler(repos)
	c.SetParamNames("id")
	c.SetParamValues("id")

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
	h = NewHandler(repos)

	err = h.GetDocumentByID(c)
	if err == nil {
		t.Errorf("error expected")
	}
}
