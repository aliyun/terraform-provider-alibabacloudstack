package alibabacloudstack

import (
	"fmt"
	"strings"

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
func (s *UniversalDnsService) DescribeUniversalDNSRecord(id string) (map[string]interface{}, error) {
	parts := strings.SplitN(id, ":", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid id format, expected {ZoneId}:{Id}")
	}
	zoneId := parts[0]
	recordId := parts[1]

	query := map[string]interface{}{
		"ZoneId":     zoneId,
		"PageNumber": 1,
		"PageSize":   100,
	}

	resp, err := s.client.DoTeaRequest("POST", "UniversalDns", "2021-06-24", "DescribeUniversalZoneRecords", "", nil, nil, query)
	if err != nil {
		return nil, err
	}

	data, ok := resp["Data"]
	if !ok || data == nil {
		return nil, errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("record %s not found", id))
	}

	records, ok := data.([]interface{})
	if !ok {
		return nil, errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("record %s not found", id))
	}

	for _, v := range records {
		record := v.(map[string]interface{})
		if record["Id"].(string) == recordId {
			return record, nil
		}
	}

	return nil, errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("record %s not found", id))
}
func (s *UniversalDnsService) DescribeUniversalDnsLine(id string) (map[string]interface{}, error) {
	query := map[string]interface{}{
		"PageNumber": 1,
		"PageSize":   100,
	}

	resp, err := s.client.DoTeaRequest("POST", "UniversalDns", "2021-06-24", "DescribeUniversalLines", "", nil, nil, query)
	if err != nil {
		return nil, err
	}

	if success, ok := resp["success"].(bool); !ok || !success {
		return nil, errmsgs.WrapError(fmt.Errorf("Failed to describe universal lines: %v", resp))
	}

	data, ok := resp["Data"].([]interface{})
	if !ok || data == nil {
		return nil, errmsgs.WrapError(errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("Resource UniversalDnsLine:%s not found", id)))
	}
	for _, item := range data {
		itemMap, isMap := item.(map[string]interface{})
		if !isMap {
			continue
		}
		lineId, idOk := itemMap["Id"].(string)
		if idOk && lineId == id {
			return itemMap, nil
		}
	}
	return nil, errmsgs.WrapError(errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("Resource UniversalDnsLine:%s not found", id)))
}
