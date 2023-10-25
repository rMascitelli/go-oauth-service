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
	CREATE_USER_CRED_TABLE = "CREATE TABLE user_credentials (userid SERIAL PRIMARY KEY, email varchar(255), password varchar(255));"
	CREATE_SESSION_TOKEN_TABLE = "CREATE TABLE session_tokens (token varcher(255), userid int, expiry_epoch int );"
)

type UserCredentialRecord struct {
	userid int
	email string
	password string
}

type SessionTokenRecord struct {
	token string
	userid int
	expiry_epoch int
}

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
	    password: "root",
	    dbname: "testdb",
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
		if true {//!pgc.TableExists(tablename) {
			_, err := pgc.db.Exec(q)
			if err != nil {
				log.Printf("	Table [%s] exists\n", tablename)
			} else {
				log.Printf("	Created table [%s]\n", tablename)
			}
		}
	}
	return nil
}

func (pgc *PostgresConnector) QueryUser(email_hash string, password_hash string) error {
	fmt.Printf("Querying %s\n", email_hash)
	q := fmt.Sprintf(`SELECT * FROM user_credentials WHERE email='%s'`, email_hash)
	rows, err := pgc.db.Query(q)
	if err != nil {
		fmt.Println("err = ", err)
		return err
	}
	for rows.Next() {
		var uc UserCredentialRecord
		rows.Scan(&uc.userid, &uc.email, &uc.password)
		fmt.Printf("  Found %+v\n", uc)
	}
	return nil
}

func (pgc *PostgresConnector) RegisterUser(user_hash string, password_hash string) error {
	fmt.Printf("Registering %s %s\n", user_hash, password_hash)
	q := fmt.Sprintf("INSERT INTO user_credentials (email, password) VALUES ('%s', '%s')", user_hash, password_hash)
	_, err := pgc.db.Exec(q)
	if err != nil {
		fmt.Println("err = ", err)
		return err
	}
	return nil
}