package main

import (
	"log"
	dbconn "github.com/rMascitelli/go-oauth-service/db_connector"
)

func main() {
	log.Println("Hello world!")
	pgc := dbconn.NewPostgresConnector()
	pgc.query_table()
}