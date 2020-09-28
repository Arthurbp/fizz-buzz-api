package main

import (
	"fizz-buzz-api/pkg/store"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

func TestHandlers(t *testing.T) {
	l := log.New(os.Stdout, "", log.Lshortfile)

	s, err := store.StartTestContainer()
	if err != nil {
		log.Fatalf("Failed to start mongo test container: %s\n", err.Error())
	}

	req1, err := http.NewRequest("GET", "/fizzbuzz?str1=value1&str2=value2&int1=1&int2=2&limit=3", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr1 := httptest.NewRecorder()
	handler(s, l).ServeHTTP(rr1, req1)
	if status := rr1.Code; status != http.StatusOK {
		fmt.Println(rr1.Result())
		t.Errorf("/fizzbuzz - handler returned wrong status code: got %v want %v\n",
			status, http.StatusOK)
	}

	req2, err := http.NewRequest("GET", "/fizzbuzz/stats", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr2 := httptest.NewRecorder()
	handler(s, l).ServeHTTP(rr2, req2)
	if status := rr2.Code; status != http.StatusOK {
		fmt.Println(rr2.Result())
		t.Errorf("/fizzbuzz/stats - handler returned wrong status code: got %v want %v\n",
			status, http.StatusOK)
	}
}

func TestHandlersIfNotData(t *testing.T) {
	l := log.New(os.Stdout, "", log.Lshortfile)

	s, err := store.StartTestContainer()
	if err != nil {
		log.Fatalf("Failed to start mongo test container: %s\n", err.Error())
	}

	req2, err := http.NewRequest("GET", "/fizzbuzz/stats", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr2 := httptest.NewRecorder()
	handler(s, l).ServeHTTP(rr2, req2)
	if status := rr2.Code; status != http.StatusNoContent {
		fmt.Println(rr2.Result())
		t.Errorf("/fizzbuzz/stats - handler returned wrong status code: got %v want %v\n",
			status, http.StatusNoContent)
	}
}

func TestFizzBuzzParser(t *testing.T) {

	url := "0.0.0.0/fizzbuzz?str1=value1&str2=value2&int1=1&int2=2&limit=3"
	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("Request issue %s\n", url)
	}
	q, err := parseFizzBuzzParams(r)
	if err != nil {
		t.Fatalf("Malformed test url %s\n", url)
	}
	if !reflect.DeepEqual(q, store.FizzBuzzQuery{Str1: "value1", Str2: "value2", Int1: 1, Int2: 2, Limit: 3}) {
		t.Fatal("Parameters don't match\n")
	}

	// Test empty cases
	url = "0.0.0.0/fizzbuzz?str1=&str2=value2&int1=1&int2=2&limit=3"
	r, err = http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("Request issue %s\n", url)
	}
	_, err = parseFizzBuzzParams(r)
	if err == nil {
		t.Fatalf("Should failed parsing str1 empty ; url: %s\n", url)
	}

	url = "0.0.0.0/fizzbuzz?str1=value1&str2=&int1=1&int2=2&limit=3"
	r, err = http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("Request issue %s\n", url)
	}
	_, err = parseFizzBuzzParams(r)
	if err == nil {
		t.Fatalf("Should failed parsing str2 empty ; url: %s\n", url)
	}

	// Test invalid type value
	url = "0.0.0.0/fizzbuzz?str1=valu1&str2=value2&int1=a&int2=2&limit=3"
	r, err = http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("Request issue %s\n", url)
	}
	_, err = parseFizzBuzzParams(r)
	if err == nil {
		t.Fatalf("Should failed parsing int1 NaN ; url: %s\n", url)
	}

	// Test integer less than 1
	url = "0.0.0.0/fizzbuzz?str1=valu1&str2=value2&int1=1&int2=2&limit=0"
	r, err = http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("Request issue %s\n", url)
	}
	_, err = parseFizzBuzzParams(r)
	if err == nil {
		t.Fatalf("Should failed parsing limit less than 1 ; url: %s\n", url)
	}

	// Test integer greater than 1000000
	url = "0.0.0.0/fizzbuzz?str1=valu1&str2=value2&int1=1&int2=2&limit=1000001"
	r, err = http.NewRequest("GET", url, nil)
	if err != nil {
		t.Fatalf("Request issue %s\n", url)
	}
	_, err = parseFizzBuzzParams(r)
	if err == nil {
		t.Fatalf("Should failed parsing limit greater than 1000000 ; url: %s\n", url)
	}
}
