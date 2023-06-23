package main

import (
	"fmt"
	"net/http"
	"sort"
	"time"
)

func http_handler_before_end(resp_writer http.ResponseWriter, r *http.Request, answer string, WriteAnswerToResp bool) {

	add_https_cros_header_for_between_domain_request(resp_writer)

	fmt.Printf("Answer: %s\n", answer)

	if WriteAnswerToResp {
		fmt.Fprintf(resp_writer, "%s", answer)
	}

	add_connect_params_to_db(r, answer)
}

func root_handler(resp_writer http.ResponseWriter, r *http.Request) {

	answer := fmt.Sprintf("Client: %s request: %s\n\n", r.RemoteAddr, r.URL.Path)

	var values []string
	for _, value := range helpMap {
		values = append(values, value)
	}
	sort.Strings(values)
	for _, value := range values {
		answer += fmt.Sprintf("%s\n", value)
	}

	http_handler_before_end(resp_writer, r, answer, true)
}

func ap_handler(resp_writer http.ResponseWriter, r *http.Request) {

	num, err := get_param_number(r)

	var answer string = ""
	if err == nil {
		answer = fmt.Sprintf("%d", ArifmeticProgression(num))
	} else {
		answer = "Bad param"
	}

	http_handler_before_end(resp_writer, r, answer, true)
}

func gp2_handler(resp_writer http.ResponseWriter, r *http.Request) {

	num, err := get_param_number(r)

	var answer string = ""
	if err == nil {
		answer = fmt.Sprintf("%d", GeometricProgression(num, 2))
	} else {
		answer = "Bad param"
	}

	http_handler_before_end(resp_writer, r, answer, true)
}

func inc_handler(resp_writer http.ResponseWriter, r *http.Request) {

	num, err := get_param_number(r)

	var answer string = ""
	if err == nil {
		num++
		answer = fmt.Sprintf("%d", num)
	} else {
		answer = "Bad param"
	}

	http_handler_before_end(resp_writer, r, answer, true)
}

func exit_handler(resp_writer http.ResponseWriter, r *http.Request) {

	answer := fmt.Sprintf("Client: %s request: %s\n", r.RemoteAddr, r.URL.Path)
	answer += "service going to shutting down\n"

	http_handler_before_end(resp_writer, r, answer, true)

	chEndService <- answer
}

func autotest_start_handler(resp_writer http.ResponseWriter, r *http.Request) {

	autotest_running = true

	answer := fmt.Sprintf("Client: %s request: %s\n", r.RemoteAddr, r.URL.Path)
	answer += "autotest mode running...\n"

	go autotest_selfconnect()

	http_handler_before_end(resp_writer, r, answer, true)
}

func autotest_stop_handler(resp_writer http.ResponseWriter, r *http.Request) {

	autotest_running = false

	answer := fmt.Sprintf("Client: %s request: %s\n", r.RemoteAddr, r.URL.Path)
	answer += "autotest mode shutingdown...\n"

	http_handler_before_end(resp_writer, r, answer, true)
}

func revers_string_handler(resp_writer http.ResponseWriter, r *http.Request) {

	stringToRevers := get_param_string(r)

	answer := revers_string(stringToRevers)

	http_handler_before_end(resp_writer, r, answer, true)
}

func echo_string_handler(resp_writer http.ResponseWriter, r *http.Request) {

	answer := get_param_string(r)

	http_handler_before_end(resp_writer, r, answer, true)
}

func timestamp_handler(resp_writer http.ResponseWriter, r *http.Request) {

	nanosec := time.Now().UnixNano()
	answer := fmt.Sprint(nanosec)

	http_handler_before_end(resp_writer, r, answer, true)
}

func lissajous_handler(resp_writer http.ResponseWriter, r *http.Request) {

	lissajous(resp_writer)

	http_handler_before_end(resp_writer, r, "lissajous", false)

}

func connections_history_handler(resp_writer http.ResponseWriter, r *http.Request) {

	answer := get_connections_history()

	http_handler_before_end(resp_writer, r, answer, true)
}

func connections_chart_handler(resp_writer http.ResponseWriter, r *http.Request) {

	drawChart(resp_writer, r)

	http_handler_before_end(resp_writer, r, "history chart", false)
}

// func performance_chart_handler(resp_writer http.ResponseWriter, r *http.Request) {
//
// 	nanosec := time.Now().UnixNano()
// 	answer := fmt.Sprint(nanosec)
//
// 	http_handler_before_end(resp_writer, r, answer, true)
// }
