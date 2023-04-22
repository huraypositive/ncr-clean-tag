package main

import (
	"fmt"
	"nct/config"
	"nct/router"
	"os"
)

var command string
var version string

func init() {
	if command == "" {
		command = "nct"
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println(config.Usage)
		os.Exit(2)
	}

	switch os.Args[1] {
	case "--help", "-h":
		fmt.Println(config.Usage)
	case "version", "--version", "-v":
		fmt.Printf("%s version: %s\n", command, version)
	case "get":
		if len(os.Args) < 3 {
			fmt.Println(config.GetUsage)
			return
		}
		router.Get()
	case "delete":
		if len(os.Args) < 3 {
			fmt.Println(config.DeleteUsage)
			return
		}
		router.Delete()
	}
}
