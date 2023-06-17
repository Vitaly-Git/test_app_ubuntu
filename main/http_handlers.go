package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func root_handler(resp_writer http.ResponseWriter, r *http.Request) {

	var answer string = ""
	var url_path string = r.URL.Path

	answer += fmt.Sprintf("Client: %s request: %s\n", r.RemoteAddr, url_path)

	is_service_shuting_down := url_path == "/exit"

	if is_service_shuting_down {
		answer += fmt.Sprintf("%s\n", "service going to shutting down")
	}

	fmt.Printf("Answer: %s\n", answer)
	fmt.Fprintf(resp_writer, answer)

	if is_service_shuting_down {
		chEndService <- "Service shutdown"
	}
}

func ap_handler(resp_writer http.ResponseWriter, r *http.Request) {

	numStr := r.URL.Query().Get("number")
	num, err := strconv.ParseInt(numStr, 0, 64)

	fmt.Printf("Input param number: %s\n", numStr)

	var answer string = ""
	if err == nil {
		answer = fmt.Sprintf("%d", ArifmeticProgression(num))
	} else {
		answer = "Bad param"
	}

	fmt.Printf("Answer: %s\n", answer)
	fmt.Fprintf(resp_writer, answer)
}

func gp2_handler(resp_writer http.ResponseWriter, r *http.Request) {

	numStr := r.URL.Query().Get("number")
	num, err := strconv.ParseInt(numStr, 0, 64)

	fmt.Printf("Input param number: %s\n", numStr)

	var answer string = ""
	if err == nil {
		answer = fmt.Sprintf("%d", GeometricProgression(num, 2))
	} else {
		answer = "Bad param"
	}

	fmt.Printf("Answer: %s\n", answer)
	fmt.Fprintf(resp_writer, answer)
}
