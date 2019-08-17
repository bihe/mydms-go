package documents

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/bihe/mydms/persistence"
)

// DocumentEntity represents a record in the persistence store
type DocumentEntity struct {
	ID          int64     `db:"id"`
	Title       string    `db:"title"`
	FileName    string    `db:"filename"`
	AltID       string    `db:"alternativeid"`
	PreviewLink string    `db:"previewlink"`
	Amount      float32   `db:"amount"`
	Created     time.Time `db:"created"`
	Modified    time.Time `db:"modified"`
	TagList     string    `db:"taglist"`
	SenderList  string    `db:"senderlist"`
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
	Get(id int64) (DocumentEntity, error)
	Save(doc DocumentEntity, a persistence.Atomic) (d DocumentEntity, err error)
	Delete(id int64, a persistence.Atomic) (err error)
	Search(s DocSearch) ([]DocumentEntity, error)
}

type dbDocumentReaderWriter struct {
	c persistence.Connection
}

func finshTX(txActive bool, atomic persistence.Atomic, err error) error {
	if !txActive {
		switch err {
		case nil:
			return atomic.Commit()
		default:
			log.Errorf("could not complete the transaction: %v", err)
			if e := atomic.Rollback(); e != nil {
				return fmt.Errorf("%v; could not rollback transaction: %v", err, e)
			}
		}
	}
	return nil
}

func (rw dbDocumentReaderWriter) Save(doc DocumentEntity, a persistence.Atomic) (d DocumentEntity, err error) {
	var atomic *persistence.Atomic

	defer func() {
		err = persistence.HandleTX(!a.Active, atomic, err)
	}()

	if atomic, err = persistence.CheckTX(rw.c, &a); err != nil {
		return
	}

	return DocumentEntity{}, nil
}
