package main

import (
	"io/ioutil"
	"net/http"
)

func get_url_body(url2get string) string {

	resp, _ := http.Get(url2get)
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	return string(body)
}
