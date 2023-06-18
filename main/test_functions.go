package main

import (
	"fmt"
	"strconv"
)

func autotest_selfconnect() {

	var siteRequest string = ""
	var testNumber int64 = 2
	//var err error

	for autotest_running {
		siteRequest = fmt.Sprintf("http://%s/inc?number=%d", serviceAddress, testNumber)
		numberStr := get_url_body(siteRequest)
		testNumber, _ = strconv.ParseInt(numberStr, 0, 64)
	}

}
