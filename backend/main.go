package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

func main() {
	var listenPort int
	flag.IntVar(&listenPort, "listenPort", 0, "port to listen on")
	flag.Parse()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from backend %d", listenPort)
	})

	fmt.Printf("i'm a backend, will listen on port %d\n", listenPort)

	addr := fmt.Sprintf(":%d", listenPort)
	log.Fatal(http.ListenAndServe(addr, handler))
}
