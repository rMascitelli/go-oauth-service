package main

import (
	"log"
	"fmt"
	"database/sql"
	"errors"
	"time"
	"encoding/hex"	

	_ "github.com/lib/pq"
)

// Make sure you have created a db with the same 'dbname' as this file
//		CREATE DATABASE dbname;

const (
	USER_CREDENTIALS = "user_credentials"
	SESSION_TOKENS = "session_tokens"

	CREATE_USER_CRED_TABLE = "CREATE TABLE user_credentials (userid SERIAL PRIMARY KEY, email varchar(255), password varchar(255));"
	CREATE_SESSION_TOKEN_TABLE = "CREATE TABLE session_tokens (token varchar(255), userid int, expiry_epoch int );"
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
		USER_CREDENTIALS: CREATE_USER_CRED_TABLE, 
		SESSION_TOKENS: CREATE_SESSION_TOKEN_TABLE,
	}
	for tablename,q := range queries {
		if true {//!pgc.TableExists(tablename) {
			_, err := pgc.db.Exec(q)
			if err != nil {
				log.Printf("	%v\n", err)
			} else {
				log.Printf("	Created table [%s]\n", tablename)
			}
		}
	}
	return nil
}

func (pgc *PostgresConnector) QueryUser(email_hash string, password_hash string) (error, UserCredentialRecord) {
	fmt.Printf("Querying [%s]...\n", email_hash[:5])
	var uc UserCredentialRecord
	q := fmt.Sprintf(`SELECT * FROM %s WHERE email='%s'`, USER_CREDENTIALS, email_hash)
	rows, err := pgc.db.Query(q)
	if err != nil {
		fmt.Println("err = ", err)
		return err, UserCredentialRecord{}
	}
	for rows.Next() {
		rows.Scan(&uc.userid, &uc.email, &uc.password)
	}
	if uc.password == password_hash {
		fmt.Println("  Success!\n")
		return nil, uc
	} else {
		return errors.New("Mismatch in password"), uc
	}
}

func (pgc *PostgresConnector) CreateSessionToken(userid int) error {
	now := time.Now().Unix()
	token := hex.EncodeToString(getSHA256Hash(string(now)))	
	log.Printf("Created token - expiry: %d, userid: %d, token: %s", now+300, userid, token[:6])
	q := fmt.Sprintf("INSERT INTO %s (userid, token, expiry_epoch) VALUES ('%d', '%s', '%d')", SESSION_TOKENS, userid, token, now+300)
	_, err := pgc.db.Exec(q)
	if err != nil {
		fmt.Println("err = ", err)
		return err
	}
	return nil
}

func (pgc *PostgresConnector) RegisterUser(user_hash string, password_hash string) error {
	q := fmt.Sprintf("INSERT INTO %s (email, password) VALUES ('%s', '%s')", USER_CREDENTIALS, user_hash, password_hash)
	_, err := pgc.db.Exec(q)
	if err != nil {
		fmt.Println("err = ", err)
		return err
	}
	return nil
}

func (pgc *PostgresConnector) DropTable(tablename string) error {
	log.Printf("Dropping table %s\n", tablename)
	q := fmt.Sprintf("DROP TABLE %s", tablename)
	_, err := pgc.db.Exec(q)
	if err != nil {
		fmt.Println("err = ", err)
		return err
	}
	return nil
}