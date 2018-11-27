package main

import (
	"fmt"
	"log"
	"net/http"
)

//Print fizz if number is multiple of 3, buzz if number is multiple of 5
//Print fizzbuzz if number is a multiple of both
//Print number if neither

type fizzbuzzParameters struct {
	limit   int
	int1    int
	int2    int
	string1 string
	string2 string
}

func fizzbuzz(p fizzbuzzParameters) []string {
	var s []string
	for i := 1; i <= p.limit; i++ {
		switch {
		case i % (p.int1 * p.int2) == 0:
			s = append(s, p.string1 + p.string2)
		case i % p.int1 == 0:
			s = append(s, p.string1)
		case i % p.int2 == 0:
			s = append(s, p.string2)
		default:
			s = append(s, fmt.Sprintf("%d", i))
		}
	}
	return s
}

type tigerHandler struct{}
func (h tigerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
	//json.Decoder{}
}

func main() {
	http.Handle("/", tigerHandler{})
	log.Print("Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
