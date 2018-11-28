package main

import "fmt"

type fbParams struct {
	Limit   int    `json:"limit"`
	Int1    int    `json:"int1"`
	Int2    int    `json:"int2"`
	String1 string `json:"string1"`
	String2 string `json:"string2"`
}

// Function fizzbuzz returns a list of strings with numbers from 1 to limit where all
// multiples of int1 are replaced by str1, all multiples of int2 are replaced by str2,
// all multiples of int1 and int2 are replaced by str1str2.
func fizzbuzz(p fbParams) []string {
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
