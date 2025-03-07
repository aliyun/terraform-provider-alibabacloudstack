package alibabacloudstack

import (
	"fmt"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

type NasService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *NasService) DoNasDescribemounttargetsRequest(id string) (object map[string]interface{}, err error) {
	return s.DescribeNasMountTarget(id)
}

func (s *NasService) DoNasDescribeaccessrulesRequest(id string) (object map[string]interface{}, err error) {
	return s.DescribeNasAccessRule(id)
}

func (s *NasService) DoNasDescribeaccessgroupsRequest(id string) (object map[string]interface{}, err error) {
	return s.DescribeNasAccessGroup(id)
}

func (s *NasService) DescribeNasFileSystem(id string) (object map[string]interface{}, err error) {
	request := map[string]interface{}{
		"FileSystemId": id,
		"PageSize":     1,
		"PageNumber":   1,
	}
	response, err := s.client.DoTeaRequest("POST", "Nas", "2017-06-26", "DescribeFileSystems", "", nil, nil, request)
	addDebug("DescribeFileSystems", response, request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidFileSystem.NotFound", "Forbidden.NasNotFound", "Resource.NotFound"}) {
			err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("NasFileSystem", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
			return object, err
		}
		return object, err
	}
	v, err := jsonpath.Get("$.FileSystems.FileSystem", response)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.FileSystems.FileSystem", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("NAS", id)), errmsgs.NotFoundWithResponse, response)
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *NasService) DescribeNasMountTarget(id string) (object map[string]interface{}, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = errmsgs.WrapError(err)
		return
	}
	request := map[string]interface{}{
		"FileSystemId":      parts[0],
		"MountTargetDomain": parts[1],
		"PageSize":          1,
		"PageNumber":        1,
	}
	response, err := s.client.DoTeaRequest("POST", "Nas", "2017-06-26", "DescribeMountTargets", "", nil, nil, request)
	addDebug("DescribeMountTargets", response, request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"Forbidden.NasNotFound", "InvalidFileSystem.NotFound", "InvalidLBid.NotFound", "InvalidMountTarget.NotFound", "VolumeUnavailable"}) {
			err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("NasMountTarget", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
			return object, err
		}
		return object, err
	}
	v, err := jsonpath.Get("$.MountTargets.MountTarget", response)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.MountTargets.MountTarget", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("NAS", id)), errmsgs.NotFoundWithResponse, response)
	} else {
		if v.([]interface{})[0].(map[string]interface{})["MountTargetDomain"].(string) != parts[1] {
			return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("NAS", id)), errmsgs.NotFoundWithResponse, response)
		}
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *NasService) DescribeNasAccessGroup(id string) (object map[string]interface{}, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = errmsgs.WrapError(err)
		return
	}
	request := map[string]interface{}{
		"AccessGroupName": parts[0],
		"FileSystemType":  parts[1],
		"PageSize":        1,
		"PageNumber":      1,
	}
	response, err := s.client.DoTeaRequest("POST", "Nas", "2017-06-26", "DescribeAccessGroups", "", nil, nil, request)
	addDebug("DescribeAccessGroups", response, request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"Forbidden.NasNotFound", "InvalidAccessGroup.NotFound", "Resource.NotFound"}) {
			err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("NasAccessGroup", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
			return object, err
		}
		return object, err
	}
	v, err := jsonpath.Get("$.AccessGroups.AccessGroup", response)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.AccessGroups.AccessGroup", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("NAS", id)), errmsgs.NotFoundWithResponse, response)
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *NasService) DescribeNasAccessRule(id string) (object map[string]interface{}, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = errmsgs.WrapError(err)
		return
	}
	request := map[string]interface{}{
		"AccessGroupName": parts[0],
		"AccessRuleId":    parts[1],
		"PageSize":        1,
		"PageNumber":      1,
	}
	response, err := s.client.DoTeaRequest("POST", "Nas", "2017-06-26", "DescribeAccessRules", "", nil, nil, request)
	addDebug("DescribeAccessRules", response, request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidAccessGroup.NotFound", "Forbidden.NasNotFound"}) {
			err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("AccessRule", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
			return object, err
		}
		return object, err
	}
	v, err := jsonpath.Get("$.AccessRules.AccessRule", response)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.AccessRules.AccessRule", response)
	}
	if len(v.([]interface{})) < 1 {
		return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("NAS", id)), errmsgs.NotFoundWithResponse, response)
	}
	object = v.([]interface{})[0].(map[string]interface{})
	return object, nil
}

func (s *NasService) NasMountTargetStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeNasMountTarget(id)
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

func (s *NasService) DescribeNasFileSystemStateRefreshFunc(id string, defaultRetryState string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeNasFileSystem(id)
		if err != nil {
			if errmsgs.NeedRetry(err) && errmsgs.IsExpectedErrors(err, []string{errmsgs.InvalidFileSystemStatus_Ordering}) {
				return nil, defaultRetryState, nil
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
