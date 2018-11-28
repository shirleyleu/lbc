package main

import (
	"fmt"
	"log"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Check out the /statistics endpoint")
}

func main() {
	var m = make(map[fbParams]int)
	http.HandleFunc("/", homeHandler)
	http.Handle("/fizzbuzz", fbHandler{m})
	http.Handle("/statistics", statHandler{m})
	log.Print("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
