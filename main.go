package main

import (
	"certbot/api"
	"certbot/config"
	"flag"
	"fmt"
	"io"
	"os"

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

	var use string
	flag.StringVar(&use, "use", "", "use")
	flag.Parse()

	switch use {
	case "auth": // 创建recordSet
		certbotValidation, certbotDomain, description := loadParams()
		recordSetId, err := api.CreateRecordSet(client, zoneId, certbotValidation, certbotDomain, description)
		if err != nil {
			fmt.Println("create record set error")
			panic(err)
		}
		saveRecordSetId(recordSetId)

	case "cleanup": // 删除recordSet
		recordSetId := loadRecordSetId()
		err := api.DeleteRecordSet(client, zoneId, recordSetId)
		if err != nil {
			fmt.Println("delete record set error")
			panic(err)
		}

	default:
		fmt.Println("use param error. use auth or cleanup")
		panic("use param error")
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

func loadParams() (certbotValidation, certbotDomain, description string) {
	mode := viper.GetString("MODE")
	description = viper.GetString("DESCRIPTION")
	if mode == "DEBUG" {
		certbotValidation = viper.GetString("CERTBOT_VALIDATION")
		certbotDomain = viper.GetString("CERTBOT_DOMAIN")
	} else if mode == "DEV" {
		certbotValidation = os.Getenv("CERTBOT_VALIDATION")
		certbotDomain = os.Getenv("CERTBOT_DOMAIN")
	} else {
		fmt.Println("mode error. mode must be DEBUG or DEV")
		panic("mode error")
	}
	return certbotValidation, certbotDomain, description
}

func saveRecordSetId(recordSetId string) {
	file, err := os.Create("recordSetId.txt")
	if err != nil {
		fmt.Println("save record set id error")
		panic(err)
	}
	file.WriteString(recordSetId)
	file.Close()
}

func loadRecordSetId() string {
	file, err := os.Open("recordSetId.txt")
	if err != nil {
		fmt.Println("load record set id error")
		panic(err)
	}
	recordSetId, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("load record set id error")
		panic(err)
	}
	file.Close()
	os.Remove("recordSetId.txt")
	return string(recordSetId)
}
