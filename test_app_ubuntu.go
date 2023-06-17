package main

import (
	"fmt"
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
	http.HandleFunc("/ap", ap_handler)
	http.HandleFunc("/gp2", gp2_handler)

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
