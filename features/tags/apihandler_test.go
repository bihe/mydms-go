package tags

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

const noTags = "no tags available"
const invalidJSON = "could not get valid json: %v"

// implement TagReader
// GetAllTags() ([]TagEntity, error)
// SearchTags(s string) ([]TagEntity, error)
type mockTagReader struct {
	tags []TagEntity
}

func (m *mockTagReader) init() {
	m.tags = []TagEntity{
		TagEntity{ID: 1, Name: "tag1"},
		TagEntity{ID: 2, Name: "tag2"},
		TagEntity{ID: 3, Name: "tag3"},
	}
}

func (m *mockTagReader) clear() {
	m.tags = []TagEntity{}
}

func (m *mockTagReader) GetAllTags() ([]TagEntity, error) {
	if len(m.tags) == 0 {
		return nil, fmt.Errorf(noTags)
	}
	return m.tags, nil
}

func (m *mockTagReader) SearchTags(s string) ([]TagEntity, error) {
	if len(m.tags) == 0 {
		return nil, fmt.Errorf(noTags)
	}
	filtered := m.tags[:0]
	for _, x := range m.tags {
		if strings.Index(x.Name, s) > -1 {
			filtered = append(filtered, x)
		}
	}

	if len(filtered) == 0 {
		return nil, fmt.Errorf(noTags)
	}

	return filtered, nil
}

func TestGetAllTags(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	m := &mockTagReader{}
	m.init()
	h := Handler{Reader: m}
	if assert.NoError(t, h.GetAllTags(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var tags []Tag
		err := json.Unmarshal(rec.Body.Bytes(), &tags)
		if err != nil {
			t.Errorf(invalidJSON, err)
		}
		assert.Equal(t, 3, len(tags))
		assert.Equal(t, "tag1", tags[0].Name)
		assert.Equal(t, "tag3", tags[len(tags)-1].Name)
	}

	m.clear()
	h = Handler{Reader: m}
	assert.Error(t, h.GetAllTags(c))
}

func TestSearchTags(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	q := req.URL.Query()
	q.Add("name", "tag")
	req.URL.RawQuery = q.Encode()
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	m := &mockTagReader{}
	m.init()
	h := Handler{Reader: m}
	if assert.NoError(t, h.SearchForTags(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var tags []Tag
		err := json.Unmarshal(rec.Body.Bytes(), &tags)
		if err != nil {
			t.Errorf(invalidJSON, err)
		}
		assert.Equal(t, 3, len(tags))
		assert.Equal(t, "tag1", tags[0].Name)
		assert.Equal(t, "tag3", tags[len(tags)-1].Name)
	}

	req = httptest.NewRequest(http.MethodGet, "/", nil)
	q = req.URL.Query()
	q.Add("name", "-")
	req.URL.RawQuery = q.Encode()
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	assert.Error(t, h.SearchForTags(c))
}
