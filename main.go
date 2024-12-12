package main

import (
	"certbot/api"
	"certbot/config"
	"fmt"
	"time"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/core/auth/basic"
	dnsM "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/dns/v2"
	regionM "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/dns/v2/region"
	"github.com/spf13/viper"
)

func main() {
	config.Init()

	// The AK and SK used for authentication are hard-coded or stored in plaintext, which has great security risks. It is recommended that the AK and SK be stored in ciphertext in configuration files or environment variables and decrypted during use to ensure security.
	// In this example, AK and SK are stored in environment variables for authentication. Before running this example, set environment variables CLOUD_SDK_AK and CLOUD_SDK_SK in the local environment
	ak := viper.GetString("CLOUD_SDK_AK")
	sk := viper.GetString("CLOUD_SDK_SK")
	region := viper.GetString("REGION")
	region_value, err := regionM.SafeValueOf(region)
	if err != nil {
		fmt.Println("region error")
		panic(err)
	}

	auth, err := basic.NewCredentialsBuilder().
		WithAk(ak).
		WithSk(sk).
		SafeBuild()
	if err != nil {
		fmt.Println("auth error")
		panic(err)
	}

	dns, err := dnsM.DnsClientBuilder().
		WithRegion(region_value).
		WithCredential(auth).
		SafeBuild()
	if err != nil {
		fmt.Println("dns error")
		panic(err)
	}

	client := dnsM.NewDnsClient(dns)

	// 获取zoneId
	zoneId := getZoneId(client, viper.GetString("DOMAIN"))
	if zoneId == "" {
		fmt.Println("zone id not found")
		panic("zone id not found")
	}

	mode := viper.GetString("MODE")
	if mode == "DEBUG" {
		// 创建recordSet
		certbotValidation := viper.GetString("CERTBOT_VALIDATION")
		certbotDomain := viper.GetString("CERTBOT_DOMAIN")
		description := viper.GetString("DESCRIPTION")
		recordSetId, err := api.CreateRecordSet(client, zoneId, certbotValidation, certbotDomain, description)
		if err != nil {
			fmt.Println("create record set error")
			panic(err)
		}

		fmt.Println("recordSetId", recordSetId)
		time.Sleep(10 * time.Second)

		// 删除recordSet
		err = api.DeleteRecordSet(client, zoneId, recordSetId)
		if err != nil {
			fmt.Println("delete record set error")
			panic(err)
		}
	} else if mode == "DEV" {
		//TODO
	}
}

func getZoneId(client *dnsM.DnsClient, domain string) string {
	zones, err := api.ListPublicZones(client)
	if err != nil {
		fmt.Println("get zones error")
		panic(err)
	}
	for _, zone := range *zones {
		if *zone.Name == domain {
			return *zone.Id
		}
	}
	return ""
}
