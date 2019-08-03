package upload

import (
	"fmt"
	"log"
	"time"

	"github.com/bihe/mydms/persistence"
)

// Upload defines an entity within the persistence store
type Upload struct {
	ID       string    `db:"id"`
	FileName string    `db:"filename"`
	MimeType string    `db:"mimetype"`
	Created  time.Time `db:"created"`
}

// ReaderWriter provides CRUD methods for uploads
type ReaderWriter interface {
	Write(item Upload, a persistence.Atomic) (err error)
	Read(id string) (Upload, error)
	Delete(id string, a persistence.Atomic) (err error)
}

type dbReaderWriter struct {
	c persistence.Connection
}

// NewReaderWriter creates a new instance using an existing connection
func NewReaderWriter(c persistence.Connection) (ReaderWriter, error) {
	if !c.Active {
		return nil, fmt.Errorf("no repository connection available")
	}
	return dbReaderWriter{c}, nil
}

// Write saves an upload item
func (rw dbReaderWriter) Write(item Upload, a persistence.Atomic) (err error) {
	var atomic persistence.Atomic

	defer func() {
		if !a.Active {
			switch err {
			case nil:
				err = atomic.Commit()
			default:
				log.Printf("could not complete the transaction: %v", err)
				if e := atomic.Rollback(); e != nil {
					err = fmt.Errorf("%v; could not rollback transaction: %v", err, e)
				}
			}
		}
	}()

	atomic = a
	if !a.Active {
		atomic, err = rw.c.CreateAtomic()
		if err != nil {
			return
		}
	}
	_, err = atomic.NamedExec("INSERT INTO UPLOADS (id,filename,mimetype,created) VALUES (:id, :filename, :mimetype, :created)", &item)
	if err != nil {
		err = fmt.Errorf("cannot write upload item: %v", err)
		return
	}
	return nil
}

// Read gets an item by it's ID
func (rw dbReaderWriter) Read(id string) (Upload, error) {
	u := Upload{}

	err := rw.c.Get(&u, "SELECT id, filename, mimetype, created FROM UPLOADS WHERE id=?", id)
	if err != nil {
		return Upload{}, fmt.Errorf("cannot get upload-item by id '%s': %v", id, err)
	}
	return u, nil
}

// Delete removes the item with the specified id from the store
func (rw dbReaderWriter) Delete(id string, a persistence.Atomic) (err error) {
	var atomic persistence.Atomic

	defer func() {
		if !a.Active {
			switch err {
			case nil:
				err = atomic.Commit()
			default:
				log.Printf("could not complete the transaction: %v", err)
				if e := atomic.Rollback(); e != nil {
					err = fmt.Errorf("%v; could not rollback transaction: %v", err, e)
				}
			}
		}
	}()

	atomic = a
	if !a.Active {
		atomic, err = rw.c.CreateAtomic()
		if err != nil {
			return
		}
	}
	_, err = atomic.Exec("DELETE FROM UPLOADS WHERE id = ?", id)
	if err != nil {
		err = fmt.Errorf("cannot delete upload item: %v", err)
		return
	}
	return nil
}
