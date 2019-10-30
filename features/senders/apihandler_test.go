package senders

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/bihe/mydms/internal/persistence"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

const errNoSenders = "no senders available"
const invalidJSON = "could not get valid json: %v"

// implement persistence.SenderReader
// GetAllSenders() ([]Sender, error)
// SearchSenders(s string) ([]Sender, error)
// SaveSenders(senders []string, a persistence.Atomic) (err error)
type mockRepository struct {
	senders []SenderEntity
}

func (m *mockRepository) init() {
	m.senders = []SenderEntity{
		SenderEntity{ID: 1, Name: "sender1"},
		SenderEntity{ID: 2, Name: "sender2"},
		SenderEntity{ID: 3, Name: "sender3"},
	}
}

func (m *mockRepository) clear() {
	m.senders = []SenderEntity{}
}

func (m *mockRepository) GetAllSenders() ([]SenderEntity, error) {
	if len(m.senders) == 0 {
		return nil, fmt.Errorf(errNoSenders)
	}
	return m.senders, nil
}

func (m *mockRepository) SearchSenders(s string) ([]SenderEntity, error) {
	if len(m.senders) == 0 {
		return nil, fmt.Errorf(errNoSenders)
	}
	filtered := m.senders[:0]
	for _, x := range m.senders {
		if strings.Index(x.Name, s) > -1 {
			filtered = append(filtered, x)
		}
	}

	if len(filtered) == 0 {
		return nil, fmt.Errorf(errNoSenders)
	}

	return filtered, nil
}

func (m *mockRepository) SaveSenders(senders []string, a persistence.Atomic) (err error) {
	return nil
}

func (m *mockRepository) GetSenderByName(name string) (SenderEntity, error) {
	return SenderEntity{}, nil
}

func (m *mockRepository) CreateSender(name string, a persistence.Atomic) (sender SenderEntity, err error) {
	return SenderEntity{}, nil
}

func TestGetAllSenders(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	m := &mockRepository{}
	m.init()
	h := Handler{R: m}
	if assert.NoError(t, h.GetAllSenders(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var senders []Sender
		err := json.Unmarshal(rec.Body.Bytes(), &senders)
		if err != nil {
			t.Errorf(invalidJSON, err)
		}
		assert.Equal(t, 3, len(senders))
		assert.Equal(t, "sender1", senders[0].Name)
		assert.Equal(t, "sender3", senders[len(senders)-1].Name)
	}

	m.clear()
	h = Handler{R: m}
	assert.Error(t, h.GetAllSenders(c))
}

func TestSearchSenders(t *testing.T) {
	// Setup
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	q := req.URL.Query()
	q.Add("name", "sender")
	req.URL.RawQuery = q.Encode()
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	m := &mockRepository{}
	m.init()
	h := Handler{R: m}
	if assert.NoError(t, h.SearchForSenders(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		var senders []Sender
		err := json.Unmarshal(rec.Body.Bytes(), &senders)
		if err != nil {
			t.Errorf(invalidJSON, err)
		}
		assert.Equal(t, 3, len(senders))
		assert.Equal(t, "sender3", senders[len(senders)-1].Name)
		assert.Equal(t, "sender2", senders[len(senders)-2].Name)
	}

	req = httptest.NewRequest(http.MethodGet, "/", nil)
	q = req.URL.Query()
	q.Add("name", "-")
	req.URL.RawQuery = q.Encode()
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	assert.Error(t, h.SearchForSenders(c))
}
