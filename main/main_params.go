package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func get_param_number(r *http.Request) (int64, error) {

	numStr := r.URL.Query().Get("number")
	num, err := strconv.ParseInt(numStr, 0, 64)

	fmt.Printf("Input param number: %s\n", numStr)

	return num, err
}

func get_param_string(r *http.Request) string {

	stringParameter := r.URL.Query().Get("string")
	fmt.Printf("Input param string: %s\n", stringParameter)

	return stringParameter
}
