package main

import (
  "fmt"
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
    return
  }

  // version
  if os.Args[1] == "version" || os.Args[1] == "--version" || os.Args[1] == "-v" {
    fmt.Printf("%s version: %s\n", command, version)
    return
  }

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
