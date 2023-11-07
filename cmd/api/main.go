package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.Println("Hello world!")

	// catch SIGETRM or SIGINTERRUPT
	cancelChan := make(chan os.Signal, 1)
	signal.Notify(cancelChan, syscall.SIGTERM, syscall.SIGINT)

	// Get cmd line args to check for "demo" mode
	demoMode := flag.Bool("demo", false, "Execute in demo mode or not")
	flag.Parse()

	// Consider moving to App() function
	var postgres PostgresConnector
	go func() {
		postgres = NewPostgresConnector(*demoMode)
		rt := NewRouter(5001, postgres)
		rt.StartRouter()
	}()

	sig := <-cancelChan
	log.Printf("Caught signal %v", sig)
	_ = postgres.DropTable(SESSION_TOKENS)
	_ = postgres.DropTable(USER_CREDENTIALS)

}
