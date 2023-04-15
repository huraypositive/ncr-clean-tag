package main

import (
  // "fmt"
  "nct/router"
)

func main() {
  // fmt.Println(*headers)
  router.GetRepos(router.MakeApiSpec())
}