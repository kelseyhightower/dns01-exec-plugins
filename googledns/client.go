// Copyright 2016 Google Inc. All Rights Reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"context"
	"crypto/tls"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/certifi/gocertifi"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/dns/v1"
)

var httpClient http.Client

func init() {
	// Use the Root Certificates bundle from the Certifi project so we don't
	// rely on the host OS or container base images for a CA Bundle.
	// See https://certifi.io for more details.
	certPool, err := gocertifi.CACerts()
	if err != nil {
		log.Fatal(err)
	}
	httpClient = http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{RootCAs: certPool},
		},
	}
}

type client struct {
	domain  string
	project string
	*dns.Service
}

func newClient(serviceAccount []byte, project, domain string) (*client, error) {
	jwtConfig, err := google.JWTConfigFromJSON(
		serviceAccount,
		dns.NdevClouddnsReadwriteScope,
	)
	if err != nil {
		return nil, err
	}

	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, &httpClient)

	jwtHTTPClient := jwtConfig.Client(ctx)
	service, err := dns.New(jwtHTTPClient)
	if err != nil {
		return nil, err
	}

	return &client{domain, project, service}, nil
}

func (c *client) create(fqdn, value string, ttl int) error {
	zones, err := c.ManagedZones.List(c.project).Do()
	if err != nil {
		return err
	}

	zoneName := ""
	for _, zone := range zones.ManagedZones {
		if strings.HasSuffix(c.domain+".", zone.DnsName) {
			zoneName = zone.Name
		}
	}
	if zoneName == "" {
		return errors.New("Zone not found")
	}

	record := &dns.ResourceRecordSet{
		Name:    fqdn,
		Rrdatas: []string{value},
		Ttl:     int64(ttl),
		Type:    "TXT",
	}

	change := &dns.Change{
		Additions: []*dns.ResourceRecordSet{record},
	}

	changesCreateCall, err := c.Changes.Create(c.project, zoneName, change).Do()
	if err != nil {
		return err
	}

	for changesCreateCall.Status == "pending" {
		time.Sleep(time.Second)
		changesCreateCall, err = c.Changes.Get(c.project, zoneName, changesCreateCall.Id).Do()
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *client) delete(fqdn string) error {
	zones, err := c.ManagedZones.List(c.project).Do()
	if err != nil {
		return err
	}

	zoneName := ""
	for _, zone := range zones.ManagedZones {
		if strings.HasSuffix(c.domain+".", zone.DnsName) {
			zoneName = zone.Name
		}
	}
	if zoneName == "" {
		return errors.New("Zone not found")
	}

	records, err := c.ResourceRecordSets.List(c.project, zoneName).Do()
	if err != nil {
		return err
	}

	matchingRecords := []*dns.ResourceRecordSet{}
	for _, record := range records.Rrsets {
		if record.Type == "TXT" && record.Name == fqdn {
			matchingRecords = append(matchingRecords, record)
		}
	}

	for _, record := range matchingRecords {
		change := &dns.Change{
			Deletions: []*dns.ResourceRecordSet{record},
		}
		_, err = c.Changes.Create(c.project, zoneName, change).Do()
		if err != nil {
			return err
		}
	}

	return nil
}
