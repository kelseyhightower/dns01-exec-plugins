// Copyright 2017 Alexander Palaistras. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	// Standard library.
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
)

func main() {
	if os.Getenv("APIVERSION") != "v1" {
		os.Exit(3)
	}

	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		io.WriteString(os.Stderr, err.Error())
		os.Exit(1)
	}

	var client = &Client{}

	err = json.Unmarshal(data, &client)
	if err != nil {
		io.WriteString(os.Stderr, err.Error())
		os.Exit(2)
	}

	var record = &Record{
		Domain: os.Getenv("DOMAIN"),
		Fqdn:   os.Getenv("FQDN"),
		Token:  os.Getenv("TOKEN"),
	}

	switch os.Getenv("COMMAND") {
	case "CREATE":
		err = client.CreateRecord(record)
	case "DELETE":
		err = client.DeleteRecord(record)
	}

	if err != nil {
		io.WriteString(os.Stderr, err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}
