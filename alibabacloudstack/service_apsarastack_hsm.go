package alibabacloudstack

import (
	"fmt"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
)

type HsmService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *HsmService) DescribeHsmInstance(id string) (object map[string]interface{}, err error) {

	request := make(map[string]interface{})
	request["InstanceId"] = id

	resp, err := s.client.DoTeaRequest("GET", "hsm-private", "2018-06-30", "DescribeInstances", "", nil, request, nil)
	if err != nil {
		return nil, errmsgs.WrapError(err)
	}

	instances, ok := resp["Instances"].([]interface{})
	if !ok || len(instances) == 0 {
		return nil, errmsgs.GetNotFoundErrorFromString("Resource not found: Hsm Instance " + id)
	}

	object = instances[0].(map[string]interface{})
	if fmt.Sprint(object["HsmStatus"]) == "3" {
		// instance is deleted
		return nil, errmsgs.GetNotFoundErrorFromString("Resource not found: Hsm Instance " + id)
	}
	return object, nil
}
