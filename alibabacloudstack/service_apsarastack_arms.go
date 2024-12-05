package alibabacloudstack

import (
	"fmt"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
)

type ArmsService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *ArmsService) DoArmsSearchalertcontactRequest(id string) (object map[string]interface{}, err error) {
	return s.DescribeArmsAlertContact(id)
}

func (s *ArmsService) DescribeArmsAlertContact(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	request := map[string]interface{}{
		"ContactIds": convertListToJsonString([]interface{}{id}),
	}
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	response, err = s.client.DoTeaRequest("POST", "ARMS", "2019-08-08", "SearchAlertContact", "", nil, request)
	if err != nil {
		return object, err
	}
	v, err := jsonpath.Get("$.PageBean.Contacts", response)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.PageBean.Contacts", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("ARMS", id)), errmsgs.NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["ContactId"]) != id {
			return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("ARMS", id)), errmsgs.NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *ArmsService) DescribeArmsAlertContactGroup(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	request := map[string]interface{}{
		"ContactGroupIds": convertListToJsonString([]interface{}{id}),
		"IsDetail":       "true",
	}
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	response, err = s.client.DoTeaRequest("POST", "ARMS", "2019-08-08", "SearchAlertContactGroup", "", nil, request)
	if err != nil {
		return object, err
	}
	v, err := jsonpath.Get("$.ContactGroups", response)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.ContactGroups", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("ARMS", id)), errmsgs.NotFoundWithResponse, response)
	} else {
		if fmt.Sprint(v.([]interface{})[0].(map[string]interface{})["ContactGroupId"]) != id {
			return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("ARMS", id)), errmsgs.NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *ArmsService) DescribeArmsDispatchRule(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	request := map[string]interface{}{
		"Id":       id,
	}
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	response, err = s.client.DoTeaRequest("POST", "ARMS", "2019-08-08", "DescribeDispatchRule", "", nil, request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"50003"}) {
			return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("ARMS", id)), errmsgs.NotFoundWithResponse, response)
		}
		return object, err
	}
	v, err := jsonpath.Get("$.DispatchRule", response)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.DispatchRule", response)
	}
	return v.(map[string]interface{}), nil
}

func (s *ArmsService) DescribeArmsPrometheusAlertRule(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = errmsgs.WrapError(err)
		return
	}
	request := map[string]interface{}{
		"ClusterId": parts[0],
	}
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1
	response, err = s.client.DoTeaRequest("POST", "ARMS", "2019-08-08", "ListPrometheusAlertRules", "", nil, request)
	if err != nil {
		return object, err
	}
	v, err := jsonpath.Get("$.PrometheusAlertRules", response)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.PrometheusAlertRules", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("ARMS", id)), errmsgs.NotFoundWithResponse, response)
	}
	var idExist bool
	for _, v := range v.([]interface{}) {
		if fmt.Sprint(v.(map[string]interface{})["AlertId"]) == parts[1] {
			idExist = true
			return v.(map[string]interface{}), nil
		}
	}
	if !idExist {
		return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("ARMS", id)), errmsgs.NotFoundWithResponse, response)
	}
	return object, nil
}
