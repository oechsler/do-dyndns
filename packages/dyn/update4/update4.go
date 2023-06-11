package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/samber/lo"
)

func Main(ctx context.Context, event map[string]any) map[string]any {
	appSecret := os.Getenv("APP_SECRET")
	if len(appSecret) == 0 {
		return map[string]any{
			"body":       "Environment variable APP_SECRET must be set for the function",
			"statusCode": 500,
		}
	}

	token, ok := event["token"].(string)
	if !ok || token != appSecret {
		return map[string]any{
			"body":       "Unauthorized",
			"statusCode": 401,
		}
	}

	clientToken := os.Getenv("DO_TOKEN")
	if len(clientToken) == 0 {
		return map[string]any{
			"body":       "Environment variable DO_TOKEN must be set for the function",
			"statusCode": 500,
		}
	}
	client := godo.NewFromToken(clientToken)

	record, ok := event["record"].(string)
	if !ok || len(record) == 0 {
		return map[string]any{
			"body":       "Param record must be provided with the request",
			"statusCode": 400,
		}
	}
	domainParts := strings.Split(record, ".")[1:]
	domain := strings.Join(domainParts, ".")

	fmt.Printf("domain: '%s'\n", domain)
	fmt.Printf("record: '%s'\n", record)

	ipv4, ok := event["ipv4"].(string)
	if !ok || len(ipv4) == 0 {
		return map[string]any{
			"body":       "Param ipv4 must be provided with the request",
			"statusCode": 400,
		}
	}

	fmt.Printf("ipv4: %s\n", ipv4)

	matchingRecords, res, err := client.Domains.RecordsByName(ctx, domain, record, nil)
	if err != nil {
		return map[string]any{
			"body":       err.Error(),
			"statusCode": res.StatusCode,
		}
	}

	ipv4Records := lo.Filter(matchingRecords, func(record godo.DomainRecord, _ int) bool {
		return record.Type == "A"
	})
	foundIpv4Record := len(ipv4Records) > 0

	if foundIpv4Record {
		ipv4Record := ipv4Records[0]

		_, res, err := client.Domains.EditRecord(ctx, domain, ipv4Record.ID, &godo.DomainRecordEditRequest{
			Type: ipv4Record.Type,
			Name: ipv4Record.Name,
			Data: ipv4,
			TTL:  ipv4Record.TTL,
		})
		if err != nil {
			return map[string]any{
				"body":       err.Error(),
				"statusCode": res.StatusCode,
			}
		}

		fmt.Printf("updated the A record for '%s'\n", record)
	} else {
		recordName := strings.Replace(record, "."+domain, "", 1)

		_, res, err := client.Domains.CreateRecord(ctx, domain, &godo.DomainRecordEditRequest{
			Type: "A",
			Name: recordName,
			Data: ipv4,
			TTL:  60,
		})
		if err != nil {
			return map[string]any{
				"body":       err.Error(),
				"statusCode": res.StatusCode,
			}
		}

		fmt.Printf("created the A record for '%s'\n", record)
	}

	return map[string]any{
		"statusCode": 200,
	}
}
