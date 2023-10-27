package main

import (
	"net/http"
	"fmt"
)

func main() {
	fmt.Println("Serving at port 4321")
	http.ListenAndServe(":4321", nil)
}