package tests

import (
	"database/sql"
	_ "github.com/lib/pq"
	"testing"
)

// TB common interface for *testing.T and *testing.B
type TB interface {
	Fatal(args ...any)
}

func ConnectToPostgres(t TB) (*sql.DB, func()) {
	switch t.(type) {
	case *testing.T, *testing.B:
	default:
		t.Fatal("argument must be instance of *testing.T or *testing.B")
	}

	dataSourceName, err := getDataSourceName()
	if err != nil {
		t.Fatal(err)
	}

	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		t.Fatal(err)
	}

	deferFunc := func() {
		if err = db.Close(); err != nil {
			t.Fatal(err)
		}
	}
	return db, deferFunc
}
