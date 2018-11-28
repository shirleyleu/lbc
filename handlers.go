package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
)

type fbHandler struct {
	m map[fbParams]int
}

type statHandler struct {
	m map[fbParams]int
}

func (h statHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Sort map
	slice := sortMap(h.m)

	// Response with parameters and number of times requested
	result, err := json.Marshal(slice)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error constructing JSON response: %s", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(result)
	if err != nil {
		log.Printf("Error writing response to header: %s", err)
		return
	}

}

type fbCount struct {
	Parameters fbParams `json:"parameters"`
	Count      int      `json:"count"`
}

func sortMap(m map[fbParams]int) []fbCount {
	var slice []fbCount
	for k, v := range m {
		slice = append(slice, fbCount{k, v})
	}
	sort.Slice(slice, func(i, j int) bool {
		return slice[i].Count > slice[j].Count
	})
	return slice
}

func (h fbHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Interpret the request
	j, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading request body: %s", err), http.StatusBadRequest)
		return
	}
	var p fbParams
	if err := json.Unmarshal(j, &p); err != nil {
		http.Error(w, fmt.Sprintf("Error parsing JSON: %s", err), http.StatusBadRequest)
		return
	}

	// Construct fizzbuzz output and send back response as JSON
	result, err := json.Marshal(fizzbuzz(p))
	if err != nil {
		http.Error(w, fmt.Sprintf("Error constructing JSON response: %s", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(result)
	if err != nil {
		log.Printf("Error writing response to header: %s", err)
		return
	}

	// Store and count the successful requests
	h.m[p] += 1
}
