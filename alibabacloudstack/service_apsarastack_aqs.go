package alibabacloudstack

import (
	"fmt"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
)

type AqsService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *AqsService) DescribeAqsOssScanConfig(id string) (object map[string]interface{}, err error) {
	request := map[string]interface{}{
		"Id":   id,
		"From": "sas",
	}
	response, err := s.client.DoTeaRequest("GET", "aegis", "2016-11-11", "GetOssScanConfig", "", nil, request, nil)
	if err != nil {
		return nil, err
	}
	data, ok := response["Data"].(map[string]interface{})
	if !ok || data == nil {
		return nil, errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("Aqs OssScanConfig %s not found", id))
	}
	return data, nil
}

func (s *AqsService) DescribeAntiBruteForceRule(id string) (map[string]interface{}, error) {
	request := map[string]interface{}{
		"Id":   id,
		"From": "sas",
	}

	action := "DescribeAntiBruteForceRules"
	response, err := s.client.DoTeaRequest("GET", "aegis", "2016-11-11", action, "", nil, request, nil)
	if err != nil {
		return nil, err
	}

	rules, ok := response["Rules"].([]interface{})
	if !ok || len(rules) == 0 {
		return nil, errmsgs.GetNotFoundErrorFromString("Resource not found")
	}
	for _, v := range rules {
		rule := v.(map[string]interface{})
		if fmt.Sprint(rule["Id"]) == id {
			return rule, nil
		}
	}
	return nil, errmsgs.GetNotFoundErrorFromString("Resource not found")
}

func (s *AqsService) DescribeCloudCenterInstances() ([]interface{}, error) {
	request := map[string]interface{}{
		"From":         "sas",
		"MachineTypes": "ecs",
	}
	response, err := s.client.DoTeaRequest("GET", "aegis", "2016-11-11", "DescribeCloudCenterInstances", "", nil, request, nil)
	if err != nil {
		return nil, err
	}

	instances, ok := response["Instances"].([]interface{})
	if !ok || len(instances) == 0 {
		return nil, errmsgs.GetNotFoundErrorFromString("Resource not found")
	}
	return instances, nil
}
