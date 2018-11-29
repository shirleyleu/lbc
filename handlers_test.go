package main

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
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

func TestSearchHandler_ServeHTTP_should_handle_empty_body(t *testing.T) {
	var c safeCounter
	c.m = make(map[fbParams]int)
	req := httptest.NewRequest("POST", "/fizzbuzz", nil)
	rr := httptest.NewRecorder()
	handler := http.Handler(fbHandler{&c})
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestStatHandler_ServeHTTP_one_request(t *testing.T) {
	var c safeCounter
	c.m = map[fbParams]int{fbParams{Limit: 10, Int1: 1, Int2: 2, String1: "Gregor", String2: "Meow"}: 5}
	req := httptest.NewRequest("GET", "/fizzbuzz", nil)
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

func TestStatHandler_ServeHTTP_multiple_requests_with_tied_counts(t *testing.T) {
	var c safeCounter
	c.m = map[fbParams]int{fbParams{Limit: 10, Int1: 1, Int2: 2, String1: "Gregor", String2: "Meow"}: 5,
		fbParams{Limit: 10, Int1: 1, Int2: 2, String1: "Dustin", String2: "Moose"}: 5,
		fbParams{Limit: 10, Int1: 1, Int2: 2, String1: "Moomie", String2: "Cow"}:   3}
	req := httptest.NewRequest("GET", "/fizzbuzz", nil)
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
    },
	{
        "parameters": {
            "limit": 10,
            "int1": 1,
            "int2": 2,
            "string1": "Dustin",
            "string2": "Moose"
        },
        "count": 5
    }
	]`, rr.Body.String())
}
