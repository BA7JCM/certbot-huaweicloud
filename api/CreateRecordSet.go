package api

import (
	"fmt"

	dnsM "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/dns/v2"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/dns/v2/model"
	"github.com/spf13/viper"
)

func CreateRecordSet(client *dnsM.DnsClient, zoneId string, certbotValidation string, certbotDomain string, description string) (string, error) {
	request := &model.CreateRecordSetRequest{}
	// valueTags := "value1"
	// var listTagsbody = []model.Tag{
	// 	{
	// 		Key:   "key1",
	// 		Value: &valueTags,
	// 	},
	// }
	var listRecordsbody = []string{
		"\"" + certbotValidation + "\"",
	}
	ttlCreateRecordSetRequestBody := int32(300)
	descriptionCreateRecordSetRequestBody := description
	request.ZoneId = zoneId
	request.Body = &model.CreateRecordSetRequestBody{
		// Tags:        &listTagsbody,
		Records:     listRecordsbody,
		Ttl:         &ttlCreateRecordSetRequestBody,
		Type:        "TXT",
		Description: &descriptionCreateRecordSetRequestBody,
		Name:        certbotDomain,
	}
	response, err := client.CreateRecordSet(request)
	if mode := viper.GetString("MODE"); mode == "DEBUG" {
		fmt.Printf("%+v\n", response)
	}
	return *response.Id, err
}
