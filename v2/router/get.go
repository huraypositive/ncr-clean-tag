package router

import (
	"encoding/json"
	"fmt"
	"nct/signature"
  "nct/config"
	"os"
)

func Get() {
  apiSpec := ApiSpec{}
  switch (os.Args[2]) {
  case "registry":
    apiSpec.getRepos()
  case "image":
    fallthrough
  case "images":
    // if os.Args[3] != 
    if len(os.Args) >= 4 && string([]rune(os.Args[3])[0]) == "-" {
      apiSpec.getImages()
    } else {
      apiSpec.getImageDetail()
    }
  }
}

func (apiSpec *ApiSpec)getRepos() {
	apiSpec.method = "GET"
	apiSpec.path = "/ncr/api/v2/repositories"
	apiSpec.headers = signature.GetHeader(&apiSpec.method, &apiSpec.path)
	apiSpec.jsonContent = true
	data, err := sendRequest(apiSpec)
	if err != nil {
		return
	}
	type Body struct {
		Results []map[string]string `json:"results"`
	}
	var body Body
	err = json.Unmarshal(*data, &body)
	fmt.Println("NAME")
	for _, v := range body.Results {
		fmt.Println(v["name"])
	}
}

func (apiSpec *ApiSpec)getImages() {
  flagConfig := getFlagConfig()
	apiSpec.method = "GET"
	apiSpec.path = "/ncr/api/v2/repositories/" + flagConfig.Registry
	apiSpec.headers = signature.GetHeader(&apiSpec.method, &apiSpec.path)
	apiSpec.jsonContent = true
	data, err := sendRequest(apiSpec)
	if err != nil {
		return
	}
  fmt.Println(string(*data))
	// type Body struct {
	// 	Results []map[string]string `json:"results"`
	// }
	// var body Body
	// err = json.Unmarshal(*data, &body)
	// fmt.Println("NAME")
	// for _, v := range body.Results {
	// 	fmt.Println(v["name"])
	// }
}

func (apiSpec *ApiSpec)getImageDetail() {
  // fmt.Println("arst")
}

func getFlagConfig() *config.Flag {
  flagConfig := config.Flag{}
  flags := flagConfig.Setup()
  flagConfig.Parse(flags,os.Args[2:])
  return &flagConfig
}