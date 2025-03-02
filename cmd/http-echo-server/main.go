package main

import (
	"flag"
	"fmt"
	"github.com/kokardy/go-echo-server/internal/server"
	"os"
)

func main() {
	port := flag.String("port", "10080", "port to run the server on")
	flag.Parse()

	if *port == "" {
		fmt.Println("Port number is required")
		os.Exit(1)
	}

	s := server.New()
	s.Run(":" + *port)
}
