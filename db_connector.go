package main

import (
	"log"
	"fmt"
	"database/sql"

	_ "github.com/lib/pq"
)

// Make sure you have created a db with the same 'dbname' as this file
//		CREATE DATABASE dbname;

const (
	CREATE_USER_CRED_TABLE = "CREATE TABLE user_credentials (userid integer, email varchar(255), password varchar(255))"
	CREATE_SESSION_TOKEN_TABLE = "CREATE TABLE session_tokens (token varchar(255), userid varchar(255), expiry_epoch int )"
)

// In case we want to add a different type of DB
type DBConnector interface {
	ConnectToDB() error
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
	    dbname: "oauth_tables",
	    tablename: "example",
	}
	pgc.conninfo = fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=disable", pgc.user, pgc.password, pgc.host, pgc.dbname)
	
	db, err := pgc.ConnectToDB()
	if err != nil {
		log.Fatalf("PostgresConnector: %v\n", err)
	}
	pgc.db = db
	
	if err := pgc.CreateRequiredTables(); err != nil {
		log.Fatalf("PostgresConnector: %v\n", err)
	}

	return pgc
}

func (pgc *PostgresConnector) ConnectToDB() (*sql.DB, error) {
	log.Printf("Connecting to DB %v...\n", pgc.dbname)
    db, err := sql.Open("postgres", pgc.conninfo)
	if err != nil {
		return nil, fmt.Errorf("Couldnt open DB connection, err: %v\n", err)
	}
	return db, nil
}

func (pgc *PostgresConnector) TableExists(tablename string) bool {
    _, err := pgc.db.Exec(fmt.Sprintf("SELECT EXISTS ( SELECT 1 FROM pg_tables WHERE tablename = '%s' ) AS table_existence;", tablename))
    return err == nil
}

func (pgc *PostgresConnector) CreateRequiredTables() error {
	log.Println("Creating required tables...")
	queries := map[string]string{
		"user_credentials": CREATE_USER_CRED_TABLE, 
		"session_tokens": CREATE_SESSION_TOKEN_TABLE,
	}
	for tablename,q := range queries {
		if !pgc.TableExists(tablename) {
			_, err := pgc.db.Exec(q)
			if err != nil {
				return err
			}
		} else {
			log.Printf("	Table [%v] exists\n", tablename)
		}
	}
	return nil
}