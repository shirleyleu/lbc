package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type statHandler struct {
	c *safeCounter
}

func (c *safeCounter) Inc(key fbParams) {
	c.mux.Lock()
	// Lock so only one goroutine at a time can access the map c.m
	c.m[key]++
	c.mux.Unlock()
}

func (c *safeCounter) highestCount(m map[fbParams]int) []fbCount {
	c.mux.Lock()
	defer c.mux.Unlock()
	var s []fbCount
	for k, v := range m {
		s = append(s, fbCount{k, v})
	}
	highestParamsCounts := []fbCount{}
	count := 0
	for _, v := range s {
		if v.Count >= count {
			count = v.Count
			highestParamsCounts = append(highestParamsCounts, v)
		}
	}
	return highestParamsCounts
}

func (h statHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Convert map to a slice and isolate the request(s) with the highest count
	slice := h.c.highestCount(h.c.m)
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
