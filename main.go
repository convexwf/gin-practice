package main

import (
	"flag"

	"github.com/convexwf/gin-practice/internal/cmd"
)

func main() {
	// Define command line flags
	isServer := flag.Bool("server", false, "Run as server")
	isClient := flag.Bool("client", false, "Run as client")

	// Parse command line flags
	flag.Parse()
	if *isServer {
		// Run as server
		cmd.RunServer()
	} else if *isClient {
		// Run as client
		cmd.RunClient()
	} else {
		// Print usage
		flag.Usage()
	}

}
