package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

type WafOpenapiService struct {
	client *connectivity.AlibabacloudStackClient
}

// waf vpc vswitch
type WafVPCVSwitch struct {
	VSwitchName   string `json:"vswitch_name"`
	VSwitch       string `json:"vswitch"`
	CIDRBlock     string `json:"cidr_block"`
	AvailableZone string `json:"available_zone"`
	VPC           string `json:"vpc"`
	VPCName       string `json:"vpc_name"`
}

func (s *WafOpenapiService) convertLogHeadersToString(v []interface{}) (string, error) {
	arrayMaps := make([]interface{}, len(v))
	for i, vv := range v {
		item := vv.(map[string]interface{})
		arrayMaps[i] = map[string]string{
			"k": item["key"].(string),
			"v": item["value"].(string),
		}
	}
	maps, err := json.Marshal(arrayMaps)
	if err != nil {
		return "", errmsgs.WrapError(err)
	}
	return string(maps), nil
}

func (s *WafOpenapiService) DescribeWafDomain(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeDomain"
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = errmsgs.WrapError(err)
		return
	}
	request := map[string]interface{}{
		"RegionId":   s.client.RegionId,
		"Domain":     parts[1],
		"InstanceId": parts[0],
	}
	response, err = client.DoTeaRequest("POST", "waf-openapi", "2019-09-10", action, "", nil, request)
	addDebug(action, response, request)
	v, err := jsonpath.Get("$.Domain", response)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.Domain", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *WafOpenapiService) DescribeWafInstance(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeWAFInstance"
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
	}
	response, err = client.DoTeaRequest("GET", "waf-onecs", "2020-07-01", action, "", nil, request)
	if err != nil {
		err = errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, id, action, errmsgs.AlibabacloudStackSdkGoERROR)
		return
	}
	addDebug(action, response, request)
	// data, err := jsonpath.Get("$.data", response)
	// if err != nil {
	// 	return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.data", response)
	// }
	result, err := jsonpath.Get("$.Result", response)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.Result", response)
	}
	items, err := jsonpath.Get("$.items", result)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.items", response)
	}
	for _, v := range items.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["instance_id"]) == id {
			return v.(map[string]interface{}), nil
		}
	}
	return object, nil
}

func (s *WafOpenapiService) DescribeWafCertificate(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeCertificates"
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		err = errmsgs.WrapError(err)
		return
	}
	request := map[string]interface{}{
		"Domain":     parts[1],
		"InstanceId": parts[0],
	}
	idExist := false
	response, err = client.DoTeaRequest("POST", "waf-openapi", "2019-09-10", action, "", nil, request)
	addDebug(action, response, request)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, id, action, errmsgs.AlibabacloudStackSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Certificates", response)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.Certificates", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("WAF", id)), errmsgs.NotFoundWithResponse, response)
	}
	for _, v := range v.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["CertificateId"]) == parts[2] {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}
	if !idExist {
		return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("WAF", id)), errmsgs.NotFoundWithResponse, response)
	}
	return object, nil
}

func (s *WafOpenapiService) DescribeProtectionModuleStatus(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeProtectionModuleStatus"
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		err = errmsgs.WrapError(err)
		return
	}
	request := map[string]interface{}{
		"DefenseType": parts[2],
		"Domain":      parts[1],
		"InstanceId":  parts[0],
	}
	response, err = client.DoTeaRequest("POST", "waf-openapi", "2019-09-10", action, "", nil, request)
	addDebug(action, response, request)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, id, action, errmsgs.AlibabacloudStackSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *WafOpenapiService) DescribeWafProtectionModule(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	client := s.client
	action := "DescribeProtectionModuleMode"
	parts, err := ParseResourceId(id, 3)
	if err != nil {
		err = errmsgs.WrapError(err)
		return
	}
	request := map[string]interface{}{
		"DefenseType": parts[2],
		"Domain":      parts[1],
		"InstanceId":  parts[0],
	}
	response, err = client.DoTeaRequest("POST", "waf-openapi", "2019-09-10", action, "", nil, request)
	addDebug(action, response, request)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, id, action, errmsgs.AlibabacloudStackSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *WafOpenapiService) DescribeWafv3Instance(id string) (object map[string]interface{}, err error) {
	client := s.client
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
	}

	var response map[string]interface{}
	action := "DescribeInstance"
	response, err = client.DoTeaRequest("POST", "waf-openapi", "2021-10-01", action, "", nil, request)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, id, action, errmsgs.AlibabacloudStackSdkGoERROR)
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$", response)
	}

	if _, ok := v.(map[string]interface{})["InstanceId"]; !ok {
		return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("Wafv3Instance", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
	}

	return v.(map[string]interface{}), nil
}

func (s *WafOpenapiService) Wafv3InstanceStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeWafInstance(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}
		log.Printf("[DEBUG] instance_make_status is %s", object["instance_make_status"])
		status84 := object["instance_make_status"]
		for _, failState := range failStates {
			if fmt.Sprint(status84) == failState {
				return object, fmt.Sprint(status84), errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, fmt.Sprint(status84)))
			}
		}
		return object, fmt.Sprint(status84), nil
	}
}

func (s *WafOpenapiService) DescribeWafv3Domain(id string) (object map[string]interface{}, err error) {
	client := s.client
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return object, errmsgs.WrapError(err)
	}

	request := map[string]interface{}{
		"InstanceId": parts[0],
		"Domain":     parts[1],
		"RegionId":   s.client.RegionId,
	}

	var response map[string]interface{}
	action := "DescribeDomainDetail"
	response, err = client.DoTeaRequest("POST", "waf-openapi", "2021-10-01", action, "", nil, request)
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$", response)
	}

	if _, ok := v.(map[string]interface{})["Domain"]; !ok {
		return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("Wafv3Domain", id)), errmsgs.NotFoundWithResponse, response)
	}

	return v.(map[string]interface{}), nil
}

// func (s *WafOpenapiService) Wafv3DomainStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
// 	return func() (interface{}, string, error) {
// 		object, err := s.DescribeWafv3Domain(id)
// 		if err != nil {
// 			if errmsgs.NotFoundError(err) {
// 				return nil, "", nil
// 			}
// 			return nil, "", errmsgs.WrapError(err)
// 		}

// 		localVar75 := object["Status"]
// 		for _, failState := range failStates {
// 			if fmt.Sprint(localVar75) == failState {
// 				return object, fmt.Sprint(localVar75), errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, fmt.Sprint(localVar75)))
// 			}
// 		}
// 		return object, fmt.Sprint(localVar75), nil
// 	}
// }
