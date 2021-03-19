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

func (db *RethinkDB) DropTables() (err error) {
	if err = r.TableDrop("messages").Exec(db.session); err != nil {
		return RethinkDBError{err}
	}

	return nil
}

func (db *RethinkDB) CreateTables() (err error) {
	if err = r.TableCreate("messages").Exec(db.session); err != nil {
		return RethinkDBError{err}
	}

	// function(var61) { return r.row("Sets")("SetID");
	// function(var1) { return r.expr([var1.getField("Sets").getField("SetID")]); },
	if err = r.Table("messages").IndexCreateFunc(
		"SetID",
		// func(row r.Term) interface{} { return interface{}{row.row("Sets").Field("SetID")} },
		// func(row r.Term) interface{} { return []interface{}{row.Field("Sets").Field("SetID")} },
		func(row r.Term) interface{} { return row.Field("Sets").Field("SetID") },
		r.IndexCreateOpts{Multi: true},
	).Exec(db.session); err != nil {
		return RethinkDBError{err}
	}

	return nil
}

func (db *RethinkDB) Write(m *encoding.Message) (err error) {
	var result r.WriteResponse

	if result, err = r.Table("messages").Insert(m).RunWrite(db.session); err != nil {
		return RethinkDBError{err}
	}

	log.Printf("RethinkDB.Write(%v): inserted %v\n", *m, result.Inserted)

	return nil
}
