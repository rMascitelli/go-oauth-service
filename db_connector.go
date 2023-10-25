package main

import (
	"log"
	"fmt"
	"database/sql"

	_ "github.com/lib/pq"
)

// In case we want to add a different type of DB
type DBConnector interface {
	ConnectToDB() error
	RunQuery(query string)
}

type PostgresConnector struct {
	host string
	port int
	user string
	password string
	dbname string
	tablename string
	conninfo string
	db 	*sql.DB
}

func NewPostgresConnector() PostgresConnector {
	pgc := PostgresConnector{
		host: "localhost",
	    port: 5432,
	    user: "postgres",
	    password: "new_password",
	    dbname: "testdb",
	    tablename: "example",
	}
	pgc.conninfo = fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=disable", pgc.user, pgc.password, pgc.host, pgc.dbname)
	db, err := pgc.ConnectToDB()
	if err != nil {
		log.Fatalf("PostgresConnector: %v\n", err)
	}
	pgc.db = db
	return pgc
}

func (pgc *PostgresConnector) ConnectToDB() (*sql.DB, error) {
    db, err := sql.Open("postgres", pgc.conninfo)
	if err != nil {
		return nil, fmt.Errorf("Couldnt open DB connection, err: %v\n", err)
	}
	return db, nil
}

func (pgc *PostgresConnector) query_table() {
    r, err := pgc.db.Exec(fmt.Sprintf("SELECT * FROM %s;", pgc.tablename))
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("r = %v\n", r)
}