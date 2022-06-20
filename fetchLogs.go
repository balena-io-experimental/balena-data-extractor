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
	Name    string
	Url     string
	Key     string
	Path    string
	Cmd     string
	Payload []byte
}

func getYaml(filename string) (map[string]yamlFields, error) {
	file, err := ioutil.ReadFile(filename)
	data := make(map[string]yamlFields)
	err1 := yaml.Unmarshal(file, &data)
	if err1 != nil {
		log.Fatal(err)
	}

	return data, err
}

func processYaml(wg sync.WaitGroup) {
	yamlData, err := getYaml("cmds.yaml")
	if err != nil {
		log.Println("cmds.yaml does not exist. Continuing...")
		return
	}

	for key, value := range yamlData {
		switch key {
		case "api":
			wg.Add(1)
			go getApiJSON(&wg, value.Url, value.Key, value.Name, value.Payload)
		case "file":
			wg.Add(1)
			go getFileContent(&wg, value.Path, value.Name)
		case "shell":
			wg.Add(1)
			go getShellOutput(&wg, value.Cmd, value.Name)
		}
	}
}

func fetchLogs() error {
	var wg sync.WaitGroup

	journalCtlValues := map[string]string{"follow": "false", "count": "10000"}
	journalCtlJson, err := json.Marshal(journalCtlValues)

	if err != nil {
		log.Fatal(err)
	}

	println("\033[34m", " == Logs == ")

	switch chosenLog {
	case "all":
		// Ensure the below wg int matches the number of executed processes in this section
		wg.Add(4)

		go getApiJSON(&wg, supervisorAddress+"/v2/journal-logs", supervisorKey, "JournalCtl", journalCtlJson)
		go getApiJSON(&wg, supervisorAddress+"/v1/device", supervisorKey, "Device Info", nil)
		go getShellOutput(&wg, "ifconfig", "Network Interfaces")
		go getShellOutput(&wg, "printenv | grep -v API_KEY", "Environment Variables")
		//go getFileContent(&wg, "/file/path", "Fetch from a file")
		processYaml(wg)
	case "deviceinfo":
		wg.Add(1)
		go getApiJSON(&wg, supervisorAddress+"/v2/local/device-info", supervisorKey, "Device Info", nil)
	case "envvars":
		wg.Add(1)
		go getShellOutput(&wg, "printenv | grep -v API_KEY", "Environment Variables")
	case "journalctl":
		wg.Add(1)
		go getApiJSON(&wg, supervisorAddress+"/v2/journal-logs", supervisorKey, "JournalCtl", journalCtlJson)
	case "networkinterfaces":
		wg.Add(1)
		go getShellOutput(&wg, "ifconfig", "Network Interfaces")
	case "yaml":
		processYaml(wg)
	}

	wg.Wait()

	return nil
}

func getApiJSON(wg *sync.WaitGroup, url string, apiKey string, source string, jsonPayload []byte) error {
	var err error
	var formattedJSON []byte
	var req *http.Request

	defer wg.Done()
	clientRequest := http.Client{
		Timeout: time.Second * 3, // Timeout after x seconds
	}

	if jsonPayload == nil {
		req, err = http.NewRequest(http.MethodGet, url, nil)
	} else {
		req, err = http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonPayload))
	}

	if err != nil {
		log.Fatal(err)
	}

	if apiKey != "" {
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Authorization", "Bearer "+apiKey)
	}

	res, getErr := clientRequest.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
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
		log.Fatal(err)
	}

	uploadContent(&fileContent, source)
	return nil
}

func getShellOutput(wg *sync.WaitGroup, shellCommand string, source string) error {
	defer wg.Done()
	cmd, err := exec.Command("sh", "-c", shellCommand).Output()

	if err != nil {
		log.Fatal(err)
	}

	uploadContent(&cmd, source)
	return nil
}
