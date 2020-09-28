package store

import (
	"context"
	"log"
	"testing"
)

var testQuery1 = FizzBuzzQuery{
	Str1:  "fizz",
	Str2:  "buzz",
	Int1:  3,
	Int2:  5,
	Limit: 20,
}

var testQuery2 = FizzBuzzQuery{
	Str1:  "azz",
	Str2:  "zoo",
	Int1:  4,
	Int2:  9,
	Limit: 40,
}

func TestInsertAggregate(t *testing.T) {
	s, err := StartTestContainer()
	if err != nil {
		log.Fatalf("Failed to start mongo test container, %s", err.Error())
	}
	err = s.C.Database("app").Collection("queries").Drop(context.Background())
	if err != nil {
		t.Errorf("Failed to drop queries Collection")
	}
	for i := 0; i < 10; i++ {
		err = s.InsertFizzBuzzQuery(context.Background(), testQuery1)
		if err != nil {
			t.Errorf("Failed to insert %v in queries Collection\n", testQuery1)
		}
	}
	for i := 0; i < 20; i++ {
		err = s.InsertFizzBuzzQuery(context.Background(), testQuery2)
		if err != nil {
			t.Errorf("Failed to insert %v in queries Collection\n", testQuery2)
		}
	}
	aggrQueries, err := s.AggregateFizzBuzzQueries(context.Background())
	if err != nil {
		t.Errorf("Failed to aggregate queries Collection\n")
	}
	for _, q := range aggrQueries {
		if q.Str1 == testQuery1.Str1 && q.NbHits != 10 {
			t.Errorf("nbHits for testQuery1 should be 10 not %d\n", q.NbHits)
		}
		if q.Str2 == testQuery2.Str2 && q.NbHits != 20 {
			t.Errorf("nbHits for testQuery2 should be 20 not %d\n", q.NbHits)
		}
	}
}
