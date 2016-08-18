package main

import (
	"encoding/json"
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

	// AWS
	zoneID := os.Getenv("ZONEID")

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

	c, err := newClient(zoneID)
	if err != nil {
		io.WriteString(os.Stderr, "Error creating google DNS client"+err.Error())
		os.Exit(1)
	}

	switch command {
	case "CREATE":
		err = c.create(domain, fqdn, token)
	case "DELETE":
		err = c.delete(domain, fqdn, token)
	}

	if err != nil {
		io.WriteString(os.Stderr, err.Error())
		os.Exit(1)
	}
	os.Exit(0)
}
