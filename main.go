package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
)

var (
	burnAfterRead bool
	chosenLog     string
	defaultURL    = os.Getenv("PRIVATEBIN_URL")
	passedURL     string
	password      string
	setExpiry     string
)

func init() {
	if defaultURL == "" {
		defaultURL = "https://privatebin.net/"
	}

	flag.BoolVar(&burnAfterRead, "burn", false, "Burn all data after being read once")
	flag.StringVar(&chosenLog, "data", "all", "Choose which data to export. Options are: 'all', 'deviceinfo', 'envvars', 'journalctl', 'networkinterfaces', 'yaml')")
	flag.StringVar(&password, "password", "", "Set a password for accessing the uploaded content")
	flag.StringVar(&passedURL, "url", defaultURL, "Override the default data host with the passed URL")
	flag.StringVar(&setExpiry, "expire", "day", "Delete all data after specified time. Options are: 'hour', 'day', 'week' or 'month'")

	flag.Parse()

	getURL, err := url.Parse(passedURL)
	if err != nil {
		panic(err)
	}

	host.api = getURL.String()
}

func main() {
	err := fetchLogs()
	if err != nil {
		panic(err)
	}
}

func uploadContent(stdInput *[]byte, message string) error {
	p, err := CraftPaste(*stdInput)
	if err != nil {
		panic(err)
	}
	p.BurnAfterRead(burnAfterRead)
	p.SetExpiry(setExpiry)

	if password != "" {
		p.SetPassword(password)
	}
	ur, _, err := p.Send()
	if err != nil {
		panic(err)
	}

	fmt.Println("\033[32m", message+":", "\033[0m", ur)

	return nil
}
