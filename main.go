package main

import (
	"log"
	"os"
	"syscall"
	"os/signal"
)

func main() {
	log.Println("Hello world!")

	// catch SIGETRM or SIGINTERRUPT
	cancelChan := make(chan os.Signal, 1)
    signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)

    // Consider moving to App() function
    go func() {
    	pgc := NewPostgresConnector()
		rt := NewRouter(8080, pgc)
		rt.StartRouter()

    }()

	sig := <-cancelChan
    log.Printf("Caught signal %v", sig)
}