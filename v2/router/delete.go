package router

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"nct/config"
	"nct/signature"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
)

func Delete() error {
	apiSpec := ApiSpec{}
	switch os.Args[2] {
	case "image", "images":
		if len(os.Args) < 4 || (os.Args[3] != "-h" && os.Args[3] != "--help" && string([]rune(os.Args[3])[0]) == "-") {
			fmt.Println(config.DeleteUsage)
			return nil
		}
		flagConfig := (*getFlagConfig("delete", "image")).(*config.DeleteFlag)
		err := apiSpec.deleteImage(flagConfig)
		if err != nil {
			return err
		}
	case "tag", "tags":
		if len(os.Args) < 4 {
			fmt.Println(config.DeleteUsage)
			return nil
		}
		flagConfig := (*getFlagConfig("delete", "tag")).(*config.DeleteFlag)
		err := apiSpec.deleteTags(flagConfig)
		if err != nil {
			return err
		}
	default:
		fmt.Println(config.DeleteUsage)
	}
	return nil
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

func (apiSpec *ApiSpec) deleteImage(flagConfig *config.DeleteFlag) error {
	if confirmDelete(flagConfig); !flagConfig.Yes {
		return nil
	}
	if flagConfig.DryRun {
		fmt.Println(flagConfig.Registry + "/" + os.Args[3] + " image deleted. - Dry run")
		return nil
	}

	apiSpec.method = "DELETE"
	apiSpec.path = "/ncr/api/v2/repositories/" + flagConfig.Registry + "/" + url.QueryEscape(os.Args[3])
	apiSpec.headers = signature.GetHeader(&apiSpec.method, &apiSpec.path)
	body, status, err := sendRequest(apiSpec)
	if err != nil {
		os.Stderr.WriteString(err.Error())
		return err
	}
	if status == 204 {
		fmt.Println(flagConfig.Registry + "/" + os.Args[3] + " image successfully deleted.")
		return nil
	} else if status == 500 {
		return errors.New(flagConfig.Registry + "/" + os.Args[3] + " image not exists.\n")
	} else {
		return errors.New(string(*body))
	}
}

func (apiSpec *ApiSpec) deleteTags(flagConfig *config.DeleteFlag) error {
	if confirmDelete(flagConfig); !flagConfig.Yes {
		return nil
	}
	apiSpec.method = "DELETE"
	if flagConfig.FileName != "" {
		configs, err := flagConfig.GetConfigFromFile()
		if err != nil {
			return err
		}
		for i := range *configs {
			if (*configs)[i].ExcludeRecent < 0 {
				os.Stderr.WriteString("exclude-recent must be greater than or equal to zero. (Image:" + (*configs)[i].Image + ")\n")
				continue
			}
			if (*configs)[i].Registry == "" {
				if flagConfig.Registry != "" {
					(*configs)[i].Registry = flagConfig.Registry
				} else {
					(*configs)[i].Registry = config.DefaultRegistry
				}
			}
			if len((*configs)[i].ExcludeTags) == 0 && len(flagConfig.ExcludeTags) != 0 {
				(*configs)[i].ExcludeTags = flagConfig.ExcludeTags
			}
			if (*configs)[i].ExcludeRecent == 0 && flagConfig.ExcludeRecent != 0 {
				(*configs)[i].ExcludeRecent = flagConfig.ExcludeRecent
			}
			if flagConfig.DryRun {
				(*configs)[i].DryRun = flagConfig.DryRun
			}
			if !(*configs)[i].Enable {
				if len((*configs)[i].Tags) == 0 {
					return errors.New("There is no list of delete tags. To delete all, use the 'all: true' option.\n")
				}
				for j := range (*configs)[i].Tags {
					err := deleteTag(apiSpec, &(*configs)[i].DryRun, &(*configs)[i].Registry, &(*configs)[i].Image, (*configs)[i].Tags[j], &(*configs)[i].ExcludeTags)
					if err != nil {
						return err
					}
				}
			} else {
				results, err := getDeleteList(&(*configs)[i].Registry, &(*configs)[i].Image, &(*configs)[i].ExcludeRecent)
				if err != nil {
					return err
				}
				for j := 0 + (*configs)[i].ExcludeRecent; j < len(results); j++ {
					err := deleteTag(apiSpec, &(*configs)[i].DryRun, &(*configs)[i].Registry, &(*configs)[i].Image, fmt.Sprintf("%s", results[j].(map[string]interface{})["name"]), &(*configs)[i].ExcludeTags)
					if err != nil {
						return err
					}
				}
			}
		}
	} else {
		if flagConfig.Image == "" {
			return errors.New("You must insert the image name.\n")
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
			if len(tags) == 0 {
				return errors.New("There is no list of delete tags. To delete all, use the '--all' option.\n")
			}
			for i := range tags {
				err := deleteTag(apiSpec, &flagConfig.DryRun, &flagConfig.Registry, &flagConfig.Image, tags[i], &flagConfig.ExcludeTags)
				if err != nil {
					return err
				}
			}
		} else {
			results, err := getDeleteList(&flagConfig.Registry, &flagConfig.Image, &flagConfig.ExcludeRecent)
			if err != nil {
				return err
			}
			for j := 0 + flagConfig.ExcludeRecent; j < len(results); j++ {
				err := deleteTag(apiSpec, &flagConfig.DryRun, &flagConfig.Registry, &flagConfig.Image, fmt.Sprintf("%s", results[j].(map[string]interface{})["name"]), &flagConfig.ExcludeTags)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func getDeleteList(registry *string, image *string, recent *int) (Results, error) {
	apiSpec := ApiSpec{}
	apiSpec.method = "GET"
	var results Results
	for i := 1; ; i++ {
		apiSpec.path = "/ncr/api/v2/repositories/" + *registry + "/" +
			url.QueryEscape(*image) + "/tags?page=" + strconv.Itoa(i)
		apiSpec.headers = signature.GetHeader(&apiSpec.method, &apiSpec.path)
		data, _, err := sendRequest(&apiSpec)
		if err != nil {
			return nil, err
		}
		var body Body
		err = json.Unmarshal(*data, &body)
		if err != nil {
			return nil, fmt.Errorf("An error occurred parsing the response body. %s\n", err)
		}
		for _, v := range body.Results {
			results = append(results, v)
		}
		if body.Next == 0 {
			break
		}
	}
	if len(results) <= *recent {
		os.Stderr.WriteString("The " + *registry + "/" + *image + " image contains fewer tags than exclude-recent. - Skipping\n")
		return nil, nil
	}
	sort.Slice(results, func(i, j int) bool {
		now, _ := strconv.Atoi(strings.Split(fmt.Sprintf("%f", results[i].(map[string]interface{})["last_updated"]), ".")[0])
		next, _ := strconv.Atoi(strings.Split(fmt.Sprintf("%f", results[j].(map[string]interface{})["last_updated"]), ".")[0])
		return now > next
	})
	return results, nil
}

func deleteTag(apiSpec *ApiSpec, dryRun *bool, registry *string, image *string, tag string, excludeTags *[]string) error {
	for i := range *excludeTags {
		if tag == (*excludeTags)[i] {
			fmt.Println("The " + tag + " tag is included in the exclude tag list. - Skipping")
			return nil
		}
	}
	if *dryRun {
		fmt.Println(*registry + "/" + *image + ":" + tag + " deleted. - Dry run")
		return nil
	}
	apiSpec.path = "/ncr/api/v2/repositories/" + *registry + "/" +
		url.QueryEscape(*image) + "/tags/" + tag
	apiSpec.headers = signature.GetHeader(&apiSpec.method, &apiSpec.path)
	_, status, err := sendRequest(apiSpec)
	if err != nil {
		os.Stderr.WriteString(err.Error())
		return err
	}
	if status == 204 {
		fmt.Println(*registry + "/" + *image + ":" + tag + " successfully deleted.")
		return nil
	} else {
		return errors.New(*registry + "/" + *image + ":" + tag + " not exists.\n")
	}
}
