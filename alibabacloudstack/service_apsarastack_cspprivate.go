package alibabacloudstack

import (
	"fmt"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
)

type CspprivateService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *CspprivateService) DescribeCspprivateHsmInstance(id string) (object map[string]interface{}, err error) {

	request := make(map[string]interface{})
	request["HsmId"] = id
	request["PageSize"] = 100
	request["PageNumber"] = 1

	resp, err := s.client.DoTeaRequest("POST", "Cspprivate", "2022-02-17", "DescribeHsms", "", nil, request, nil)
	if err != nil {
		return nil, errmsgs.WrapError(err)
	}

	instances, err := jsonpath.Get("$.HsmInfos.HsmInfo", resp)
	if err != nil || len(instances.([]interface{})) == 0 {
		return nil, errmsgs.WrapError(err)
	}
	for _, v := range instances.([]interface{}) {
		instance := v.(map[string]interface{})
		if fmt.Sprint(instance["HsmId"]) == id {
			return instance, nil
		}
	}
	return nil, errmsgs.GetNotFoundErrorFromString("Resource not found: Cspprivate Hsm Instance " + id)
}

func (s *CspprivateService) DescribeProxyCryptoService(id string) (objects []map[string]interface{}, err error) {

	request := map[string]interface{}{
		"CryptoServiceId": "testid",
		"Data":            "{\"HttpMethod\":\"post\",\"ContentType\":\"application/x-www-form-urlencoded\",\"Uri\":\"/api/hsm/describeproducts\",\"Data\":{\"vendorCode\":\"jnta\"}}",
	}
	request["InstanceId"] = id

	resp, err := s.client.DoTeaRequest("GET", "Cspprivate", "2022-02-17", "DescribeProxyCryptoService", "", nil, request, nil)
	if err != nil {
		return nil, errmsgs.WrapError(err)
	}
	vendors, err := jsonpath.Get("$.data.vendors", resp)
	if err != nil || len(vendors.([]interface{})) == 0 {
		return nil, errmsgs.WrapError(err)
	}
	return vendors.([]map[string]interface{}), nil
}

func (s *CspprivateService) DescribeCspprivateHsmGroup(id string) (map[string]interface{}, error) {
	request := map[string]interface{}{
		"PageSize":   100,
		"PageNumber": 1,
	}

	resp, err := s.client.DoTeaRequest("GET", "Cspprivate", "2022-02-17", "DescribeHsmGroups", "", nil, request, nil)
	if err != nil {
		return nil, errmsgs.WrapError(err)
	}

	hsmGroups, err := jsonpath.Get("$.HsmGroups.HsmGroup", resp)
	if err != nil || len(hsmGroups.([]interface{})) == 0 {
		return nil, errmsgs.WrapError(err)
	}
	for _, v := range hsmGroups.([]interface{}) {
		hsmGroup := v.(map[string]interface{})
		if hsmGroup["GroupName"].(string) == id {
			return hsmGroup, nil
		}
	}
	return nil, errmsgs.GetNotFoundErrorFromString("Resource not found: Cspprivate Hsm Group " + id)
}

func (s *CspprivateService) GetCspprivateHsmList(groupId string) (hsm_list []string, err error) {

	request := make(map[string]interface{})
	request["PageSize"] = 100
	request["PageNumber"] = 1
	request["HsmGroup"] = groupId

	resp, err := s.client.DoTeaRequest("POST", "Cspprivate", "2022-02-17", "DescribeHsms", "", nil, request, nil)
	if err != nil {
		return nil, errmsgs.WrapError(err)
	}
	instances, err := jsonpath.Get("$.HsmInfos.HsmInfo", resp)
	if err != nil || len(instances.([]interface{})) == 0 {
		return nil, errmsgs.WrapError(err)
	}
	for _, v := range instances.([]interface{}) {
		instance := v.(map[string]interface{})
		if fmt.Sprint(instance["HsmGroup"]) == groupId {
			hsm_list = append(hsm_list, fmt.Sprint(instance["HsmId"]))
		}
	}
	return hsm_list, nil
}
