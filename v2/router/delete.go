package router

import (
	"bufio"
	"encoding/json"
	"fmt"
	"nct/config"
	"nct/signature"
	"net/url"
	"os"
	"strconv"
)

func Delete() {
	apiSpec := ApiSpec{}
	switch os.Args[2] {
	case "image":
		fallthrough
	case "images":
		if len(os.Args) < 4 || (os.Args[3] != "-h" && os.Args[3] != "--help" && string([]rune(os.Args[3])[0]) == "-") {
			return
		}
		flagConfig := (*getFlagConfig("delete", "image")).(*config.DeleteFlag)
		apiSpec.deleteImage(flagConfig)
	case "tag":
		fallthrough
	case "tags":
		if len(os.Args) < 4 {
			return
		}
		flagConfig := (*getFlagConfig("delete", "tag")).(*config.DeleteFlag)
		apiSpec.deleteTags(flagConfig)
	}
}

func confirmDelete(flagConfig *config.DeleteFlag) {
	if flagConfig.Yes {
		return
	}

	fmt.Println("Confirm Delete (yes/no)")
	var err error
	var confirm string
	stdin := bufio.NewReader(os.Stdin)
	for {
		_, err = fmt.Scanln(&confirm)
		if err != nil {
			stdin.ReadString('\n')
		} else if confirm == "yes" {
			flagConfig.Yes = true
			break
		} else if confirm == "no" {
			break
		} else {
			fmt.Println("Type yes or no")
		}
	}
}

func (apiSpec *ApiSpec) deleteImage(flagConfig *config.DeleteFlag) {
	if confirmDelete(flagConfig); !flagConfig.Yes {
		return
	}
	apiSpec.method = "DELETE"
	apiSpec.path = "/ncr/api/v2/repositories/" + flagConfig.Registry + "/" + url.QueryEscape(os.Args[3])
	apiSpec.headers = signature.GetHeader(&apiSpec.method, &apiSpec.path)
	body, status, err := sendRequest(apiSpec)
	if err != nil {
		return
	}
	if status == 204 {
		fmt.Println(flagConfig.Registry + "/" + os.Args[3] + " image successfully deleted.")
	} else if status == 500 {
		fmt.Println(flagConfig.Registry + "/" + os.Args[3] + " image not exists.")
	} else {
		fmt.Println(string(*body))
	}
}

func (apiSpec *ApiSpec) deleteTags(flagConfig *config.DeleteFlag) {
	if confirmDelete(flagConfig); !flagConfig.Yes {
		return
	}
	apiSpec.method = "DELETE"
	if flagConfig.FileName != "" {
		configs, err := flagConfig.GetConfigFromFile()
		if err != nil {
			return
		}
		for i := range *configs {
			if (*configs)[i].Registry == "" {
				(*configs)[i].Registry = config.DefaultRegistry
			}
			if !(*configs)[i].Enable {
				fmt.Println("not all")
			} else {
				fmt.Println("all")
			}
		}
	}
}

func deleteTagsAll(registry *string, image *string, tag *string) {
	apiSpec := ApiSpec{}
	var results Results
	for i := 1; ; i++ {
		apiSpec.method = "GET"
		apiSpec.path = "/ncr/api/v2/repositories/" + *registry + "/" +
			url.QueryEscape(*image) + "/tags?page=" + strconv.Itoa(i)
		apiSpec.headers = signature.GetHeader(&apiSpec.method, &apiSpec.path)
		data, _, err := sendRequest(&apiSpec)
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
	fmt.Println(results)
}

func deleteTag(apiSpec *ApiSpec, registry *string, image *string, tag *string) {
	apiSpec.path = "/ncr/api/v2/repositories/" + *registry + "/" +
		url.QueryEscape(*image) + "/tags/" + *tag
	apiSpec.headers = signature.GetHeader(&apiSpec.method, &apiSpec.path)
	// fmt.Println(apiSpec)
	// body,status,err := sendRequest(apiSpec)
	// if err != nil {
	// 	return
	// }
	// if status == 204 {
	// 	fmt.Println(flagConfig.Registry + "/" + os.Args[3] + " image successfully deleted.")
	// } else if status == 500 {
	// 	fmt.Println(flagConfig.Registry + "/" + os.Args[3] + " image not exists.")
	// } else {
	// 	fmt.Println(string(*body))
	// }
}
