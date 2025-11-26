package alibabacloudstack

import (
	"fmt"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
)

type UniversalDnsService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *UniversalDnsService) DescribeUniversalZones(id string) (map[string]interface{}, error) {

	query := map[string]interface{}{
		"PageNumber": 1,
		"PageSize":   100,
	}

	resp, err := s.client.DoTeaRequest("POST", "UniversalDns", "2021-06-24", "DescribeUniversalZones", "", nil, nil, query)
	if err != nil {
		return nil, err
	}

	if success, ok := resp["success"].(bool); !ok || !success {
		return nil, errmsgs.WrapError(fmt.Errorf("Failed to describe universal zones: %v", resp))
	}

	data, ok := resp["Data"].([]interface{})
	if !ok || data == nil {
		return nil, errmsgs.WrapError(errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("Resource UniversalZones:%s not found", id)))
	}
	for _, item := range data {
		itemMap, isMap := item.(map[string]interface{})
		if !isMap {
			continue
		}
		domainid, idOk := itemMap["Id"].(string)
		if idOk && domainid == id {
			return itemMap, nil
		}
	}
	return nil, errmsgs.WrapError(errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("Resource UniversalZones:%s not found", id)))
}
