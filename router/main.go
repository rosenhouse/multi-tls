package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
)

func main() {
	var listenPort int
	flag.IntVar(&listenPort, "listenPort", 0, "port to listen on")
	flag.Parse()

	reverseProxy := &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			log.Printf("received request for %s\n", r.Host)
			r.URL.Host = strings.Replace(r.Host, "backend-", "127.0.0.1:", -1)
			r.URL.Scheme = "http"
			log.Printf("transformed into request for %s\n", r.URL.String())
		},
	}

	fmt.Printf("i'm the router, will listen on port %d\n", listenPort)

	addr := fmt.Sprintf(":%d", listenPort)
	log.Fatal(http.ListenAndServe(addr, reverseProxy))
}
