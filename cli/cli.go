package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/josh3021/nomadcoin/explorer"
	"github.com/josh3021/nomadcoin/rest"
)

func usage() {
	fmt.Printf("Welcome to 노마드 코인\n\n")
	fmt.Printf("Please use the following flags:\n\n")
	fmt.Printf("-port:		Set the PORT of the server\n")
	fmt.Printf("-mode:		Choose between 'html' and 'rest'\n\n")
	os.Exit(0)
}

func Start() {
	if len(os.Args) == 1 {
		usage()
	}

	// rest := flag.NewFlagSet("rest", flag.ExitOnError)
	// portFlag := rest.Int("port", 4000, "Sets the port of the server")
	port := flag.Int("port", 4000, "Sets the port of the server.")
	mode := flag.String("mode", "rest", "Sets the mode of the server.")
	flag.Parse()

	switch *mode {
	case "rest":
		rest.Start(*port)
	case "html":
		explorer.Start(*port)
	default:
		usage()
	}
}
