package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time"

	"gopkg.in/yaml.v3"
)

var (
	supervisorAddress = os.Getenv("BALENA_SUPERVISOR_ADDRESS")
	supervisorKey     = os.Getenv("BALENA_SUPERVISOR_API_KEY")
)

type yamlFields struct {
	Name       string
	Url        string
	Path       string
	Supervisor bool
	Cmd        string
	Cmd_type   string
	Payload    map[string]interface{}
}

func getYaml(filename string) (map[string]yamlFields, error, error) {
	file, err := ioutil.ReadFile(filename)
	data := make(map[string]yamlFields)
	err1 := yaml.Unmarshal(file, &data)

	return data, err, err1
}

func fetchData() error {
	var wg sync.WaitGroup

	yamlData, err, err1 := getYaml("cmds.yaml")
	if err != nil {
		log.Fatal("cmds.yaml does not exist", err)
	}
	if err1 != nil {
		log.Fatal("cmds.yaml is invalid", err1)
	}

	println("\033[34m", " == Extracted Data == ")

	for _, value := range yamlData {
		switch value.Cmd_type {
		case "api":
			if value.Supervisor {
				if supervisorAddress == "" {
					log.Println("The Balena Supervisor environment variables are missing from the device container. ",
						"Check the `io.balena.features.supervisor-api: 1` label is present in your Docker Compose ",
						"file or run command.")
					break
				}
				wg.Add(1)
				go getApiJSON(&wg, supervisorAddress+value.Url, supervisorKey, value.Name, value.Payload)
			} else {
				wg.Add(1)
				go getApiJSON(&wg, value.Url, "", value.Name, value.Payload)
			}

		case "file":
			wg.Add(1)
			go getFileContent(&wg, value.Path, value.Name)
		case "shell":
			wg.Add(1)
			go getShellOutput(&wg, value.Cmd, value.Name)
		}
	}

	wg.Wait()

	return nil
}

func getApiJSON(wg *sync.WaitGroup, url string, apiKey string, source string, jsonPayload map[string]interface{}) error {
	var err error
	var formattedJSON []byte
	var req *http.Request

	defer wg.Done()
	clientRequest := http.Client{
		Timeout: time.Second * 3,
	}

	if jsonPayload == nil {
		req, err = http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			log.Println(err)
			return nil
		}
	} else {
		mJson, err := json.Marshal(jsonPayload)
		if err != nil {
			log.Println("Failed processing JSON payload", err)
			return nil
		}

		req, err = http.NewRequest(http.MethodPost, url, bytes.NewReader(mJson))
		if err != nil {
			log.Println("Failed sending http request", err)
			return nil
		}
	}

	if apiKey != "" {
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", "Bearer "+apiKey)
	}

	res, err := clientRequest.Do(req)
	if err != nil {
		log.Println(err)
		return nil
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return nil
	}

	if json.Valid(body) {
		var prettyJSON bytes.Buffer
		json.Indent(&prettyJSON, body, "", "\t")
		formattedJSON = prettyJSON.Bytes()
		uploadContent(&formattedJSON, source)
	} else {
		uploadContent(&body, source)
	}

	return nil
}

func getFileContent(wg *sync.WaitGroup, file string, source string) error {
	defer wg.Done()
	fileContent, err := ioutil.ReadFile(file)
	if err != nil {
		log.Println(err)
		return nil
	}

	uploadContent(&fileContent, source)
	return nil
}

func getShellOutput(wg *sync.WaitGroup, shellCommand string, source string) error {
	defer wg.Done()
	cmd, err := exec.Command("sh", "-c", shellCommand).Output()

	if err != nil {
		log.Println(err)
		return nil
	}

	uploadContent(&cmd, source)
	return nil
}
