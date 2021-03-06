package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

type safeCounter struct {
	m   map[fbParams]int
	mux sync.Mutex
}

func (c *safeCounter) inc(key fbParams) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.m[key]++
}

// Convert map to a slice and isolate the request(s) with the highest count
func (c *safeCounter) highestCount() []fbCount {
	s := c.mapToSlice()
	highestParamsCounts := []fbCount{}
	highestCount := 0
	for _, v := range s {
		switch {
		case v.Count > highestCount:
			highestCount = v.Count
			highestParamsCounts = []fbCount{v}
		case v.Count == highestCount:
			highestParamsCounts = append(highestParamsCounts, v)
		}
	}
	return highestParamsCounts
}

func (c *safeCounter) mapToSlice() []fbCount {
	c.mux.Lock()
	defer c.mux.Unlock()
	var s []fbCount
	for k, v := range c.m {
		s = append(s, fbCount{k, v})
	}
	return s
}

type statHandler struct {
	c *safeCounter
}

func (h statHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slice := h.c.highestCount()
	// If no requests were made to fizzbuzz, return a http 204 No Content
	if len(slice) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	// Else, respond with parameters and number of times requested
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

type fbHandler struct {
	c *safeCounter
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
	// Protect against running out of memory while appending to slice
	if p.Limit > 900000 {
		http.Error(w, "Limit cannot be higher than 900000", http.StatusBadRequest)
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
	h.c.inc(p)
}
