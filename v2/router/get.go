package router

import (
	"encoding/json"
	"fmt"
	"nct/config"
	"nct/signature"
	"net/url"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

func Get() {
	apiSpec := ApiSpec{}
	switch os.Args[2] {
	case "registry":
		flagConfig := (*getFlagConfig("get", "registry")).(*config.GetFlag)
		apiSpec.getRegistry(flagConfig)
	case "image":
		fallthrough
	case "images":
		if len(os.Args) < 4 || (len(os.Args) >= 4 && string([]rune(os.Args[3])[0]) == "-") {
			flagConfig := (*getFlagConfig("get", "image")).(*config.GetFlag)
			apiSpec.getImages(flagConfig)
		} else {
			flagConfig := (*getFlagConfig("get", "image")).(*config.GetFlag)
			apiSpec.getImageDetail(flagConfig)
		}
	case "tag":
		fallthrough
	case "tags":
		if len(os.Args) < 4 || (len(os.Args) >= 4 && string([]rune(os.Args[3])[0]) == "-") {
			flagConfig := (*getFlagConfig("get", "tag")).(*config.GetFlag)
			apiSpec.getTags(flagConfig)
		} else {
			flagConfig := (*getFlagConfig("get", "tag")).(*config.GetFlag)
			apiSpec.getTagDetail(flagConfig)
		}
	}
}

func (apiSpec *ApiSpec) getRegistry(flagConfig *config.GetFlag) {
	apiSpec.method = "GET"
	apiSpec.path = "/ncr/api/v2/repositories"
	apiSpec.headers = signature.GetHeader(&apiSpec.method, &apiSpec.path)
	apiSpec.jsonContent = true
	data, _, err := sendRequest(apiSpec)
	if err != nil {
		return
	}
	if flagConfig.Output == "json" {
		fmt.Println(string(*data))
		return
	}
	var body Body
	err = json.Unmarshal(*data, &body)
	if err != nil {
		return
	}
	if flagConfig.Output == "yaml" {
		yamlString, err := yaml.Marshal(body)
		if err != nil {
			return
		}
		fmt.Println(string(yamlString))
		return
	}
	if !flagConfig.NoHeaders {
		fmt.Println("NAME")
	}
	for _, v := range body.Results {
		fmt.Println(v.(map[string]interface{})["name"])
	}
}

func (apiSpec *ApiSpec) getImages(flagConfig *config.GetFlag) {
	var results Results
	for i := 1; ; i++ {
		apiSpec.method = "GET"
		apiSpec.path = "/ncr/api/v2/repositories/" + flagConfig.Registry + "?page=" + strconv.Itoa(i)
		apiSpec.headers = signature.GetHeader(&apiSpec.method, &apiSpec.path)
		data, _, err := sendRequest(apiSpec)
		if err != nil {
			return
		}
		var body Body
		err = json.Unmarshal(*data, &body)
		for _, v := range body.Results {
			results = append(results, v)
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
	if flagConfig.Output == "yaml" {
		yamlString, err := yaml.Marshal(results)
		if err != nil {
			return
		}
		fmt.Println(string(yamlString))
		return
	}
	if !flagConfig.NoHeaders {
		fmt.Println("NAME")
	}
	for _, v := range results {
		fmt.Println(v.(map[string]interface{})["name"])
	}
}

func (apiSpec *ApiSpec) getImageDetail(flagConfig *config.GetFlag) {
	apiSpec.method = "GET"
	apiSpec.path = "/ncr/api/v2/repositories/" + flagConfig.Registry + "/" + url.QueryEscape(os.Args[3])
	apiSpec.headers = signature.GetHeader(&apiSpec.method, &apiSpec.path)
	data, _, err := sendRequest(apiSpec)
	if err != nil {
		return
	}
	if flagConfig.Output == "json" {
		fmt.Println(string(*data))
		return
	}
	body := make(map[string]interface{})
	err = json.Unmarshal(*data, &body)
	if err != nil {
		return
	}
	if flagConfig.Output == "yaml" {
		yamlString, err := yaml.Marshal(body)
		if err != nil {
			return
		}
		fmt.Println(string(yamlString))
		return
	}
	if !flagConfig.NoHeaders {
		fmt.Println("NAME")
	}
	fmt.Println(body["name"])
}

func (apiSpec *ApiSpec) getTags(flagConfig *config.GetFlag) {
	var results Results
	for i := 1; ; i++ {
		apiSpec.method = "GET"
		apiSpec.path = "/ncr/api/v2/repositories/" + flagConfig.Registry + "/" +
			url.QueryEscape(flagConfig.Image) + "/tags?page=" + strconv.Itoa(i)
		apiSpec.headers = signature.GetHeader(&apiSpec.method, &apiSpec.path)
		data, _, err := sendRequest(apiSpec)
		if err != nil {
			return
		}
		var body Body
		err = json.Unmarshal(*data, &body)
		for _, v := range body.Results {
			results = append(results, v)
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
	if flagConfig.Output == "yaml" {
		yamlString, err := yaml.Marshal(results)
		if err != nil {
			return
		}
		fmt.Println(string(yamlString))
		return
	}
	if !flagConfig.NoHeaders {
		fmt.Println("NAME")
	}
	for _, v := range results {
		fmt.Println(v.(map[string]interface{})["name"])
	}
}

func (apiSpec *ApiSpec) getTagDetail(flagConfig *config.GetFlag) {
	apiSpec.method = "GET"
	apiSpec.path = "/ncr/api/v2/repositories/" + flagConfig.Registry + "/" +
		url.QueryEscape(flagConfig.Image) + "/tags/" + os.Args[3]
	apiSpec.headers = signature.GetHeader(&apiSpec.method, &apiSpec.path)
	data, _, err := sendRequest(apiSpec)
	if err != nil {
		return
	}
	if flagConfig.Output == "json" {
		fmt.Println(string(*data))
		return
	}
	body := make(map[string]interface{})
	err = json.Unmarshal(*data, &body)
	if err != nil {
		return
	}
	if flagConfig.Output == "yaml" {
		yamlString, err := yaml.Marshal(body)
		if err != nil {
			return
		}
		fmt.Println(string(yamlString))
		return
	}
	if !flagConfig.NoHeaders {
		fmt.Println("NAME")
	}
	fmt.Println(body["name"])
}
