package main

import (
	"fmt"
	"net"
	"net/http"
)

// var serviceAddress = "95.163.241.176:8000"
// var serviceAddress = "95.163.233.63:8000"
// var serviceAddress = "51.250.67.25:8000"

var succefulAppInit bool

var serviceIpAddr net.IP
var servicePortAddr string
var serviceAddress string

var chEndService = make(chan string)
var helpMap = make(map[string]string)
var autotest_running bool

func init() {

	succefulAppInit = true

	var err error

	serviceIpAddr, err = GetOutboundIP()
	if err != nil {
		errorStr := fmt.Sprintf("IP address recognition error: %s", err)
		fmt.Println(errorStr)
		chEndService <- errorStr
		succefulAppInit = false
		return
	}

	servicePortAddr = "8000"

	serviceAddress = fmt.Sprintf("%s:%s", serviceIpAddr, servicePortAddr)
}

func main() {

	if !succefulAppInit {
		fmt.Println("Service stopped")
		return
	}

	fmt.Println("Starting service...")
	go start_http_server()
	waiting_for_gorutines_stops()
	fmt.Println("Service stopped")
}

func start_http_server() {

	fmt.Printf("Service starting on address: %s\n", serviceAddress)

	add_http_handler("/", root_handler, fmt.Sprintf("001. http://%s/,\t\t\t get descriptions of services", serviceAddress))
	add_http_handler("/exit", exit_handler, fmt.Sprintf("002. http://%s/exit,\t\t\t shutdown service", serviceAddress))
	add_http_handler("/ap", ap_handler, fmt.Sprintf("003. http://%s/ap?number=3,\t\t arithmetic progression", serviceAddress))
	add_http_handler("/gp2", gp2_handler, fmt.Sprintf("004. http://%s/gp2?number=3,\t\t geometric progression (with common ratio 2)", serviceAddress))
	add_http_handler("/inc", inc_handler, fmt.Sprintf("005. http://%s/inc?number=3,\t\t increment", serviceAddress))
	add_http_handler("/autotest_start", autotest_start_handler, fmt.Sprintf("006. http://%s/autotest_start,\t\t start autotest", serviceAddress))
	add_http_handler("/autotest_stop", autotest_stop_handler, fmt.Sprintf("007. http://%s/autotest_stop,\t\t stop autotest", serviceAddress))
	add_http_handler("/strrev", revers_string_handler, fmt.Sprintf("008. http://%s/strrev?string=\"abc\",\t revers string", serviceAddress))
	add_http_handler("/echo", echo_string_handler, fmt.Sprintf("009. http://%s/echo?string=\"abc\",\t echo string", serviceAddress))
	add_http_handler("/timestamp", timestamp_handler, fmt.Sprintf("010. http://%s/timestamp,\t\t unix timestamp", serviceAddress))

	err := http.ListenAndServe(serviceAddress, nil)
	if err != nil {
		errorStr := fmt.Sprintf("Error opening listening port: %s", err)
		fmt.Println(errorStr)
		chEndService <- errorStr
		return
	}
}

func add_http_handler(pattern string, handler func(http.ResponseWriter, *http.Request), helpString string) {
	helpMap[pattern] = helpString
	http.HandleFunc(pattern, handler)
}

func waiting_for_gorutines_stops() {
	<-chEndService
}
