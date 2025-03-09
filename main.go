package main

import (
	"bufio"
	"certbot/api"
	"certbot/config"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
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

	var use string
	flag.StringVar(&use, "use", "", "use")
	flag.Parse()

	switch use {
	case "auth": // 创建recordSet
		certbotValidation, certbotDomain, description := loadParams()
		saveValidation(certbotDomain, certbotValidation)

		// 检查是否所有验证值都已收集完毕
		remainingChallengesStr := os.Getenv("CERTBOT_REMAINING_CHALLENGES")
		if remainingChallengesStr == "" {
			fmt.Println("CERTBOT_REMAINING_CHALLENGES is not set")
			panic("CERTBOT_REMAINING_CHALLENGES is not set")
		}

		remainingChallenges, err := strconv.Atoi(remainingChallengesStr)
		if err != nil {
			fmt.Println("invalid CERTBOT_REMAINING_CHALLENGES value")
			panic(err)
		}

		if remainingChallenges == 0 {
			validations := loadValidations()
			for domain, validation := range validations {
				// fmt.Println("creating record for domain:", domain, "value:", strings.Join(validation, ","))
				recordSetId, err := api.CreateRecordSet(client, zoneId, strings.Join(validation, ","), domain, description)
				if err != nil {
					fmt.Println("create record set error")
					panic(err)
				}
				fmt.Println("creating record")
				saveRecordSetId(recordSetId) // 保存实际的 recordSetId
			}
			// 等一分钟让dns解析生效
			fmt.Println("wating 1 min for dns")
			time.Sleep(60 * time.Second)
			os.Remove("validations.txt") // 删除validations.txt文件
		}

	case "cleanup": // 删除recordSet
		remainingChallengesStr := os.Getenv("CERTBOT_REMAINING_CHALLENGES")
		if remainingChallengesStr == "" {
			fmt.Println("CERTBOT_REMAINING_CHALLENGES is not set, skipping cleanup")
			return
		}

		remainingChallenges, err := strconv.Atoi(remainingChallengesStr)
		if err != nil {
			fmt.Println("invalid CERTBOT_REMAINING_CHALLENGES value, skipping cleanup")
			return
		}

		if remainingChallenges != 0 {
			fmt.Println("remaining challenges:", remainingChallenges, "skipping cleanup")
			return
		}

		if _, err := os.Stat("recordSetId.txt"); os.IsNotExist(err) {
			fmt.Println("recordSetId.txt does not exist, skipping cleanup")
			return
		}
		recordSetIds := loadRecordSetId()
		for _, recordSetId := range recordSetIds {
			err := api.DeleteRecordSet(client, zoneId, recordSetId)
			if err != nil {
				fmt.Println("delete record set error")
				panic(err)
			}
		}
		// 删除recordSetId.txt文件
		fmt.Println("delete recordSetId.txt")
		err = os.Remove("recordSetId.txt")
		if err != nil {
			fmt.Println("delete recordSetId.txt error")
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
	file, err := os.OpenFile("recordSetId.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("save record set id error")
		panic(err)
	}
	_, err = file.WriteString(recordSetId + "\n")
	if err != nil {
		fmt.Println("write record set id error")
		panic(err)
	}
	file.Close()
}

func loadRecordSetId() []string {
	file, err := os.Open("recordSetId.txt")
	if err != nil {
		fmt.Println("load record set id error")
		panic(err)
	}
	defer file.Close()

	var recordSetIds []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		recordSetIds = append(recordSetIds, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("read record set id error")
		panic(err)
	}

	return recordSetIds
}

func saveValidation(certbotDomain, certbotValidation string) {
	file, err := os.OpenFile("validations.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("save validation error")
		panic(err)
	}
	_, err = file.WriteString(certbotDomain + ":" + certbotValidation + "\n")
	if err != nil {
		fmt.Println("write validation error")
		panic(err)
	}
	file.Close()
}

func loadValidations() map[string][]string {
	file, err := os.Open("validations.txt")
	if err != nil {
		fmt.Println("load validations error")
		panic(err)
	}
	defer file.Close()

	validations := make(map[string][]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			domain := parts[0]
			validation := parts[1]
			validations[domain] = append(validations[domain], validation)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("read validations error")
		panic(err)
	}

	return validations
}
