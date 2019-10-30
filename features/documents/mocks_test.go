package documents

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/bihe/mydms/features/filestore"
	"github.com/bihe/mydms/features/senders"
	"github.com/bihe/mydms/features/tags"
	"github.com/bihe/mydms/features/upload"
	"github.com/bihe/mydms/internal/persistence"
)

var errTx = fmt.Errorf("start transaction failed")

// --------------------------------------------------------------------------
// MOCK: documents.Repository
// --------------------------------------------------------------------------

type mockRepository struct {
	c         persistence.Connection
	fail      bool
	errMap    map[int]error
	callCount int
}

func newDocRepo(c persistence.Connection) *mockRepository {
	return &mockRepository{
		c:      c,
		errMap: make(map[int]error),
	}
}

func (m *mockRepository) Get(id string) (d DocumentEntity, err error) {
	m.callCount++
	if id == "" {
		return DocumentEntity{}, fmt.Errorf("no document")
	}
	return DocumentEntity{
		Modified:    sql.NullTime{Time: time.Now().UTC(), Valid: true},
		PreviewLink: sql.NullString{String: "string", Valid: true},
	}, m.errMap[m.callCount]
}

func (m *mockRepository) Save(doc DocumentEntity, a persistence.Atomic) (d DocumentEntity, err error) {
	m.callCount++
	return doc, m.errMap[m.callCount]
}

func (m *mockRepository) Delete(id string, a persistence.Atomic) (err error) {
	m.callCount++
	if id == noDelete {
		return fmt.Errorf("delete error")
	}
	return m.errMap[m.callCount]
}

func (m *mockRepository) Search(s DocSearch, order []OrderBy) (PagedDocuments, error) {
	m.callCount++
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
	m.callCount++
	if id == notExists {
		return "", fmt.Errorf("exists error")
	}
	if id == noFileDelete {
		return noFileDelete, nil
	}
	return "file", nil
}

func (m *mockRepository) CreateAtomic() (persistence.Atomic, error) {
	m.callCount++
	if m.fail {
		return persistence.Atomic{}, errTx
	}
	return m.c.CreateAtomic()
}

func (m *mockRepository) SaveReferences(id string, tagIds, senderIds []int, a persistence.Atomic) (err error) {
	m.callCount++
	return m.errMap[m.callCount]
}

// --------------------------------------------------------------------------
// MOCK: filestore.FileService
// --------------------------------------------------------------------------

type mockFileService struct {
	errMap    map[int]error
	callCount int
}

func newFileService() *mockFileService {
	return &mockFileService{
		errMap: make(map[int]error),
	}
}

// SaveFile(file FileItem) error
// GetFile(filePath string) (FileItem, error)
// DeleteFile(filePath string) error

func (m *mockFileService) SaveFile(file filestore.FileItem) error {
	m.callCount++
	return m.errMap[m.callCount]
}
func (m *mockFileService) GetFile(filePath string) (filestore.FileItem, error) {
	m.callCount++
	return filestore.FileItem{
		FileName:   "test.pdf",
		FolderName: "PATH",
		MimeType:   "application/pdf",
		Payload:    []byte(pdfPayload),
	}, m.errMap[m.callCount]
}
func (m *mockFileService) DeleteFile(filePath string) error {
	m.callCount++
	return m.errMap[m.callCount]
}

// --------------------------------------------------------------------------
// MOCK: tags.Repository
// --------------------------------------------------------------------------

type mockTagsRepository struct {
	errMap    map[int]error
	callCount int
}

func newTagRepo() *mockTagsRepository {
	return &mockTagsRepository{
		errMap: make(map[int]error),
	}
}

// GetAllTags() ([]TagEntity, error)
// SearchTags(s string) ([]TagEntity, error)
// SaveTags(tags []string, a persistence.Atomic) (err error)
// CreateTag(name string, a persistence.Atomic) (tag TagEntity, err error)
// GetTagByName(name string) (TagEntity, error)

func (m *mockTagsRepository) GetAllTags() ([]tags.TagEntity, error) {
	m.callCount++
	return nil, m.errMap[m.callCount]
}
func (m *mockTagsRepository) SearchTags(s string) ([]tags.TagEntity, error) {
	m.callCount++
	return nil, m.errMap[m.callCount]
}
func (m *mockTagsRepository) GetTagByName(name string) (tags.TagEntity, error) {
	m.callCount++
	return tags.TagEntity{}, m.errMap[m.callCount]
}
func (m *mockTagsRepository) SaveTags(tags []string, a persistence.Atomic) (err error) {
	m.callCount++
	return m.errMap[m.callCount]
}
func (m *mockTagsRepository) CreateTag(name string, a persistence.Atomic) (tag tags.TagEntity, err error) {
	m.callCount++
	return tags.TagEntity{}, m.errMap[m.callCount]
}

// --------------------------------------------------------------------------
// MOCK: senders.Repository
// --------------------------------------------------------------------------

type mockSendersRepository struct {
	errMap    map[int]error
	callCount int
}

func newSenderRepo() *mockSendersRepository {
	return &mockSendersRepository{
		errMap: make(map[int]error),
	}
}

// GetAllSenders() ([]SenderEntity, error)
// SearchSenders(s string) ([]SenderEntity, error)
// SaveSenders(senders []string, a persistence.Atomic) (err error)
// CreateSender(name string, a persistence.Atomic) (sender SenderEntity, err error)
// GetSenderByName(name string) (SenderEntity, error)

func (m *mockSendersRepository) GetAllSenders() ([]senders.SenderEntity, error) {
	m.callCount++
	return nil, m.errMap[m.callCount]
}
func (m *mockSendersRepository) SearchSenders(s string) ([]senders.SenderEntity, error) {
	m.callCount++
	return nil, m.errMap[m.callCount]
}
func (m *mockSendersRepository) GetSenderByName(name string) (senders.SenderEntity, error) {
	m.callCount++
	return senders.SenderEntity{}, m.errMap[m.callCount]
}
func (m *mockSendersRepository) SaveSenders(senders []string, a persistence.Atomic) (err error) {
	m.callCount++
	return m.errMap[m.callCount]
}
func (m *mockSendersRepository) CreateSender(name string, a persistence.Atomic) (sender senders.SenderEntity, err error) {
	m.callCount++
	return senders.SenderEntity{}, m.errMap[m.callCount]
}

// --------------------------------------------------------------------------
// MOCK: upload.Repository
// --------------------------------------------------------------------------

type mockUploadRepository struct {
	c         persistence.Connection
	errMap    map[int]error
	resultMap map[int]upload.Upload
	callCount int
}

func newUploadRepo() *mockUploadRepository {
	return &mockUploadRepository{
		errMap:    make(map[int]error),
		resultMap: make(map[int]upload.Upload),
	}
}

// persistence.BaseRepository
// Write(item Upload, a persistence.Atomic) (err error)
// Read(id string) (Upload, error)
// Delete(id string, a persistence.Atomic) (err error)

func (m *mockUploadRepository) Write(item upload.Upload, a persistence.Atomic) (err error) {
	m.callCount++
	return m.errMap[m.callCount]
}
func (m *mockUploadRepository) Read(id string) (upload.Upload, error) {
	m.callCount++
	return m.resultMap[m.callCount], m.errMap[m.callCount]
}
func (m *mockUploadRepository) Delete(id string, a persistence.Atomic) (err error) {
	m.callCount++
	return m.errMap[m.callCount]
}
func (m *mockUploadRepository) CreateAtomic() (persistence.Atomic, error) {
	m.callCount++
	if err := m.errMap[m.callCount]; err != nil {
		return persistence.Atomic{}, errTx
	}
	return m.c.CreateAtomic()
}
