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
var serviceIpAddrResolve string
var servicePortAddr string
var serviceAddress string
var serviceResolvedAddress string
var serviceHttpAddress string

var chEndService = make(chan string)
var helpMap = make(map[string]string)
var autotest_running bool

var useHttps bool
var debug_local_ip net.IP = []byte{192, 168, 1, 187}

func init() {

	succefulAppInit = true

	useHttps = true

	var err error

	serviceIpAddr, err = GetOutboundIP()
	if err != nil {
		errorStr := fmt.Sprintf("IP address recognition error: %s", err)
		fmt.Println(errorStr)
		chEndService <- errorStr
		succefulAppInit = false
		return
	} else {
		if serviceIpAddr.Equal(net.IP(debug_local_ip)) {
			serviceIpAddrResolve = fmt.Sprintf("%s", serviceIpAddr)
		} else {
			serviceIpAddrResolve = "89117771690-apps.site"
		}
	}

	// host, err := net.LookupAddr("50.62.227.1")
	// if err == nil {
	//    fmt.Println(host)
	// }

	if useHttps {
		servicePortAddr = "12443"
	} else {
		servicePortAddr = "12445"
	}

	serviceAddress = fmt.Sprintf("%s:%s", serviceIpAddr, servicePortAddr)
	serviceResolvedAddress = fmt.Sprintf("%s:%s", serviceIpAddrResolve, servicePortAddr)

	if useHttps {
		serviceHttpAddress = fmt.Sprintf("https://%s", serviceResolvedAddress)
	} else {
		serviceHttpAddress = fmt.Sprintf("http://%s", serviceResolvedAddress)
	}

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

	fmt.Printf("Service starting: %s\n", serviceHttpAddress)

	add_http_handler("/", root_handler, fmt.Sprintf("001. %s/,\t\t\t get descriptions of services", serviceHttpAddress))
	add_http_handler("/exit", exit_handler, fmt.Sprintf("002. %s/exit,\t\t\t shutdown service", serviceHttpAddress))
	add_http_handler("/ap", ap_handler, fmt.Sprintf("003. %s/ap?number=3,\t\t arithmetic progression", serviceHttpAddress))
	add_http_handler("/gp2", gp2_handler, fmt.Sprintf("004. %s/gp2?number=3,\t\t geometric progression (with common ratio 2)", serviceHttpAddress))
	add_http_handler("/inc", inc_handler, fmt.Sprintf("005. %s/inc?number=3,\t\t increment", serviceHttpAddress))
	add_http_handler("/autotest_start", autotest_start_handler, fmt.Sprintf("006. %s/autotest_start,\t start autotest", serviceHttpAddress))
	add_http_handler("/autotest_stop", autotest_stop_handler, fmt.Sprintf("007. %s/autotest_stop,\t\t stop autotest", serviceHttpAddress))
	add_http_handler("/strrev", revers_string_handler, fmt.Sprintf("008. %s/strrev?string=\"abc\",\t revers string", serviceHttpAddress))
	add_http_handler("/echo", echo_string_handler, fmt.Sprintf("009. %s/echo?string=\"abc\",\t echo string", serviceHttpAddress))
	add_http_handler("/timestamp", timestamp_handler, fmt.Sprintf("010. %s/timestamp,\t\t unix timestamp", serviceHttpAddress))
	add_http_handler("/lissajous", lissajous_handler, fmt.Sprintf("011. %s/lissajous,\t\t lissajous figures", serviceHttpAddress))
	add_http_handler("/connections_chart", connections_chart_handler, fmt.Sprintf("012. %s/connections_chart,\t connections chart", serviceHttpAddress))
	add_http_handler("/connections_history", connections_history_handler, fmt.Sprintf("013. %s/connections_history,\t connections history", serviceHttpAddress))
	add_http_handler("/particles", particles_handler, fmt.Sprintf("014. %s/particles,\t\t\t particles", serviceHttpAddress))

	// add_http_handler("/", root_handler, fmt.Sprintf("001. %30.30/%-30.30s%30.30s", serviceHttpAddress, " ", "get descriptions of services"))
	// add_http_handler("/exit", exit_handler, fmt.Sprintf("002. %22.22s/%-22.22s%22.22s", serviceHttpAddress, "exit", "shutdown service"))
	// add_http_handler("/ap", ap_handler, fmt.Sprintf("003. %s/ap?number=3,\t\t arithmetic progression", serviceHttpAddress))
	// add_http_handler("/gp2", gp2_handler, fmt.Sprintf("004. %s/gp2?number=3,\t\t geometric progression (with common ratio 2)", serviceHttpAddress))
	// add_http_handler("/inc", inc_handler, fmt.Sprintf("005. %s/inc?number=3,\t\t increment", serviceHttpAddress))
	// add_http_handler("/autotest_start", autotest_start_handler, fmt.Sprintf("006. %s/autotest_start,\t start autotest", serviceHttpAddress))
	// add_http_handler("/autotest_stop", autotest_stop_handler, fmt.Sprintf("007. %s/autotest_stop,\t\t stop autotest", serviceHttpAddress))
	// add_http_handler("/strrev", revers_string_handler, fmt.Sprintf("008. %s/strrev?string=\"abc\",\t revers string", serviceHttpAddress))
	// add_http_handler("/echo", echo_string_handler, fmt.Sprintf("009. %s/echo?string=\"abc\",\t echo string", serviceHttpAddress))
	// add_http_handler("/timestamp", timestamp_handler, fmt.Sprintf("010. %s/timestamp,\t\t unix timestamp", serviceHttpAddress))
	// add_http_handler("/lissajous", lissajous_handler, fmt.Sprintf("011. %s/lissajous,\t\t lissajous figures", serviceHttpAddress))
	// add_http_handler("/connections_history", connections_history_handler, fmt.Sprintf("012. %s/connections_history,\t connections history", serviceHttpAddress))

	// https://pkg.go.dev/net/http#ListenAndServeTLS

	var err error

	if useHttps {
		err = http.ListenAndServeTLS(serviceAddress, "cert.pem", "key.pem", nil)
	} else {
		err = http.ListenAndServe(serviceAddress, nil)
	}

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
