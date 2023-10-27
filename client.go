package main

import (
	"net/http"
	"fmt"
)

func main() {
	fmt.Println("Serving at port 8080")
	http.ListenAndServe(":8080", nil)
}