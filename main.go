package main

import (
	"flag"

	"github.com/convexwf/gin-practice/internal"
)

func main() {
	isServer := flag.Bool("server", false, "Run as server")
	isClient := flag.Bool("client", false, "Run as client")
	caFile := flag.String("ca", "", "CA file")
	clientCertFile := flag.String("cert", "", "Cert file for client")
	clientKeyFile := flag.String("key", "", "Key file for client")
	serverAddr := flag.String("addr", "localhost:8080", "Server address")

	flag.Parse()
	if *isServer {
		internal.RunServer()
	} else if *isClient {
		internal.RunClient(*caFile, *clientCertFile, *clientKeyFile, *serverAddr)
	} else {
		flag.PrintDefaults()
	}
}
