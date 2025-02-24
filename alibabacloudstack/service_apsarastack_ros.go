package alibabacloudstack

import (
	"fmt"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type RosService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *RosService) DoRosGettemplateRequest(id string) (object map[string]interface{}, err error) {
	return s.DescribeRosTemplate(id)
}

func (s *RosService) DescribeRosChangeSet(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	request := map[string]interface{}{
		"ChangeSetId":  id,
		"ShowTemplate": true,
	}
	response, err = s.client.DoTeaRequest("POST", "ROS", "2019-09-10", "GetChangeSet", "", nil, nil, request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"ChangeSetNotFound"}) {
			err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("RosChangeSet", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
			return object, err
		}
		return object, err
	}
	addDebug("GetChangeSet", response, request)
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *RosService) RosChangeSetStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeRosChangeSet(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}

		for _, failState := range failStates {
			if object["Status"].(string) == failState {
				return object, object["Status"].(string), errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, object["Status"].(string)))
			}
		}
		return object, object["Status"].(string), nil
	}
}

func (s *RosService) DescribeRosStack(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	request := map[string]interface{}{
		"StackId": id,
	}
	request["ClientToken"] = buildClientToken("GetStack")
	response, err = s.client.DoTeaRequest("POST", "ROS", "2019-09-10", "GetStack", "", nil, nil, request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"StackNotFound"}) {
			err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("RosStack", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
			return object, err
		}
		return object, err
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *RosService) RosStackStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeRosStack(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}

		for _, failState := range failStates {
			if object["Status"].(string) == failState {
				return object, object["Status"].(string), errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, object["Status"].(string)))
			}
		}
		return object, object["Status"].(string), nil
	}
}

func (s *RosService) GetStackPolicy(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	request := map[string]interface{}{
		"StackId": id,
	}
	response, err = s.client.DoTeaRequest("POST", "ROS", "2019-09-10", "GetStackPolicy", "", nil, nil, request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"StackNotFound"}) {
			err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("RosStack", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
			return object, err
		}
		return object, err
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *RosService) ListTagResources(id string, resourceType string) (object interface{}, err error) {
	var response map[string]interface{}
	request := map[string]interface{}{
		"ResourceType": resourceType,
		"ResourceId.1": id,
	}
	tags := make([]interface{}, 0)

	for {
		response, err = s.client.DoTeaRequest("POST", "ROS", "2019-09-10", "ListTagResources", "", nil, nil, request)
		if err != nil {
			return tags, err
		}
		v, err := jsonpath.Get("$.TagResources", response)
		if err != nil {
			return tags, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.TagResources.TagResource", response)
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

func (s *RosService) SetResourceTags(d *schema.ResourceData, resourceType string) error {
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
				"ResourceId.1": d.Id(),
			}
			for i, key := range removedTagKeys {
				request[fmt.Sprintf("TagKey.%d", i+1)] = key
			}
			_, err := s.client.DoTeaRequest("POST", "ROS", "2019-09-10", "UntagResources", "", nil, nil, request)
			if err != nil {
				return err
			}
		}
		if len(added) > 0 {
			request := map[string]interface{}{
				"ResourceType": string(resourceType),
				"ResourceId.1": d.Id(),
			}
			count := 1
			for key, value := range added {
				request[fmt.Sprintf("Tag.%d.Key", count)] = key
				request[fmt.Sprintf("Tag.%d.Value", count)] = value
				count++
			}

			_, err := s.client.DoTeaRequest("POST", "ROS", "2019-09-10", "TagResources", "", nil, nil, request)
			if err != nil {
				return err
			}
		}
		//d.SetPartial("tags")
	}
	return nil
}

func (s *RosService) DescribeRosStackGroup(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	request := map[string]interface{}{
		"StackGroupName": id,
	}
	response, err = s.client.DoTeaRequest("POST", "ROS", "2019-09-10", "GetStackGroup", "", nil, nil, request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"StackGroupNotFound"}) {
			err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("RosStackGroup", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
			return object, err
		}
		return object, err
	}
	addDebug("GetStackGroup", response, request)
	v, err := jsonpath.Get("$.StackGroup", response)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.StackGroup", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *RosService) RosStackGroupStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeRosStackGroup(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}

		for _, failState := range failStates {
			if object["Status"].(string) == failState {
				return object, object["Status"].(string), errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, object["Status"].(string)))
			}
		}
		return object, object["Status"].(string), nil
	}
}

func (s *RosService) DescribeRosTemplate(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	request := map[string]interface{}{
		"TemplateId": id,
	}
	response, err = s.client.DoTeaRequest("POST", "ROS", "2019-09-10", "GetTemplate", "", nil, nil, request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"ChangeSetNotFound", "StackNotFound", "TemplateNotFound"}) {
			err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("RosTemplate", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
			return object, err
		}
		return object, err
	}
	addDebug("GetTemplate", response, request)
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}
