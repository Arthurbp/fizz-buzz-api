package main

import (
	"fizz-buzz-api/pkg/store"
	"fmt"
)

//FizzBuzzResponse generates response from a fizzbuzz query
func FizzBuzzResponse(q store.FizzBuzzQuery) ([]string, error) {
	strList := make([]string, q.Limit)
	for i := range strList {
		strList[i] = fmt.Sprint(1 + i)
		if (i+1)%q.Int1 == 0 && (i+1)%q.Int2 != 0 {
			strList[i] = q.Str1
		}
		if (i+1)%q.Int1 != 0 && (i+1)%q.Int2 == 0 {
			strList[i] = q.Str2
		}
		if (i+1)%q.Int1 == 0 && (i+1)%q.Int2 == 0 {
			strList[i] = q.Str1 + q.Str2
		}
	}
	return strList, nil
}
