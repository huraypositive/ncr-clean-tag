package router

import (
	"encoding/json"
	"fmt"
  "nct/config"
  "strconv"
	"nct/signature"
	"os"
)

type Results []map[string]string
type Body struct {
  Next    int         `json:"next"`
	Results []map[string]string `json:"results"`
}

func Get() {
  apiSpec := ApiSpec{}
  switch (os.Args[2]) {
  case "registry":
    flagConfig := getFlagConfig("registry")
    apiSpec.getRepos(flagConfig)
  case "image":
    fallthrough
  case "images":
    if len(os.Args) < 4 || (len(os.Args) >= 4 && string([]rune(os.Args[3])[0]) == "-") {
      flagConfig := getFlagConfig("image")
      apiSpec.getImages(flagConfig)
    } else {
      flagConfig := getFlagConfig("image")
      apiSpec.getImageDetail(flagConfig)
    }
  }
}

func (apiSpec *ApiSpec)getRepos(flagConfig *config.Flag) {
	apiSpec.method = "GET"
	apiSpec.path = "/ncr/api/v2/repositories"
	apiSpec.headers = signature.GetHeader(&apiSpec.method, &apiSpec.path)
	apiSpec.jsonContent = true
	data, err := sendRequest(apiSpec)
	if err != nil {
		return
	}
  if flagConfig.Output == "json" {
    fmt.Println(string(*data))
    return
  }
	var body Body
	err = json.Unmarshal(*data, &body)
  if !flagConfig.NoHeaders {
    fmt.Println("NAME")
  }
	for _, v := range body.Results {
		fmt.Println(v["name"])
	}
}

func (apiSpec *ApiSpec)getImages(flagConfig *config.Flag) {
  var results Results
  for i:=1;;i++ {
    apiSpec.method = "GET"
    apiSpec.path = "/ncr/api/v2/repositories/" + flagConfig.Registry + "?page=" + strconv.Itoa(i)
    apiSpec.headers = signature.GetHeader(&apiSpec.method, &apiSpec.path)
    apiSpec.jsonContent = true
    data, err := sendRequest(apiSpec)
    if err != nil {
      return
    }
    var body Body
    err = json.Unmarshal(*data, &body)
    for _,v := range body.Results {
      results = append(results,v)
    }
    if body.Next == 0 {
      break
    }
  }
  if flagConfig.Output == "json" {
    jsonString, err := json.Marshal(results)
    if err != nil {
      return
    }
    fmt.Println(string(jsonString))
    return
  }
  if !flagConfig.NoHeaders {
    fmt.Println("NAME")
  }
	for _, v := range results {
		fmt.Println(v["name"])
	}
}

func (apiSpec *ApiSpec)getImageDetail(flagConfig *config.Flag) {
  // fmt.Println("arst")
}