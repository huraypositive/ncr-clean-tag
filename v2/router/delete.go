package router

import (
	"bufio"
	"encoding/json"
	"fmt"
	"nct/config"
	"nct/signature"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
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
	case "tag", "tags":
		if len(os.Args) < 4 {
			return
		}
		flagConfig := (*getFlagConfig("delete", "tag")).(*config.DeleteFlag)
		apiSpec.deleteTags(flagConfig)
	default:
		fmt.Println(config.DeleteUsage)
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
	if flagConfig.DryRun {
		fmt.Println(flagConfig.Registry + "/" + os.Args[3] + " image deleted. - Dry run")
		return
	}

	apiSpec.method = "DELETE"
	apiSpec.path = "/ncr/api/v2/repositories/" + flagConfig.Registry + "/" + url.QueryEscape(os.Args[3])
	apiSpec.headers = signature.GetHeader(&apiSpec.method, &apiSpec.path)
	body, status, err := sendRequest(apiSpec)
	if err != nil {
		os.Stderr.WriteString(err.Error())
		return
	}
	if status == 204 {
		fmt.Println(flagConfig.Registry + "/" + os.Args[3] + " image successfully deleted.")
	} else if status == 500 {
		os.Stderr.WriteString(flagConfig.Registry + "/" + os.Args[3] + " image not exists.\n")
	} else {
		os.Stderr.WriteString(string(*body))
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
			if (*configs)[i].Recent < 0 {
				os.Stderr.WriteString("exclude-recent must be greater than or equal to zero. (Image:" + (*configs)[i].Image + ")\n")
				continue
			}
			if (*configs)[i].Registry == "" {
				(*configs)[i].Registry = config.DefaultRegistry
			}
			if flagConfig.DryRun {
				(*configs)[i].DryRun = flagConfig.DryRun
			}
			if !(*configs)[i].Enable {
				for j := range (*configs)[i].Tags {
					deleteTag(apiSpec, &(*configs)[i].DryRun, &(*configs)[i].Registry, &(*configs)[i].Image, (*configs)[i].Tags[j])
				}
			} else {
				results := getDeleteList(&(*configs)[i].Registry, &(*configs)[i].Image, &(*configs)[i].Recent)
				for j := 0 + (*configs)[i].Recent; j < len(results); j++ {
					deleteTag(apiSpec, &(*configs)[i].DryRun, &(*configs)[i].Registry, &(*configs)[i].Image, fmt.Sprintf("%s", results[j].(map[string]interface{})["name"]))
				}
			}
		}
	} else {
		if flagConfig.Image == "" {
			os.Stderr.WriteString("You must insert the image name.\n")
			return
		}
		if !flagConfig.Enable {
			index := len(os.Args[3:])
			for i := 0; i < len(os.Args[3:]); i++ {
				if string([]rune(os.Args[i+3])[0]) == "-" {
					index = i
					break
				}
			}
			tags := make([]string, len(os.Args[3:index+3]))
			for i, v := range os.Args[3 : index+3] {
				tags[i] = v
			}
			for i := range tags {
				deleteTag(apiSpec, &flagConfig.DryRun, &flagConfig.Registry, &flagConfig.Image, tags[i])
			}
		} else {
			results := getDeleteList(&flagConfig.Registry, &flagConfig.Image, &flagConfig.Recent)
			for j := 0 + flagConfig.Recent; j < len(results); j++ {
				deleteTag(apiSpec, &flagConfig.DryRun, &flagConfig.Registry, &flagConfig.Image, fmt.Sprintf("%s", results[j].(map[string]interface{})["name"]))
			}
		}
	}
}

func getDeleteList(registry *string, image *string, recent *int) Results {
	apiSpec := ApiSpec{}
	apiSpec.method = "GET"
	var results Results
	for i := 1; ; i++ {
		apiSpec.path = "/ncr/api/v2/repositories/" + *registry + "/" +
			url.QueryEscape(*image) + "/tags?page=" + strconv.Itoa(i)
		apiSpec.headers = signature.GetHeader(&apiSpec.method, &apiSpec.path)
		data, _, err := sendRequest(&apiSpec)
		if err != nil {
			os.Stderr.WriteString(err.Error())
			continue
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
	if len(results) <= *recent {
		os.Stderr.WriteString("The " + *image + " image contains fewer tags than exclude-recent. - Skipping\n")
		return nil
	}
	sort.Slice(results, func(i, j int) bool {
		now, _ := strconv.Atoi(strings.Split(fmt.Sprintf("%f", results[i].(map[string]interface{})["last_updated"]), ".")[0])
		next, _ := strconv.Atoi(strings.Split(fmt.Sprintf("%f", results[j].(map[string]interface{})["last_updated"]), ".")[0])
		return now > next
	})
	return results
}

func deleteTag(apiSpec *ApiSpec, dryRun *bool, registry *string, image *string, tag string) {
	if *dryRun {
		fmt.Println(*registry + "/" + *image + ":" + tag + " deleted. - Dry run")
		return
	}
	apiSpec.path = "/ncr/api/v2/repositories/" + *registry + "/" +
		url.QueryEscape(*image) + "/tags/" + tag
	apiSpec.headers = signature.GetHeader(&apiSpec.method, &apiSpec.path)
	_, status, err := sendRequest(apiSpec)
	if err != nil {
		os.Stderr.WriteString(err.Error())
		return
	}
	if status == 204 {
		fmt.Println(*registry + "/" + *image + ":" + tag + " successfully deleted.")
	} else {
		os.Stderr.WriteString(*registry + "/" + *image + ":" + tag + " not exists.\n")
	}
}
