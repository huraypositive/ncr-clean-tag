package router

import (
	"fmt"
	"io/ioutil"
	"nct/config"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

var apigw string
var debug string

type Results []interface{}
type Body struct {
	Next    int     `json:"next"`
	Results Results `json:"results"`
}

type ApiSpec struct {
	headers     *map[string]string
	method      string
	path        string
	jsonContent bool
	requestBody map[string]interface{}
}

func init() {
	if _, err := os.Stat(".env"); err == nil {
		err = godotenv.Load()
		if err != nil {
			os.Stderr.WriteString(err.Error())
		}
	}
	if apigw = os.Getenv("NCR_API_GATEWAY_URL"); apigw == "" {
		apigw = "https://ncr.apigw.ntruss.com"
	}
	if debug = os.Getenv("DEBUG"); debug == "" || strings.ToLower(debug) != "true" {
		debug = "false"
	}
}

func makeApiSpec() *ApiSpec {
	var apiSpec ApiSpec
	return &apiSpec
}

func getFlagConfig(cmd ...string) *config.Flag {
	var flagConfig config.Flag
	switch cmd[0] {
	case "get":
		flagConfig = &config.GetFlag{}
		flags := flagConfig.Setup(&cmd)
		flagConfig.Parse(flags, os.Args[2:])
	case "delete":
		flagConfig = &config.DeleteFlag{}
		flags := flagConfig.Setup(&cmd)
		flagConfig.Parse(flags, os.Args[2:])
	}
	return &flagConfig
}

func sendRequest(apiSpec *ApiSpec) (*[]byte, int, error) {
	if apiSpec.jsonContent {
		(*apiSpec.headers)["Content-Type"] = "application/json; charset=utf-8"
	}
	var req *http.Request
	var err error
	if apiSpec.requestBody != nil {
		// repoBytes, err := json.Marshal(repo)
		// if err != nil {
		//   return err
		// }
		// buff = bytes.NewBuffer(repoBytes)
		// req, err := http.NewRequest(apiSpec.method, apigw + apiSpec.path, nil)
	} else {
		req, err = http.NewRequest(apiSpec.method, apigw+apiSpec.path, nil)
	}
	if err != nil {
		return nil, 0, err
	}
	for k, v := range *apiSpec.headers {
		req.Header.Set(k, v)
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer res.Body.Close()

	debugMode, err := strconv.ParseBool(debug)
	if err != nil {
		return nil, 0, err
	}
	if debugMode {
		data, err := httputil.DumpResponse(res, true)
		if err != nil {
			return nil, 0, err
		}
		fmt.Println(string(data))
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, 0, err
	}
	return &data, res.StatusCode, nil
}
