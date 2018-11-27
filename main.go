package main

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
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

func fizzbuzz(p fbparams) ([]string, error) {
	// Reject 0 as
	if p.Int1 == 0 || p.Int2 == 0 {
		return nil, errors.New("No multiples of 0")
	}
	var s []string
	for i := 1; i <= p.Limit; i++ {
		switch {
		case i%(p.Int1*p.Int2) == 0:
			s = append(s, p.String1+p.String2)
		case i%p.Int1 == 0:
			s = append(s, p.String1)
		case i%p.Int2 == 0:
			s = append(s, p.String2)
		default:
			s = append(s, fmt.Sprintf("%d", i))
		}
	}
	return s, nil
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

	// Construct fizzbuzz output
	result, err := fizzbuzz(m)
	if err != nil {
		http.Error(w, "No multiples of 0", http.StatusBadRequest)
		return
	}

	// Send back response as JSON
	resultJSON, err := json.Marshal(result)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error constructing JSON response: %s", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(resultJSON)
	if err != nil {
		log.Printf("Error writing response to header: %s", err)
	}
}

func main() {
	http.Handle("/", tigerHandler{})
	log.Print("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
