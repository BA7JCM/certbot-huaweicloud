package api

import (
	"fmt"

	dnsM "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/dns/v2"
	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/dns/v2/model"
	"github.com/spf13/viper"
)

func ListPublicZones(client *dnsM.DnsClient) (*[]model.PublicZoneResp, error) {

	request := &model.ListPublicZonesRequest{}
	typeRequest := "public"
	request.Type = &typeRequest
	// limitRequest:= int32(<limit>)
	// request.Limit = &limitRequest
	// markerRequest:= "<marker>"
	// request.Marker = &markerRequest
	// offsetRequest:= int32(<offset>)
	// request.Offset = &offsetRequest
	// tagsRequest:= "<tags>"
	// request.Tags = &tagsRequest
	// nameRequest:= "<name>"
	// request.Name = &nameRequest
	// statusRequest:= "<status>"
	// request.Status = &statusRequest
	// searchModeRequest:= "<search_mode>"
	// request.SearchMode = &searchModeRequest
	// enterpriseProjectIdRequest:= "<enterprise_project_id>"
	// request.EnterpriseProjectId = &enterpriseProjectIdRequest
	response, err := client.ListPublicZones(request)
	if mode := viper.GetString("MODE"); mode == "DEBUG" {
		fmt.Printf("%+v\n", response)
	}
	return response.Zones, err
}
