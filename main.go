package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type fbparams struct {
	Limit   int    `json:"limit"`
	Int1    int    `json:"int1"`
	Int2    int    `json:"int2"`
	String1 string `json:"string1"`
	String2 string `json:"string2"`
}

func fizzbuzz(p fbparams) []string {
	var s []string
	for i := 1; i <= p.Limit; i++ {
		switch {
		// 0 is a valid integer and nothing is a multiple of 0
		case p.Int1 != 0 && p.Int2 != 0 && i%(p.Int1*p.Int2) == 0:
			s = append(s, p.String1+p.String2)
		case p.Int1 != 0 && i%p.Int1 == 0:
			s = append(s, p.String1)
		case p.Int2 != 0 && i%p.Int2 == 0:
			s = append(s, p.String2)
		default:
			s = append(s, fmt.Sprintf("%d", i))
		}
	}
	return s
}

type tigerHandler struct{}

func (h tigerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Interpret the request
	j, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading request body: %s", err), http.StatusBadRequest)
		return
	}
	var m fbparams
	if err := json.Unmarshal(j, &m); err != nil {
		http.Error(w, fmt.Sprintf("Error parsing JSON: %s", err), http.StatusBadRequest)
		return
	}
	// Construct fizzbuzz output and send back response as JSON
	result, err := json.Marshal(fizzbuzz(m))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error constructing JSON response: %s", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(result)
	if err != nil {
		log.Printf("Error writing response to header: %s", err)
	}
}

func main() {
	http.Handle("/", tigerHandler{})
	log.Print("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
