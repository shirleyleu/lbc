package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_fizzbuzz_3_5_until_100(t *testing.T) {
	arg := fbparams{100, 3, 5, "Fizz", "Buzz"}
	expected := []string{
		"1",
		"2",
		"Fizz",
		"4",
		"Buzz",
		"Fizz",
		"7",
		"8",
		"Fizz",
		"Buzz",
		"11",
		"Fizz",
		"13",
		"14",
		"FizzBuzz",
		"16",
		"17",
		"Fizz",
		"19",
		"Buzz",
		"Fizz",
		"22",
		"23",
		"Fizz",
		"Buzz",
		"26",
		"Fizz",
		"28",
		"29",
		"FizzBuzz",
		"31",
		"32",
		"Fizz",
		"34",
		"Buzz",
		"Fizz",
		"37",
		"38",
		"Fizz",
		"Buzz",
		"41",
		"Fizz",
		"43",
		"44",
		"FizzBuzz",
		"46",
		"47",
		"Fizz",
		"49",
		"Buzz",
		"Fizz",
		"52",
		"53",
		"Fizz",
		"Buzz",
		"56",
		"Fizz",
		"58",
		"59",
		"FizzBuzz",
		"61",
		"62",
		"Fizz",
		"64",
		"Buzz",
		"Fizz",
		"67",
		"68",
		"Fizz",
		"Buzz",
		"71",
		"Fizz",
		"73",
		"74",
		"FizzBuzz",
		"76",
		"77",
		"Fizz",
		"79",
		"Buzz",
		"Fizz",
		"82",
		"83",
		"Fizz",
		"Buzz",
		"86",
		"Fizz",
		"88",
		"89",
		"FizzBuzz",
		"91",
		"92",
		"Fizz",
		"94",
		"Buzz",
		"Fizz",
		"97",
		"98",
		"Fizz",
		"Buzz",
	}
	actual, _ := fizzbuzz(arg)
	assert.Equal(t, expected, actual)
}

func Test_fizzbuzz_0_as_int1(t *testing.T) {
	arg := fbparams{5, 0, 2, "Fizz", "Buzz"}
	expected := []string{"1", "Buzz", "3", "Buzz", "5"}
	actual, _ := fizzbuzz(arg)
	assert.Equal(t, expected, actual)
}

func Test_fizzbuzz_0_as_int2(t *testing.T) {
	arg := fbparams{5, 2, 0, "Fizz", "Buzz"}
	expected := []string{"1", "Fizz", "3", "Fizz", "5"}
	actual, _ := fizzbuzz(arg)
	assert.Equal(t, expected, actual)
}

func Test_fizzbuzz_negative_numbers(t *testing.T) {
	arg := fbparams{10, -2, -5, "Fizz", "Buzz"}
	expected := []string{"1", "Fizz", "3", "Fizz", "Buzz", "Fizz", "7", "Fizz", "9", "FizzBuzz"}
	actual, _ := fizzbuzz(arg)
	assert.Equal(t, expected, actual)
}
