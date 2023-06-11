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
	if len(token) == 0 {
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

	ipv6Prefix, ok := event["ipv6_prefix"].(string)
	if !ok || len(ipv6Prefix) == 0 {
		return map[string]any{
			"body":       "Param ipv6_prefix must be provided with the request",
			"statusCode": 400,
		}
	}

	fmt.Printf("ipv6_prefix: %s\n", ipv6Prefix)

	ipv6PrefixParts := strings.Split(ipv6Prefix, "/")
	ipv6Prefix = ipv6PrefixParts[0]
	ipv6Prefix = strings.TrimSuffix(ipv6Prefix, "::")
	ipv6Prefix = strings.TrimSuffix(ipv6Prefix, ":")

	ipv6, ok := event["ipv6"].(string)
	if !ok || len(ipv6) == 0 {
		return map[string]any{
			"body":       "Param ipv6 must be provided with the request",
			"statusCode": 400,
		}
	}

	interfaceId, ok := event["interface_id"]
	if ok {
		interfaceId, ok := interfaceId.(string)
		if ok && len(interfaceId) > 0 {
			interfaceId = strings.TrimPrefix(interfaceId, "::")
			interfaceId = strings.TrimPrefix(interfaceId, ":")

			ipv6 = ipv6Prefix + ":" + interfaceId

			ipv6Parts := strings.Split(ipv6, ":")
			if len(ipv6Parts) < 8 {
				ipv6 = ipv6Prefix + "::" + interfaceId
			}
		}
	}

	fmt.Printf("ipv6: %s\n", ipv6)

	matchingRecords, res, err := client.Domains.RecordsByName(ctx, domain, record, nil)
	if err != nil {
		return map[string]any{
			"body":       err.Error(),
			"statusCode": res.StatusCode,
		}
	}

	ipv6Records := lo.Filter(matchingRecords, func(record godo.DomainRecord, _ int) bool {
		return record.Type == "AAAA"
	})
	foundIpv6Record := len(ipv6Records) > 0

	if foundIpv6Record {
		ipv6Record := ipv6Records[0]

		_, res, err := client.Domains.EditRecord(ctx, domain, ipv6Record.ID, &godo.DomainRecordEditRequest{
			Type: ipv6Record.Type,
			Name: ipv6Record.Name,
			Data: ipv6,
			TTL:  ipv6Record.TTL,
		})
		if err != nil {
			return map[string]any{
				"body":       err.Error(),
				"statusCode": res.StatusCode,
			}
		}

		fmt.Printf("updated the AAAA record '%s'\n", record)
	} else {
		recordName := strings.Replace(record, "."+domain, "", 1)

		_, res, err := client.Domains.CreateRecord(ctx, domain, &godo.DomainRecordEditRequest{
			Type: "AAAA",
			Name: recordName,
			Data: ipv6,
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
