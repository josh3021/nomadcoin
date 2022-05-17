package cli

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/josh3021/nomadcoin/explorer"
	"github.com/josh3021/nomadcoin/rest"
)

func usage() {
	fmt.Printf("Welcome to 노마드 코인\n\n")
	fmt.Printf("Please use the following flags:\n\n")
	fmt.Printf("-restPort:		Sets the \"port\" of the REST API SERVER.\n")
	fmt.Printf("-htmlPort:		Sets the \"port\" of the HTML EXPLORER SERVER.\n")
	fmt.Printf("-mode:		Choose between \"html\" and \"rest\" and \"both\".\n\n")
	runtime.Goexit()
}

func Start() {
	if len(os.Args) == 1 {
		usage()
	}

	// rest := flag.NewFlagSet("rest", flag.ExitOnError)
	// portFlag := rest.Int("port", 4000, "Sets the port of the server")
	restPort := flag.Int("restPort", 4000, "Sets the \"port\" of the REST API SERVER.")
	htmlPort := flag.Int("htmlPort", 3000, "Sets the \"port\" of the HTML EXPLORER SERVER.")
	mode := flag.String("mode", "both", "Sets the \"mode\" of the server.")
	flag.Parse()

	switch *mode {
	case "both":
		go rest.Start(*restPort)
		explorer.Start(*htmlPort)
	case "rest":
		rest.Start(*restPort)
	case "html":
		explorer.Start(*htmlPort)
	default:
		usage()
	}
}
