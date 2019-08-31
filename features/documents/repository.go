package documents

import (
	"database/sql"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"

	"github.com/bihe/mydms/persistence"
)

// DocumentEntity represents a record in the persistence store
type DocumentEntity struct {
	ID          string         `db:"id"`
	Title       string         `db:"title"`
	FileName    string         `db:"filename"`
	AltID       string         `db:"alternativeid"`
	PreviewLink sql.NullString `db:"previewlink"`
	Amount      float32        `db:"amount"`
	Created     time.Time      `db:"created"`
	Modified    mysql.NullTime `db:"modified"` // go1.13 https://tip.golang.org/pkg/database/sql/#NullTime
	TagList     string         `db:"taglist"`
	SenderList  string         `db:"senderlist"`
}

// PagedDocuments wraps a list of documents and returns the total number of documents
type PagedDocuments struct {
	Documents []DocumentEntity
	Count     int
}

// SortDirection can either by ASC or DESC
type SortDirection uint

const (
	// ASC as ascending sort direction
	ASC SortDirection = iota
	// DESC is descending sort direction
	DESC
)

func (s SortDirection) String() string {
	str := ""
	switch s {
	case ASC:
		str = "ASC"
		break
	case DESC:
		str = "DESC"
		break
	}
	return str
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

// OrderBy is used to sort a result list
type OrderBy struct {
	Field string
	Order SortDirection
}

// Respository is the CRUD interface for documents in the persistence store
type Respository interface {
	Get(id string) (d DocumentEntity, err error)
	Save(doc DocumentEntity, a persistence.Atomic) (d DocumentEntity, err error)
	Delete(id string, a persistence.Atomic) (err error)
	Search(s DocSearch, order []OrderBy) (PagedDocuments, error)
}

// NewRepository creates a new instance using an existing connection
func NewRepository(c persistence.Connection) (Respository, error) {
	if !c.Active {
		return nil, fmt.Errorf("no repository connection available")
	}
	return dbRepository{c}, nil
}

type dbRepository struct {
	c persistence.Connection
}

// Save a document entry. Either create or update the entry, based on availability
// if a valid/active atomic object is supplied the transaction handling is done by the caller
// otherwise a new transaction is created for the scope of the method
func (rw dbRepository) Save(doc DocumentEntity, a persistence.Atomic) (d DocumentEntity, err error) {
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
func (rw dbRepository) Get(id string) (d DocumentEntity, err error) {
	err = rw.c.Get(&d, "SELECT id,title,filename,alternativeid,previewlink,amount,taglist,senderlist,created,modified FROM DOCUMENTS WHERE id=?", id)
	if err != nil {
		err = fmt.Errorf("cannot get document by id '%s': %v", id, err)
		return
	}
	return d, nil
}

// Delete a document by its id
func (rw dbRepository) Delete(id string, a persistence.Atomic) (err error) {
	var (
		atomic *persistence.Atomic
	)

	defer func() {
		err = persistence.HandleTX(!a.Active, atomic, err)
	}()

	if atomic, err = persistence.CheckTX(rw.c, &a); err != nil {
		return
	}

	_, err = atomic.Exec("DELETE FROM DOCUMENTS WHERE id = ?", id)
	if err != nil {
		err = fmt.Errorf("cannot delete document item: %v", err)
	}
	return
}

// Search for documents based on the supplied search-object 'DocSearch'
// the slice of order-bys is used to defined the query sort-order
func (rw dbRepository) Search(s DocSearch, order []OrderBy) (d PagedDocuments, err error) {
	var query string
	q := "SELECT id,title,filename,alternativeid,previewlink,amount,taglist,senderlist,created,modified FROM DOCUMENTS"
	qc := "SELECT count(id) FROM DOCUMENTS"
	where := "\nWHERE 1=1"
	paging := ""
	arg := make(map[string]interface{})

	// use the supplied search-object to create the query
	if s.Title != "" {
		where += "\nAND ( lower(title) LIKE :search OR lower(taglist) LIKE :search OR lower(senderlist) LIKE :search)"
		arg["search"] = "%" + strings.ToLower(s.Title) + "%"
	}
	if s.Tag != "" {
		where += "\nAND lower(taglist) LIKE :tag"
		arg["tag"] = "%" + strings.ToLower(s.Tag) + "%"
	}
	if s.Sender != "" {
		where += "\nAND lower(senderlist) LIKE :sender"
		arg["sender"] = "%" + strings.ToLower(s.Sender) + "%"
	}
	if !s.From.IsZero() {
		where += "\nAND created >= :from"
		arg["from"] = s.From
	}
	if !s.Until.IsZero() {
		where += "\nAND created <= :until"
		arg["until"] = s.Until
	}
	if s.Limit > 0 {
		paging += fmt.Sprintf("\nLIMIT %d", s.Limit)
	}
	if s.Skip > 0 {
		paging += fmt.Sprintf("\nOFFSET %d", s.Skip)
	}

	// get the number of affected documents
	query = qc + where
	var c int
	query, args, err := prepareQuery(rw.c, query, arg)
	if err != nil {
		return
	}

	if err = rw.c.Get(&c, query, args...); err != nil {
		err = fmt.Errorf("could not get the total number of documents: %v", err)
		return
	}

	// query the documents
	orderby := ""
	if order != nil && len(order) > 0 {
		orderby = "\nORDER BY "
		for i, o := range order {
			if i > 0 {
				orderby += ", "
			}
			orderby += fmt.Sprintf("%s %s", o.Field, o.Order)
		}
	}

	// retrieve the documents
	query = q + where + orderby + paging
	log.Debugf("QUERY: %s", query)
	query, args, err = prepareQuery(rw.c, query, arg)
	if err != nil {
		return
	}
	var docs []DocumentEntity
	if err = rw.c.Select(&docs, query, args...); err != nil {
		err = fmt.Errorf("could not get the documents: %v", err)
		return
	}
	return PagedDocuments{Documents: docs, Count: c}, nil
}

func prepareQuery(c persistence.Connection, q string, args map[string]interface{}) (string, []interface{}, error) {
	namedq, namedargs, err := sqlx.Named(q, args)
	if err != nil {
		return "", nil, fmt.Errorf("query error: %v", err)
	}
	query := c.Rebind(namedq)
	return query, namedargs, nil
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
