package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

type client struct {
	r53    *route53.Route53
	zoneID string
}

func newClient(zoneID string) (*client, error) {
	return &client{
		r53:    route53.New(session.New()),
		zoneID: zoneID,
	}, nil
}

func (c *client) create(domain, fqdn, token string) error {
	_, err := c.r53.ChangeResourceRecordSets(c.params(route53.ChangeActionCreate, domain, fqdn, token))
	if err != nil {
		return err
	}

	return nil
}

func (c *client) delete(domain, fqdn, token string) error {
	_, err := c.r53.ChangeResourceRecordSets(c.params(route53.ChangeActionDelete, domain, fqdn, token))
	if err != nil {
		return err
	}

	return nil
}

func (c *client) params(action, domain, fqdn, token string) *route53.ChangeResourceRecordSetsInput {
	return &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action: aws.String(action),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: aws.String(fqdn),
						Type: aws.String("TXT"),
						ResourceRecords: []*route53.ResourceRecord{
							// AWS wants the string in quotes
							{Value: aws.String(fmt.Sprintf(`"%s"`, token))},
						},
						TTL: aws.Int64(30),
					},
				},
			},
		},
		HostedZoneId: aws.String(c.zoneID),
	}
}
