package router

import (
	"bufio"
	"fmt"
	"os"
	"nct/config"
	"nct/signature"
	"net/url"
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
		flagConfig := (*getFlagConfig("delete","image")).(*config.DeleteFlag)
		apiSpec.deleteImage(flagConfig)
	case "tag":
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
			fmt.Println("type yes or no")
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
	body,status,err := sendRequest(apiSpec)
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