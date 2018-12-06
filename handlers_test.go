package main

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"
)

func TestFbHandler_ServeHTTP_should_handle_valid_request(t *testing.T) {
	var c safeCounter
	c.m = make(map[fbParams]int)
	req := httptest.NewRequest("POST", "/fizzbuzz", strings.NewReader(`{	"limit": 20,
	"int1": 2,
	"int2": 0,
	"string1": "Gregor",
	"string2": "Meow"
}`))
	rr := httptest.NewRecorder()
	handler := http.Handler(fbHandler{&c})
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, `[
    "1",
    "Gregor",
    "3",
    "Gregor",
    "5",
    "Gregor",
    "7",
    "Gregor",
    "9",
    "Gregor",
    "11",
    "Gregor",
    "13",
    "Gregor",
    "15",
    "Gregor",
    "17",
    "Gregor",
    "19",
    "Gregor"
]`, rr.Body.String())
}

func TestFbHandler_ServeHTTP_should_handle_empty_body(t *testing.T) {
	var c safeCounter
	c.m = make(map[fbParams]int)
	req := httptest.NewRequest("POST", "/fizzbuzz", nil)
	rr := httptest.NewRecorder()
	handler := http.Handler(fbHandler{&c})
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestFbHandler_ServeHTTP_limit_too_high(t *testing.T) {
	var c safeCounter
	c.m = make(map[fbParams]int)
	req := httptest.NewRequest("POST", "/fizzbuzz", strings.NewReader(`{	"limit": 900001,
	"int1": 2,
	"int2": 0,
	"string1": "Gregor",
	"string2": "Meow"
}`))
	rr := httptest.NewRecorder()
	handler := http.Handler(fbHandler{&c})
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestStatHandler_ServeHTTP_one_request(t *testing.T) {
	var c safeCounter
	c.m = map[fbParams]int{fbParams{Limit: 10, Int1: 1, Int2: 2, String1: "Gregor", String2: "Meow"}: 5}
	req := httptest.NewRequest("GET", "/statistics", nil)
	rr := httptest.NewRecorder()
	handler := http.Handler(statHandler{&c})
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.JSONEq(t, `[
    {
        "parameters": {
            "limit": 10,
            "int1": 1,
            "int2": 2,
            "string1": "Gregor",
            "string2": "Meow"
        },
        "count": 5
    }
	]`, rr.Body.String())
}

func TestStatHandler_ServeHTTP_no_requests(t *testing.T) {
	var c safeCounter
	c.m = map[fbParams]int{}
	req := httptest.NewRequest("GET", "/statistics", nil)
	rr := httptest.NewRecorder()
	handler := http.Handler(statHandler{&c})
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNoContent, rr.Code)
	assert.Equal(t, "", rr.Body.String())
}

func TestHighestCount_when_two_are_tied_for_highest_returns_both(t *testing.T) {
	var c safeCounter
	c.m = map[fbParams]int{fbParams{Limit: 20, Int1: 1, Int2: 2, String1: "Gregor", String2: "Meow"}: 5,
		fbParams{Limit: 10, Int1: 1, Int2: 2, String1: "Dustin", String2: "Moose"}: 5,
		fbParams{Limit: 10, Int1: 1, Int2: 2, String1: "Moomie", String2: "Cow"}:   3}

	// Goal is to verify that the two elements with highest count of 5 are returned by the function highestCount
	// In order to deal with random map order, resulting slice is ordered by the fbParams Limit
	var slice []fbCount
	slice = c.highestCount()
	sort.Slice(slice, func(i, j int) bool {
		return slice[i].Parameters.Limit > slice[j].Parameters.Limit
	})

	assert.Equal(t, []fbCount{{Parameters: fbParams{20, 1, 2, "Gregor", "Meow"}, Count: 5},
		{Parameters: fbParams{10, 1, 2, "Dustin", "Moose"}, Count: 5}}, slice)
}

func TestSafeCounter_increments_after_posting_to_fbHandler_and_read_by_statHandler(t *testing.T) {
	var c safeCounter
	c.m = make(map[fbParams]int)

	// Make same request to /fizzbuzz 3 times (write to shared map in safeCounter)
	rr := httptest.NewRecorder()
	handler := http.Handler(fbHandler{&c})
	req := httptest.NewRequest("POST", "/fizzbuzz", strings.NewReader(`{	"limit": 20,
	"int1": 2,
	"int2": 0,
	"string1": "Gregor",
	"string2": "Meow"
}`))
	handler.ServeHTTP(rr, req)
	req = httptest.NewRequest("POST", "/fizzbuzz", strings.NewReader(`{	"limit": 20,
	"int1": 2,
	"int2": 0,
	"string1": "Gregor",
	"string2": "Meow"
}`))
	handler.ServeHTTP(rr, req)
	req = httptest.NewRequest("POST", "/fizzbuzz", strings.NewReader(`{	"limit": 20,
	"int1": 2,
	"int2": 0,
	"string1": "Gregor",
	"string2": "Meow"
}`))
	handler.ServeHTTP(rr, req)

	// Check that /statistics endpoint returns count for 3 requests (read from shared map in safeCounter)
	req2 := httptest.NewRequest("GET", "/statistics", nil)
	rr2 := httptest.NewRecorder()
	statHandler := http.Handler(statHandler{&c})
	statHandler.ServeHTTP(rr2, req2)
	assert.Equal(t, http.StatusOK, rr2.Code)
	assert.JSONEq(t, `[
  {
      "parameters": {
          "limit": 20,
          "int1": 2,
          "int2": 0,
          "string1": "Gregor",
          "string2": "Meow"
      },
      "count": 3
  }
	]`, rr2.Body.String())
}
