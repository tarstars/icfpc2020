package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	//"strconv"
	"strings"

)

func communicate(serverURL, message string) string {
	res, err := http.Post(serverURL, "text/plain", strings.NewReader(message))
	if err != nil {
		log.Printf("Unexpected server response:\n%v", err)
		os.Exit(1)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Unexpected server response:\n%v", err)
		os.Exit(1)
	}

	if res.StatusCode != http.StatusOK {
		log.Printf("Unexpected server response:")
		log.Printf("HTTP code: %d", res.StatusCode)
		log.Printf("Response body: %s", body)
		os.Exit(2)
	}

	return string(body)
}

func main() {
	serverURL := os.Args[1]
	playerKey := os.Args[2]

	log.Printf("ServerUrl: %s; PlayerKey: %s", serverURL, playerKey)

	r1 := communicate(serverURL, playerKey)
	fmt.Println("response1", r1)

	m := "1101000"
	r2 := communicate(serverURL, m)
	fmt.Println("message ", m, " response ", r2)

	/*
		var meter int64
		for meter = 0; meter < 2000; meter += 1 {
			s := strconv.FormatInt(meter, 2) // fmt.Sprintf("%d", meter)
			r2 := communicate(serverURL, s)
			fmt.Println("request: ", s, " response2 ", r2)
		}
	*/
}
