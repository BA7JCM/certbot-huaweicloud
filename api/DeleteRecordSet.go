package api

import (
	"fmt"

	dnsM "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/dns/v2"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/dns/v2/model"
	"github.com/spf13/viper"
)

func DeleteRecordSet(client *dnsM.DnsClient, zoneId string, recordSetId string) error {
	request := &model.DeleteRecordSetRequest{}
	request.ZoneId = zoneId
	request.RecordsetId = recordSetId
	response, err := client.DeleteRecordSet(request)
	if mode := viper.GetString("MODE"); mode == "DEBUG" {
		fmt.Printf("%+v\n", response)
	}
	return err
}
