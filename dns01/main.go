package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type Config struct {
	Username string `json:"username"`
	Token    string `json:"token"`
}

func main() {
	apiVersion := os.Getenv("APIVERSION")
	command := os.Getenv("COMMAND")
	domain := os.Getenv("DOMAIN")
	fqdn := os.Getenv("FQDN")
	token := os.Getenv("TOKEN")

	if apiVersion != "v1" {
		os.Exit(3)
	}

	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		io.WriteString(os.Stderr, err.Error())
		os.Exit(1)
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		io.WriteString(os.Stderr, err.Error())
		os.Exit(2)
	}

	switch command {
	case "CREATE":
		err = createRecord(domain, fqdn, token, config)
	case "DELETE":
		err = deleteRecord(domain, fqdn, config)
	}
	if err != nil {
		io.WriteString(os.Stderr, err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}

func createRecord(domain, fqdn, token string, config Config) error {
	fmt.Printf("%s TXT %s\n", fqdn, token)
	return nil
}

func deleteRecord(domain, fqdn string, config Config) error {
	return nil
}
