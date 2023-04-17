package main

import (
	"nct/router"
	"os"
)

func main() {
	if len(os.Args) < 3 {
		os.Stderr.WriteString("arg required\n")
		os.Exit(2)
	}

	switch os.Args[1] {
	case "get":
		router.Get()
	}
	// if os.Args[1] == "get" {
	//   router.Get()
	// }
}
