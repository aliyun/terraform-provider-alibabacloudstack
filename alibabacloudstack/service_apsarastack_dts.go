package alibabacloudstack

import (
	"fmt"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type DtsService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *DtsService) DoDtsDescribedtsjobsRequest(id string) (object map[string]interface{}, err error) {
	return s.DescribeDtsSynchronizationInstance(id)
}

func (s *DtsService) DescribeDtsJobMonitorRule(id string) (object map[string]interface{}, err error) {
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"DtsJobId": id,
	}
	response, err := s.client.DoTeaRequest("POST", "Dts", "2020-01-01", "DescribeJobMonitorRule", "", nil, request)
	if err != nil {
		return object, err
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return object, errmsgs.WrapError(fmt.Errorf("%s failed, response: %v", "DescribeJobMonitorRule", response))
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *DtsService) DescribeDtsSubscriptionJob(id string) (object map[string]interface{}, err error) {
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"DtsJobId": id,
	}
	response, err := s.client.DoTeaRequest("POST", "Dts", "2020-01-01", "DescribeDtsJobDetail", "", nil, request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"Forbidden.InstanceNotFound"}) {
			return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("DTS:SubscriptionJob", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, err
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return object, errmsgs.WrapError(fmt.Errorf("%s failed, response: %v", "DescribeDtsJobDetail", response))
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	if object["Status"] == "Starting" {
		object["Status"] = "Normal"
	}
	return object, nil
}

func (s *DtsService) DtsSubscriptionJobStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeDtsSubscriptionJob(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}

func (s *DtsService) DescribeDtsSynchronizationInstance(id string) (object map[string]interface{}, err error) {
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"PageNum":  1,
		"PageSize": 30,
	}
	var response map[string]interface{}
	idExist := false
	for {
		response, err = s.client.DoTeaRequest("POST", "Dts", "2020-01-01", "DescribeSynchronizationJobs", "", nil, request)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"Forbidden.InstanceNotFound"}) {
				return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("DTS:SynchronizationInstance", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR, fmt.Sprint(response["RequestId"]))
			}
			return object, err
		}
		v, err := jsonpath.Get("$.SynchronizationInstances", response)
		if err != nil {
			return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.SynchronizationInstances", response)
		}
		if len(v.([]interface{})) < 1 {
			return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("DTS", id)), errmsgs.NotFoundWithResponse, response)
		}
		for _, v := range v.([]interface{}) {
			if fmt.Sprint(v.(map[string]interface{})["SynchronizationJobId"]) == id {
				idExist = true
				return v.(map[string]interface{}), nil
			}
		}
		if len(v.([]interface{})) < request["PageSize"].(int) {
			break
		}
		request["PageNum"] = request["PageNum"].(int) + 1
	}
	if !idExist {
		return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("DTS", id)), errmsgs.NotFoundWithResponse, response)
	}
	return
}

func (s *DtsService) ListTagResources(id string, resourceType string) (object interface{}, err error) {
	request := map[string]interface{}{
		"ResourceType": resourceType,
		"ResourceId.1": id,
	}
	tags := make([]interface{}, 0)
	var response map[string]interface{}

	for {
		response, err = s.client.DoTeaRequest("POST", "Dts", "2020-01-01", "ListTagResources", "", nil, request)
		if err != nil {
			return object, err
		}
		v, err := jsonpath.Get("$.TagResources.TagResource", response)
		if err != nil {
			return nil, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.TagResources.TagResource", response)
		}
		if v != nil {
			tags = append(tags, v.([]interface{})...)
		}
		if response["NextToken"] == nil {
			break
		}
		request["NextToken"] = response["NextToken"]
	}

	return tags, nil
}

func (s *DtsService) SetResourceTags(d *schema.ResourceData, resourceType string) (err error) {
	if d.HasChange("tags") {
		added, removed := parsingTags(d)
		removedTagKeys := make([]string, 0)
		for _, v := range removed {
			if !ignoredTags(v, "") {
				removedTagKeys = append(removedTagKeys, v)
			}
		}
		if len(removedTagKeys) > 0 {
			request := map[string]interface{}{
				"ResourceType": resourceType,
				"ResourceId.1": d.Get("dts_instance_id"),
			}
			for i, key := range removedTagKeys {
				request[fmt.Sprintf("TagKey.%d", i+1)] = key
			}
			_, err = s.client.DoTeaRequest("POST", "Dts", "2020-01-01", "UntagResources", "", nil, request)
			if err != nil {
				return err
			}
		}
		if len(added) > 0 {
			request := map[string]interface{}{
				"RegionId":     s.client.RegionId,
				"ResourceType": resourceType,
				"ResourceId.1": d.Get("dts_instance_id"),
			}
			count := 1
			for key, value := range added {
				request[fmt.Sprintf("Tag.%d.Key", count)] = key
				request[fmt.Sprintf("Tag.%d.Value", count)] = value
				count++
			}
			_, err = s.client.DoTeaRequest("POST", "Dts", "2020-01-01", "TagResources", "", nil, request)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *DtsService) DescribeDtsSynchronizationJob(id string) (object map[string]interface{}, err error) {
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"DtsJobId": id,
	}
	action := "DescribeDtsJobDetail"
	response, err := s.client.DoTeaRequest("POST", "Dts", "2020-01-01", action, "", nil, request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"Forbidden.InstanceNotFound"}) {
			return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("DTS:SynchronizationJob", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, err
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return object, errmsgs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	if object["Status"] == "synchronizing" || object["Status"] == "Initializing" {
		object["Status"] = "Synchronizing"
	}
	return object, nil
}

func (s *DtsService) DescribeDtsJobDetail(id string) (object map[string]interface{}, err error) {
	request := map[string]interface{}{
		"RegionId": s.client.RegionId,
		"DtsJobId": id,
	}
	action := "DescribeDtsJobDetail"
	response, err := s.client.DoTeaRequest("POST", "Dts", "2020-01-01", action, "", nil, request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"Forbidden.InstanceNotFound"}) {
			return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("DTS:SynchronizationJob", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, err
	}
	if fmt.Sprint(response["Success"]) == "false" {
		return object, errmsgs.WrapError(fmt.Errorf("%s failed, response: %v", action, response))
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *DtsService) DtsSynchronizationJobStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeDtsSynchronizationJob(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}

		for _, failState := range failStates {
			if fmt.Sprint(object["Status"]) == failState {
				return object, fmt.Sprint(object["Status"]), errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, fmt.Sprint(object["Status"])))
			}
		}
		return object, fmt.Sprint(object["Status"]), nil
	}
}
