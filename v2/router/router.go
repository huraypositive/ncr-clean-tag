package router

import (
  "encoding/json"
  // "bytes"
  "net/http"
  "net/http/httputil"
  "fmt"
  "io/ioutil"
  "os"
  "strconv"
  "strings"
  "github.com/joho/godotenv"
  "nct/signature"
)

var apigw string
var debug string
type ApiSpec struct {
  headers     *map[string]string
  method	    string
  path		    string
  jsonContent bool
  requestBody map[string]interface{}
}

func init() {
  err := godotenv.Load()
  if err != nil {
    os.Stderr.WriteString(err.Error())
  }
  if apigw = os.Getenv("NCR_API_GATEWAY_URL"); apigw == "" {
    apigw = "https://ncr.apigw.ntruss.com"
  }
  if debug = os.Getenv("DEBUG"); debug == "" || strings.ToLower(debug) != "true" {
    debug = "false"
  }
}

func MakeApiSpec() *ApiSpec {
  var apiSpec ApiSpec
  return &apiSpec
}

func GetRepos(apiSpec *ApiSpec) {
  apiSpec.method = "GET"
  apiSpec.path = "/ncr/api/v2/repositories"
  apiSpec.headers = signature.GetHeader(&apiSpec.method,&apiSpec.path) 
  apiSpec.jsonContent = true
  data,err := sendRequest(apiSpec)
  if err != nil {
    return
  }
  type Body struct {
    Results []map[string]string `json:"results"`
  }
  var body Body
  err = json.Unmarshal(*data, &body)
  fmt.Println("NAME")
  for _,v := range body.Results {
    fmt.Println(v["name"])
  }
}

func sendRequest(apiSpec *ApiSpec) (*[]byte, error) {
  if (apiSpec.jsonContent) {
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
    req, err = http.NewRequest(apiSpec.method, apigw + apiSpec.path, nil)
  }
  if err != nil {
    return nil,err
  }
  for k,v := range *apiSpec.headers {
    req.Header.Set(k,v)
  }
  client := &http.Client{}
  res, err := client.Do(req)
  if err != nil {
    return nil,err
  }
  defer res.Body.Close()

  debugMode, err := strconv.ParseBool(debug)
  if err != nil {
    return nil,err
  }
  if debugMode {
    data, err := httputil.DumpResponse(res, true)
    if err != nil {
      return nil,err
    }
    fmt.Println(string(data))
  } 
  data, err := ioutil.ReadAll(res.Body)
  if err != nil {
    return nil,err
  }
  return &data,nil
}