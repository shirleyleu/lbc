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

func (c *safeCounter) Inc(key fbParams) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.m[key]++
}

type fbCount struct {
	Parameters fbParams `json:"parameters"`
	Count      int      `json:"count"`
}

func (c *safeCounter) highestCount() []fbCount {
	c.mux.Lock()
	defer c.mux.Unlock()
	var highestParamsCounts []fbCount
	count := 0
	for k, v := range c.m {
		switch {
		case v > count:
			count = v
			highestParamsCounts = []fbCount{{Parameters: k, Count: v}}
		case v == count:
			count = v
			highestParamsCounts = append(highestParamsCounts, fbCount{Parameters: k, Count: v})
		}
	}
	return highestParamsCounts
}

type statHandler struct {
	c *safeCounter
}

func (h statHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Isolate the request(s) with the highest count
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

type fbHandler struct {
	c *safeCounter
}

func (h fbHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	h.c.Inc(p)
}
