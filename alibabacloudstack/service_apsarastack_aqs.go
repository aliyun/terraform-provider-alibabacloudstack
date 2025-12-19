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
