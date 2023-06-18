package main

import (
	"net"
	"net/http"
)

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

func add_https_cros_header_for_between_domain_request(resp_writer http.ResponseWriter) {

	// https://stackoverflow.com/questions/39507065/enable-cors-in-golang
	resp_writer.Header().Set("Access-Control-Allow-Origin", "*")
}
