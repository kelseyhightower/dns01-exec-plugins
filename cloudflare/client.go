// Copyright 2017 Alexander Palaistras. All rights reserved.
// Use of this source code is governed by the MIT license that can be found in
// the LICENSE file.

package main

import (
	// Standard library.
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// The URI against all requests will be made.
const baseURI = "https://api.cloudflare.com/client/v4"

// A Zone represents a host underneath which records are managed.
type Zone struct {
	ID string `json:"id"`
}

// A Record represents a DNS record, as managed by Cloudflare DNS.
type Record struct {
	ID     string `json:"id,omitempty"`
	Type   string `json:"type"`
	Domain string `json:"zone_name"`
	Fqdn   string `json:"name"`
	Token  string `json:"content"`
}

// A Client represents a Cloudflare API client, and can be used for making
// generic, authenticated API requests.
type Client struct {
	AuthEmail string `json:"email"`
	AuthKey   string `json:"key"`
}

// A Response represents a complete API response from the Cloudflare API.
type Response struct {
	Result  interface{} `json:"result"`
	Success bool        `json:"success"`
	Errors  Errors      `json:"errors"`
}

// Errors represent a list of numeric error codes and descriptive error messages,
// which can help with determining request failures.
type Errors []struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Chain   Errors `json:"error_chain"`
}

// Error returns a comma-separated list of error messages and their numeric codes.
func (e Errors) Error() string {
	var msg string
	for _, err := range e {
		msg += fmt.Sprintf(", %s (Code %d)", err.Message, err.Code)
		if len(err.Chain) > 0 {
			msg += ": " + err.Chain.Error()
		}
	}

	if len(msg) > 2 {
		return msg[2:]
	}

	return ""
}

// CreateRecord adds or updates a DNS record for the FQDN and token provided, as
// a TXT record for the zone pointed to by the record domain name. An error is
// returned if the process fails at any point.
func (c *Client) CreateRecord(record *Record) error {
	// Get parent zone for domain.
	zone, err := c.fetchZone(record.Domain)
	if err != nil {
		return fmt.Errorf("Fetching zone for record failed: %s", err)
	} else if zone == nil {
		return fmt.Errorf("No zone found for host '%s'", record.Domain)
	}

	// Get existing record, if any.
	existing, err := c.fetchRecord(zone, "TXT", record.Fqdn)
	if err != nil {
		return fmt.Errorf("Fetching record failed: %s", err)
	}

	// Update existing record, if one exists, otherwise create new record.
	if existing != nil {
		existing.Token = record.Token
		err := c.sendRequest("PUT", ("/zones/" + zone.ID + "/dns_records/" + existing.ID), existing)
		if err != nil {
			return fmt.Errorf("Updating record failed: %s", err)
		}
	} else {
		record.Type = "TXT"
		err := c.sendRequest("POST", ("/zones/" + zone.ID + "/dns_records/"), record)
		if err != nil {
			return fmt.Errorf("Creating record failed: %s", err)
		}
	}

	return nil
}

// DeleteRecord removes a DNS record of type TXT for the FQDN and zone domain name
// pointed to.
func (c *Client) DeleteRecord(record *Record) error {
	// Get parent zone for domain.
	zone, err := c.fetchZone(record.Domain)
	if err != nil {
		return fmt.Errorf("Fetching zone for record failed: %s", err)
	} else if zone == nil {
		return fmt.Errorf("No zone found for host '%s'", record.Domain)
	}

	// Get and delete existing record, if any.
	existing, err := c.fetchRecord(zone, "TXT", record.Fqdn)
	if err != nil {
		return fmt.Errorf("Fetching record failed: %s", err)
	} else if existing == nil {
		return nil
	}

	err = c.sendRequest("DELETE", ("/zones/" + zone.ID + "/dns_records/" + existing.ID), nil)
	if err != nil {
		return fmt.Errorf("Deleting record failed: %s", err)
	}

	return nil
}

// FetchZone makes an API request for zone information stored against the domain
// name provided. A nil Zone and no error is returned if no zone exists for the
// domain name.
func (c *Client) fetchZone(domain string) (*Zone, error) {
	var zones []*Zone

	if err := c.sendRequest("GET", ("/zones/?name=" + domain), &zones); err != nil {
		return nil, fmt.Errorf("API request failed: %s", err)
	} else if len(zones) == 0 {
		return nil, nil
	}

	return zones[0], nil
}

// FetchRecord makes an API request for DNS record information stored against the
// record type and name provided, in the parent zone specified. If more than one
// records are returned, an error is returned.
func (c *Client) fetchRecord(z *Zone, kind, name string) (*Record, error) {
	var records []*Record
	var uri = fmt.Sprintf("/zones/%s/dns_records?type=%s&name=%s", z.ID, kind, name)

	if err := c.sendRequest("GET", uri, &records); err != nil {
		return nil, fmt.Errorf("API request failed: %s", err)
	} else if len(records) > 1 {
		return nil, fmt.Errorf("More than one records returned for type %s and name %s", kind, name)
	} else if len(records) == 0 {
		return nil, nil
	}

	return records[0], nil
}

// SendRequest performs an API request for the method (GET, PUT, etc) and URI
// provided. The object provided is placed in the request body where applicable
// (i.e. in PUT and POST requests) and is used for determining the correct result
// type, if any is returned.
func (c *Client) sendRequest(method, uri string, object interface{}) error {
	var req *http.Request
	var err error

	// Prepare request body, if required.
	switch method {
	case "GET", "DELETE":
		req, err = http.NewRequest(method, (baseURI + uri), nil)
		if err != nil {
			return err
		}
	case "POST", "PUT":
		buf, err := json.Marshal(object)
		if err != nil {
			return err
		}

		req, err = http.NewRequest(method, (baseURI + uri), bytes.NewBuffer(buf))
		if err != nil {
			return err
		}
	}

	// Build and make remote request.
	req.Header = map[string][]string{
		"X-Auth-Email": {c.AuthEmail},
		"X-Auth-Key":   {c.AuthKey},
	}

	resp, err := (new(http.Client)).Do(req)
	if err != nil {
		return fmt.Errorf("Failed sending request: %s", err)
	}

	defer resp.Body.Close()

	// Read response for object type requested.
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var result = &Response{Result: object}
	if err = json.Unmarshal(data, result); err != nil {
		return err
	}

	if result.Success == false {
		return result.Errors
	}

	return nil
}
