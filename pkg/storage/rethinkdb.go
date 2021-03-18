package storage

import (
	"fmt"
	"log"

	"github.com/mateumann/clouseau/pkg/encoding"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

type RethinkDB struct {
	address  string
	database string
	session  *r.Session
}

type RethinkDBError struct {
	err error
}

func (e RethinkDBError) Unwrap() error {
	return e.err
}

func (e RethinkDBError) Error() string {
	return fmt.Sprintf("rethinkdb error: %v", e.err)
}

func NewRethinkDB(url string, database string) (db *RethinkDB, err error) {
	var session *r.Session
	session, err = r.Connect(r.ConnectOpts{
		Address:  url, // endpoint without http
		Database: database,
	})

	if err != nil {
		return nil, RethinkDBError{err}
	}

	db = &RethinkDB{address: url, database: database, session: session}

	return db, nil
}

func (db *RethinkDB) Write(m *encoding.Message) (err error) {
	var result r.WriteResponse

	if result, err = r.Table("messages").Insert(m.Marshal()).RunWrite(db.session); err != nil {
		return RethinkDBError{err}
	}

	log.Printf("RethinkDB.Write(%v): inserted %v\n", *m, result.Inserted)

	return nil
}
