package router

import (
	"encoding/json"
	"errors"
	"fmt"
	"nct/config"
	"nct/signature"
	"net/url"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

func Get() error {
	apiSpec := ApiSpec{}
	switch os.Args[2] {
	case "registry":
		flagConfig := (*getFlagConfig("get", "registry")).(*config.GetFlag)
		err := apiSpec.getRegistry(flagConfig)
		if err != nil {
			return err
		}
	case "image", "images":
		if len(os.Args) < 4 || (len(os.Args) >= 4 && string([]rune(os.Args[3])[0]) == "-") {
			flagConfig := (*getFlagConfig("get", "image")).(*config.GetFlag)
			err := apiSpec.getImages(flagConfig)
			if err != nil {
				return err
			}
		} else {
			flagConfig := (*getFlagConfig("get", "image")).(*config.GetFlag)
			err := apiSpec.getImageDetail(flagConfig)
			if err != nil {
				return err
			}
		}
	case "tag", "tags":
		if len(os.Args) < 4 || (len(os.Args) >= 4 && string([]rune(os.Args[3])[0]) == "-") {
			flagConfig := (*getFlagConfig("get", "tag")).(*config.GetFlag)
			err := apiSpec.getTags(flagConfig)
			if err != nil {
				return err
			}
		} else {
			flagConfig := (*getFlagConfig("get", "tag")).(*config.GetFlag)
			err := apiSpec.getTagDetail(flagConfig)
			if err != nil {
				return err
			}
		}
	default:
		fmt.Println(config.GetUsage)
	}
	return nil
}

func (apiSpec *ApiSpec) getRegistry(flagConfig *config.GetFlag) error {
	apiSpec.method = "GET"
	apiSpec.path = "/ncr/api/v2/repositories"
	apiSpec.headers = signature.GetHeader(&apiSpec.method, &apiSpec.path)
	apiSpec.jsonContent = true
	data, _, err := sendRequest(apiSpec)
	if err != nil {
		return err
	}
	if flagConfig.Output == "json" {
		fmt.Println(string(*data))
		return nil
	}
	var body Body
	err = json.Unmarshal(*data, &body)
	if err != nil {
		return err
	}
	if flagConfig.Output == "yaml" {
		yamlString, err := yaml.Marshal(body)
		if err != nil {
			return err
		}
		fmt.Println(string(yamlString))
		return nil
	}
	if !flagConfig.NoHeaders {
		fmt.Println("NAME")
	}
	for _, v := range body.Results {
		fmt.Println(v.(map[string]interface{})["name"])
	}
	return nil
}

func (apiSpec *ApiSpec) getImages(flagConfig *config.GetFlag) error {
	var results Results
	for i := 1; ; i++ {
		apiSpec.method = "GET"
		apiSpec.path = "/ncr/api/v2/repositories/" + flagConfig.Registry + "?page=" + strconv.Itoa(i)
		apiSpec.headers = signature.GetHeader(&apiSpec.method, &apiSpec.path)
		data, _, err := sendRequest(apiSpec)
		if err != nil {
			return err
		}
		var body Body
		err = json.Unmarshal(*data, &body)
		if err != nil {
			return err
		}
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
			return err
		}
		fmt.Println(string(jsonString))
		return nil
	}
	if flagConfig.Output == "yaml" {
		yamlString, err := yaml.Marshal(results)
		if err != nil {
			return err
		}
		fmt.Println(string(yamlString))
		return nil
	}
	if !flagConfig.NoHeaders {
		fmt.Println("NAME")
	}
	for _, v := range results {
		fmt.Println(v.(map[string]interface{})["name"])
	}
	return nil
}

func (apiSpec *ApiSpec) getImageDetail(flagConfig *config.GetFlag) error {
	apiSpec.method = "GET"
	apiSpec.path = "/ncr/api/v2/repositories/" + flagConfig.Registry + "/" + url.QueryEscape(os.Args[3])
	apiSpec.headers = signature.GetHeader(&apiSpec.method, &apiSpec.path)
	data, _, err := sendRequest(apiSpec)
	if err != nil {
		return err
	}
	if flagConfig.Output == "json" {
		fmt.Println(string(*data))
		return nil
	}
	body := make(map[string]interface{})
	err = json.Unmarshal(*data, &body)
	if err != nil {
		return err
	}
	if flagConfig.Output == "yaml" {
		yamlString, err := yaml.Marshal(body)
		if err != nil {
			return err
		}
		fmt.Println(string(yamlString))
		return nil
	}
	if !flagConfig.NoHeaders {
		fmt.Println("NAME")
	}
	fmt.Println(body["name"])
	return nil
}

func (apiSpec *ApiSpec) getTags(flagConfig *config.GetFlag) error {
	if flagConfig.Image == "" {
		return errors.New("You must insert the image name.\n")
	}
	var results Results
	for i := 1; ; i++ {
		apiSpec.method = "GET"
		apiSpec.path = "/ncr/api/v2/repositories/" + flagConfig.Registry + "/" +
			url.QueryEscape(flagConfig.Image) + "/tags?page=" + strconv.Itoa(i)
		apiSpec.headers = signature.GetHeader(&apiSpec.method, &apiSpec.path)
		data, _, err := sendRequest(apiSpec)
		if err != nil {
			return err
		}
		var body Body
		err = json.Unmarshal(*data, &body)
		if err != nil {
			return err
		}
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
			return err
		}
		fmt.Println(string(jsonString))
		return nil
	}
	if flagConfig.Output == "yaml" {
		yamlString, err := yaml.Marshal(results)
		if err != nil {
			return err
		}
		fmt.Println(string(yamlString))
		return nil
	}
	if !flagConfig.NoHeaders {
		fmt.Println("NAME")
	}
	for _, v := range results {
		fmt.Println(v.(map[string]interface{})["name"])
	}
	return nil
}

func (apiSpec *ApiSpec) getTagDetail(flagConfig *config.GetFlag) error {
	apiSpec.method = "GET"
	apiSpec.path = "/ncr/api/v2/repositories/" + flagConfig.Registry + "/" +
		url.QueryEscape(flagConfig.Image) + "/tags/" + os.Args[3]
	apiSpec.headers = signature.GetHeader(&apiSpec.method, &apiSpec.path)
	data, _, err := sendRequest(apiSpec)
	if err != nil {
		return err
	}
	if flagConfig.Output == "json" {
		fmt.Println(string(*data))
		return nil
	}
	body := make(map[string]interface{})
	err = json.Unmarshal(*data, &body)
	if err != nil {
		return err
	}
	if flagConfig.Output == "yaml" {
		yamlString, err := yaml.Marshal(body)
		if err != nil {
			return err
		}
		fmt.Println(string(yamlString))
		return nil
	}
	if !flagConfig.NoHeaders {
		fmt.Println("NAME")
	}
	fmt.Println(body["name"])
	return nil
}
