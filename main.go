package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
)

func main() {
	nett := flag.String("net", "tcp", "The class of network")
	host := flag.String("addr", "127.0.0.1", "The host address to listen on")
	port := flag.Int("port", 0, "The port to bind to")
	flag.Parse()

	l, err := net.Listen(*nett, fmt.Sprintf("%s:%d", *host, *port))
	if err == nil {
		fmt.Println(l.Addr())
		err = http.Serve(l, nil)
	}
	if err != nil {
		log.Fatalf("Error: %s\n", err)
	}
}
