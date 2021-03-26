package storage

import (
	"fmt"
	"log"

	"github.com/mateumann/clouseau/pkg/encoding/ipfix"
	r "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

// RethinkDB represents a rethinkdb database connection.
type RethinkDB struct {
	address  string
	database string
	session  *r.Session
}

// RethinkDBError represents an error that occurs during communication with rethinkdb.
type RethinkDBError struct {
	err error
}

// Unwrap returns the original error returned by rethinkdb.
func (e RethinkDBError) Unwrap() error {
	return e.err
}

func (e RethinkDBError) Error() string {
	return fmt.Sprintf("rethinkdb error: %v", e.err)
}

// NewRethinkDB creates a new connection to the rethinkdb.
func NewRethinkDB(url, database string, recreate bool) (db *RethinkDB, err error) {
	var session *r.Session
	session, err = r.Connect(r.ConnectOpts{
		Address:  url, // endpoint without http
		Database: database,
	})

	if err != nil {
		return nil, RethinkDBError{err}
	}

	db = &RethinkDB{address: url, database: database, session: session}

	if recreate {
		if err = db.dropTables(); err != nil {
			log.Printf("drop tables error: %v", err)
		} else if err = db.createTables(); err != nil {
			log.Printf("create tables error: %v", err)

			return nil, RethinkDBError{err}
		}
	}

	return db, nil
}

func (db RethinkDB) dropTables() (err error) {
	if err = r.TableDrop("messages").Exec(db.session); err != nil {
		return RethinkDBError{err}
	}

	return nil
}

func (db RethinkDB) createTables() (err error) {
	if err = r.TableCreate("messages").Exec(db.session); err != nil {
		return RethinkDBError{err}
	}

	if err = r.Table("messages").IndexCreateFunc(
		"set_id",
		func(row r.Term) interface{} { return row.Field("sets").Field("set_id") },
		r.IndexCreateOpts{Multi: true},
	).Exec(db.session); err != nil {
		return RethinkDBError{err}
	}

	return nil
}

func (db RethinkDB) Write(m *ipfix.Message) (err error) {
	if _, err = r.Table("messages").Insert(m).RunWrite(db.session); err != nil {
		return RethinkDBError{err}
	}

	return nil
}
