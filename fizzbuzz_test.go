package main

import (
	"fizz-buzz-api/pkg/store"
	"reflect"
	"testing"
)

var testOutput = []string{"1", "2", "fizz", "4", "buzz", "fizz", "7", "8", "fizz", "buzz", "11", "fizz", "13", "14", "fizzbuzz", "16", "17", "fizz", "19", "buzz"}

var testQuery = store.FizzBuzzQuery{
	Str1:  "fizz",
	Str2:  "buzz",
	Int1:  3,
	Int2:  5,
	Limit: 20,
}

// TestFizzBuzzResponse generates test for FizzBuzzResponse
func TestFizzBuzzResponse(t *testing.T) {
	resp, err := FizzBuzzResponse(testQuery)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(resp, testOutput) {
		t.Errorf("Invalid, got %s instead of %s\n", resp, testOutput)
		t.Fail()
	}
}
