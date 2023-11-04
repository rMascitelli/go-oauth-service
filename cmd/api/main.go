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
	var pgc PostgresConnector
	go func() {
		pgc = NewPostgresConnector(*demoMode)
		rt := NewRouter(8080, pgc)
		rt.StartRouter()
	}()

	sig := <-cancelChan
	log.Printf("Caught signal %v", sig)
	_ = pgc.DropTable(SESSION_TOKENS)
	_ = pgc.DropTable(USER_CREDENTIALS)

}
