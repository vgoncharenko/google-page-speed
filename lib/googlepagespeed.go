package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"encoding/json"
	"flag"
)

const PAGE_SPEED_URL = "https://www.googleapis.com/pagespeedonline/v4/runPagespeed?url=%s/%s&strategy=%s&key=%s"

const OUTPUT_FOLDER = "/root/googlepagespeed/"

var TEST_TYPES = [2]string{"mobile", "desktop"}

type Scenario struct {
	Name string `json:"scenario_name"`
	URL  string `json:"sub_url"`
}

type Response struct {
	Kind         string `json:"kind"`
	ID           string `json:"id"`
	ResponseCode int    `json:"responseCode"`
	Title        string `json:"title"`
	RuleGroups   struct {
			     SPEED struct {
					   Score int `json:"score"`
				   } `json:"SPEED"`
		     } `json:"ruleGroups"`
}

func main() {
	url := flag.String("url", "http://example.com", "a string")
	apiKey := flag.String("apiKey", "00000000000000000000000000000", "a string")
	configPath := flag.String("scenarios", "default", "a string")
	flag.Parse()
	validationParams(*url, *apiKey, *configPath)

	var scenariosConfig []Scenario
	parseConfig(*configPath, &scenariosConfig)

	var wg sync.WaitGroup
	scenariosCount := len(scenariosConfig)
	testTypeCount := len(TEST_TYPES)
	count := scenariosCount * testTypeCount
	result := make(map[string]Response)
	response := Response{}

	wg.Add(count)

	for _, scenario := range scenariosConfig {
		folderName := OUTPUT_FOLDER + scenario.Name
		flushFolder(folderName)
		createFolder(folderName)
		for _, testType := range TEST_TYPES {
			go func(url string, scenario Scenario, testType string, apiKey string, folderName string) {
				defer wg.Done()
				content := request(url, scenario, testType, apiKey)
				name := folderName + "/" + testType + ".json"
				json.Unmarshal(content, &response)
				result[scenario.Name + " " + testType] = response
				storeResults(content, name)
				fmt.Printf("%s %s COMPLETE\n", scenario.Name, testType)
			}(*url, scenario, testType, *apiKey, folderName)
		}
	}
	wg.Wait()

	mapB, _ := json.Marshal(result)
	storeResults([]byte(mapB), OUTPUT_FOLDER + "result.json")
}

func request (url string, scenario Scenario, testType string, apiKey string) ([]byte) {
	requestURL := fmt.Sprintf(PAGE_SPEED_URL, url, scenario.URL, testType, apiKey)
	response, err := http.Get(requestURL)
	check(err)
	content, err := ioutil.ReadAll(response.Body)
	response.Body.Close()
	check(err)

	return content
}

func parseConfig(configPath string, scenariosConfig *[]Scenario)  {
	content, err := ioutil.ReadFile(configPath)
	check(err)
	json.Unmarshal([]byte(content), &scenariosConfig)
}

func validationParams(url string, apiKey string, configPath string) {
	if (url == "http://example.com") {
		log.Fatal("Please set an URL by -url=http://example.com")
	}
	if (apiKey == "00000000000000000000000000000") {
		log.Fatal("Please set an API_KEY by -apiKey=00000000000000000000000000000")
	}
	if (configPath == "default") {
		log.Fatal("Please set a file path to scenarios config by -apiKey=/path/config.json")
	}
}

func flushFolder(folderName string) {
	os.RemoveAll(folderName)
}

func createFolder(folderName string) {
	err := os.Mkdir(folderName, os.ModePerm)
	check(err)
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
		panic(e)
	}
}

func storeResults(robots []byte, name string) {
	err := ioutil.WriteFile(name, robots, 0777)
	check(err)
}
