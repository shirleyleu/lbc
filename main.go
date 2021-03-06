package main

import (
	"log"
	"net/http"
)

func main() {
	var c safeCounter
	c.m = make(map[fbParams]int)
	http.Handle("/fizzbuzz", fbHandler{&c})
	http.Handle("/statistics", statHandler{&c})
	log.Print("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
