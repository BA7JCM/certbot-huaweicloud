package api

import (
	"fmt"
	"strings"

	dnsM "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/dns/v2"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/dns/v2/model"
	"github.com/spf13/viper"
)

func CreateRecordSet(client *dnsM.DnsClient, zoneId string, certbotValidation string, certbotDomain string, description string) (string, error) {
	request := &model.CreateRecordSetWithLineRequest{}
	// valueTags := "value1"
	// var listTagsbody = []model.Tag{
	// 	{
	// 		Key:   "key1",
	// 		Value: &valueTags,
	// 	},
	// }
	validationValues := strings.Split(certbotValidation, ",")
	var listRecordsbody []string
	for _, value := range validationValues {
		listRecordsbody = append(listRecordsbody, "\""+value+"\"")
	}
	ttlCreateRecordSetRequestBody := int32(300)
	descriptionCreateRecordSetRequestBody := description
	request.ZoneId = zoneId
	weight := int32(1)
	request.Body = &model.CreateRecordSetWithLineRequestBody{
		// Tags:        &listTagsbody,
		Records:     &listRecordsbody,
		Ttl:         &ttlCreateRecordSetRequestBody,
		Type:        "TXT",
		Description: &descriptionCreateRecordSetRequestBody,
		Name:        "_acme-challenge." + certbotDomain,
		Weight:      &weight,
	}

	// 打印 name 字段内容
	fmt.Println("CreateRecordSetWithLineRequestBody name:", request.Body.Name)

	response, err := client.CreateRecordSetWithLine(request)
	if err != nil {
		return "", err
	}
	if mode := viper.GetString("MODE"); mode == "DEBUG" {
		fmt.Printf("%+v\n", response)
	}
	return *response.Id, err
}
