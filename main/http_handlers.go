package main

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"
)

func root_handler(resp_writer http.ResponseWriter, r *http.Request) {

	add_https_cros_header_for_between_domain_request(resp_writer)

	var answer string = ""
	var url_path string = r.URL.Path

	answer += fmt.Sprintf("Client: %s request: %s\n\n", r.RemoteAddr, url_path)

	// helpStringCount := len(helpMap)

	var values []string

	for _, value := range helpMap {
		values = append(values, value)
	}

	sort.Strings(values)

	for _, value := range values {
		answer += fmt.Sprintf("%s\n", value)
	}

	// for key, value := range helpMap {
	// 	answer += fmt.Sprintf("%s\n", value)
	// }

	// for c:=0; c < helpStringCount; c++ {
	// 	answer += fmt.Sprintf("%s\n", helpMap[])
	// }

	fmt.Printf("Answer: %s\n", answer)
	fmt.Fprintf(resp_writer, answer)
}

func ap_handler(resp_writer http.ResponseWriter, r *http.Request) {

	add_https_cros_header_for_between_domain_request(resp_writer)

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

	add_https_cros_header_for_between_domain_request(resp_writer)

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

func inc_handler(resp_writer http.ResponseWriter, r *http.Request) {

	add_https_cros_header_for_between_domain_request(resp_writer)

	numStr := r.URL.Query().Get("number")
	num, err := strconv.ParseInt(numStr, 0, 64)

	fmt.Printf("Input param number: %s\n", numStr)

	var answer string = ""
	if err == nil {
		num++
		answer = fmt.Sprintf("%d", num)
	} else {
		answer = "Bad param"
	}

	fmt.Printf("Answer: %s\n", answer)
	fmt.Fprintf(resp_writer, answer)
}

func exit_handler(resp_writer http.ResponseWriter, r *http.Request) {

	add_https_cros_header_for_between_domain_request(resp_writer)

	var answer string = ""
	var url_path string = r.URL.Path

	answer += fmt.Sprintf("Client: %s request: %s\n", r.RemoteAddr, url_path)
	answer += "service going to shutting down\n"

	fmt.Printf("Answer: %s\n", answer)
	fmt.Fprintf(resp_writer, answer)

	chEndService <- "Service shutdown"
}

func autotest_start_handler(resp_writer http.ResponseWriter, r *http.Request) {

	add_https_cros_header_for_between_domain_request(resp_writer)

	var answer string = ""
	var url_path string = r.URL.Path

	answer += fmt.Sprintf("Client: %s request: %s\n", r.RemoteAddr, url_path)
	answer += "autotest mode running...\n"

	fmt.Printf("Answer: %s\n", answer)
	fmt.Fprintf(resp_writer, answer)

	autotest_running = true

	go autotest_selfconnect()
}

func autotest_stop_handler(resp_writer http.ResponseWriter, r *http.Request) {

	add_https_cros_header_for_between_domain_request(resp_writer)

	var answer string = ""
	var url_path string = r.URL.Path

	answer += fmt.Sprintf("Client: %s request: %s\n", r.RemoteAddr, url_path)
	answer += "autotest mode shutingdown...\n"

	fmt.Printf("Answer: %s\n", answer)
	fmt.Fprintf(resp_writer, answer)

	autotest_running = false
}

func revers_string_handler(resp_writer http.ResponseWriter, r *http.Request) {

	add_https_cros_header_for_between_domain_request(resp_writer)

	stringToRevers := r.URL.Query().Get("string")

	answer := revers_string(stringToRevers)

	fmt.Printf("Input param string: %s\n", stringToRevers)

	fmt.Printf("Answer: %s\n", answer)
	fmt.Fprintf(resp_writer, answer)
}

func echo_string_handler(resp_writer http.ResponseWriter, r *http.Request) {

	add_https_cros_header_for_between_domain_request(resp_writer)

	answer := r.URL.Query().Get("string")
	fmt.Printf("Input param string: %s\n", answer)

	fmt.Printf("Answer: %s\n", answer)
	fmt.Fprintf(resp_writer, answer)
}

func timestamp_handler(resp_writer http.ResponseWriter, r *http.Request) {

	add_https_cros_header_for_between_domain_request(resp_writer)

	nanosec := time.Now().UnixNano()

	answer := fmt.Sprint(nanosec)

	fmt.Printf("Answer: %s\n", answer)
	fmt.Fprintf(resp_writer, answer)
}
