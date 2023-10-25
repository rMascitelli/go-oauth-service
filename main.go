package main

import (
	"log"
)

func main() {
	log.Println("Hello world!")
	pgc := NewPostgresConnector()
	pgc.query_table()
}