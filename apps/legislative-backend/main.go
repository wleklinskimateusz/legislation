package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/api/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from Go!")
	})

	err := http.ListenAndServe(":8080", nil)

	log.Fatal(err)

}
