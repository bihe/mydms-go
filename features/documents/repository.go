package documents

import (
	"database/sql"
	"fmt"
	"math/rand"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"github.com/bihe/mydms/persistence"
)

// DocumentEntity represents a record in the persistence store
type DocumentEntity struct {
	ID          string         `db:"id"`
	Title       string         `db:"title"`
	FileName    string         `db:"filename"`
	AltID       string         `db:"alternativeid"`
	PreviewLink string         `db:"previewlink"`
	Amount      float32        `db:"amount"`
	Created     time.Time      `db:"created"`
	Modified    mysql.NullTime `db:"modified"` // go1.13 https://tip.golang.org/pkg/database/sql/#NullTime
	TagList     string         `db:"taglist"`
	SenderList  string         `db:"senderlist"`
}

// DocSearch is used to search for documents
type DocSearch struct {
	Title  string
	Tag    string
	Sender string
	From   time.Time
	Until  time.Time
	Limit  int
	Skip   int
}

// ReaderWriter is the CRUD interface for documents in the persistence store
type ReaderWriter interface {
	Get(id string) (d DocumentEntity, err error)
	Save(doc DocumentEntity, a persistence.Atomic) (d DocumentEntity, err error)
	Delete(id string, a persistence.Atomic) (err error)
	Search(s DocSearch) ([]DocumentEntity, error)
}

type dbDocumentReaderWriter struct {
	c persistence.Connection
}

// Save a document entry. Either create or update the entry, based on availability
// if a valid atomic object is supplied use it for transaction handling - otherwise
// the operation is completed in an atomic manner
func (rw dbDocumentReaderWriter) Save(doc DocumentEntity, a persistence.Atomic) (d DocumentEntity, err error) {
	var (
		atomic  *persistence.Atomic
		newEnty bool
		r       sql.Result
	)

	defer func() {
		err = persistence.HandleTX(!a.Active, atomic, err)
	}()

	if atomic, err = persistence.CheckTX(rw.c, &a); err != nil {
		return
	}

	// try to fetch a document if an ID is supplied
	// the supplied ID is checked against an existing item
	// if the item is not found the provided data is used to create a new entry
	newEnty = true
	if doc.ID != "" {
		var find DocumentEntity
		err = rw.c.Get(&find, "SELECT id,title,filename,alternativeid,previewlink,amount,taglist,senderlist,created,modified FROM DOCUMENTS WHERE id=?", doc.ID)
		if err != nil {
			log.Warnf("could not get a Document by ID '%s' - a new entry will be created", doc.ID)
			newEnty = true
		} else {
			newEnty = false
			doc.Created = find.Created
		}
	}

	if newEnty {
		doc.ID = uuid.New().String()
		doc.Created = time.Now().UTC()
		doc.AltID = randomString(8)
		r, err = atomic.NamedExec("INSERT INTO DOCUMENTS (id,title,filename,alternativeid,previewlink,amount,taglist,senderlist,created) VALUES (:id,:title,:filename,:alternativeid,:previewlink,:amount,:taglist,:senderlist,:created)", &doc)
	} else {
		m := mysql.NullTime{Time: time.Now().UTC(), Valid: true}
		doc.Modified = m
		r, err = atomic.NamedExec("UPDATE DOCUMENTS SET title=:title,filename=:filename,alternativeid=:alternativeid,previewlink=:previewlink,amount=:amount,taglist=:taglist,senderlist=:senderlist,modified=:modified WHERE id=:id", &doc)
	}

	if err != nil {
		err = fmt.Errorf("could not create new entry: %v", err)
		return
	}
	c, err := r.RowsAffected()
	if err != nil {
		err = fmt.Errorf("could not get affected rows: %v", err)
		return
	}
	if c != 1 {
		err = fmt.Errorf("invalid number of rows affected, got %d", c)
		return
	}

	return doc, nil
}

// Get retuns a document by the given id
func (rw dbDocumentReaderWriter) Get(id string) (d DocumentEntity, err error) {
	err = rw.c.Get(&d, "SELECT id,title,filename,alternativeid,previewlink,amount,taglist,senderlist,created,modified FROM DOCUMENTS WHERE id=?", id)
	if err != nil {
		err = fmt.Errorf("cannot get upload-item by id '%s': %v", id, err)
		return
	}
	return d, nil
}

// found: https://www.admfactory.com/how-to-generate-a-fixed-length-random-string-using-golang/
func randomString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}
