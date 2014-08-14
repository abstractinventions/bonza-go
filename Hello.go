package main

import (
	"os"
	"io"
	"net/http"
	"log"
	"fmt"
	"strconv"
	"bufio"
	"strings"
)

type Config struct {
	Port int
	Mappings map[string]string

}

func readDotBonza() []string {
	return readFile(".bonza")
}

func readFile(filename string) []string {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println(lines)
	return lines

}

func parseConfig(tokens []string) Config {
	var port, _ = strconv.ParseInt(tokens[0], 10, 16)
	mappings := make(map[string]string)
	for _, value := range tokens[1:] {
		fmt.Println(value)
		parts := strings.SplitN(value, "=>", 2)
		mappings[parts[0]] = parts[1]
	}

	return Config{Port:int(port), Mappings:mappings}
}

func HelloServer(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello, world!\n")
}

func createHandlerFunc(uri string, proxy_url string) func(http.ResponseWriter, *http.Request) {

	client := &http.Client{

	}

	return func(w http.ResponseWriter, req *http.Request) {
		fmt.Println(req)
		requestURI := req.URL.Path

		// ignore non 8 bit strings for now
		proxyRequestUrl := proxy_url + requestURI[len(uri):]

		proxyRequest, _ := http.NewRequest(req.Method, proxyRequestUrl, nil)
		proxyResponse, _ := client.Do(proxyRequest)
		defer proxyResponse.Body.Close()

		//copy headers
		for k, hdr := range proxyResponse.Header {
			for _, v := range hdr {
				w.Header().Add(k, v)
			}
		}
		w.WriteHeader(proxyResponse.StatusCode)
		//copy body
		io.Copy(w, proxyResponse.Body)

	}
}

func main() {

	args := os.Args[1:]
	var config Config
	if (len(args) == 0) {
		config = parseConfig(readDotBonza())
	} else {
		_, err := strconv.ParseInt(args[0], 10, 16)
		if (err != nil) {
			config = parseConfig(readFile(args[0]));
		} else {
			config = parseConfig(args);
		}
	}

	fmt.Println(config)

	for uri, proxy_url := range config.Mappings {
		http.HandleFunc(uri, createHandlerFunc(uri, proxy_url))
	}

	err := http.ListenAndServe((":"+strconv.Itoa(config.Port)), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}


