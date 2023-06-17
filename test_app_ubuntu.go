package main

import (
	"fmt"
	"net"
	"net/http"
)

// var serviceAddress = "95.163.241.176:8000"
// var serviceAddress = "95.163.233.63:8000"
// var serviceAddress = "51.250.67.25:8000"
var chEndService = make(chan string)

func main() {
	fmt.Println("Starting service...")
	go start_http_server()
	waiting_for_gorutines_stops()
	fmt.Println("Service stoped")
}

func start_http_server() {

	http.HandleFunc("/", root_handler)

	serviceIpAddr, err := GetOutboundIP()
	if err != nil {
		errorStr := fmt.Sprintf("IP address recognition error: %s", err)
		fmt.Println(errorStr)
		chEndService <- errorStr
		return
	}

	servicePortAddr := "8000"
	serviceAddress := fmt.Sprintf("%s:%s", serviceIpAddr, servicePortAddr)
	fmt.Printf("Service starting on address: %s\n", serviceAddress)

	err = http.ListenAndServe(serviceAddress, nil)
	if err != nil {
		errorStr := fmt.Sprintf("Error opening listening port: %s", err)
		fmt.Println(errorStr)
		chEndService <- errorStr
		return
	}
}

func waiting_for_gorutines_stops() {
	<-chEndService
}

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

func get_local_ip() (net.IP, error) {

	var ip net.IP
	var err error

	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, i := range ifaces {

		var addrs []net.Addr

		addrs, err = i.Addrs()
		if err != nil {
			return nil, err
		}

		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			default:
				continue
			}
			//fmt.Println(ip)

			return ip, nil
		}
	}

	return ip, nil
}

// Get preferred outbound ip of this machine
func GetOutboundIP() (net.IP, error) {

	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, nil
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP, nil
}
