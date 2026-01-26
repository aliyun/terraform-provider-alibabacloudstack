package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"strings"

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
			if errmsgs.NotFoundError(err) {
				return nil, "", nil
			}

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

type NasDescribelifecyclepoliciesResponse struct {
	LifecyclePolicies []struct {
		FileSystemId        string `json:"FileSystemId"`
		LifecyclePolicyName string `json:"LifecyclePolicyName"`
		Path                string `json:"Path"`
		Recursive           bool   `json:"Recursive"`
		LifecycleRuleName   string `json:"LifecycleRuleName"`
		StorageType         string `json:"StorageType"`
		CreateTime          string `json:"CreateTime"`
		OssBucket           string `json:"OssBucket"`
	} `json:"LifecyclePolicies"`
	RequestId  string `json:"RequestId"`
	TotalCount int    `json:"TotalCount"`
	PageSize   int    `json:"PageSize"`
	PageNumber int    `json:"PageNumber"`
}

func (s *NasService) DoNasDescribelifecyclepoliciesRequest(id string) (*NasDescribelifecyclepoliciesResponse, error) {
	// api: NAS - 2017-06-26 - DescribeLifecyclePolicies
	request := s.client.NewCommonRequest("GET", "NAS", "2017-06-26", "DescribeLifecyclePolicies", "")
	NasDescribelifecyclepoliciesResponseObj := &NasDescribelifecyclepoliciesResponse{}
	params := strings.Split(id, ":")
	request.QueryParams["FileSystemId"] = params[0]
	bresponse, err := s.client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "DescribeLifecyclePolicies", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &NasDescribelifecyclepoliciesResponseObj)

	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", "DescribeLifecyclePolicies", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return NasDescribelifecyclepoliciesResponseObj, nil
}

func (s *NasService) DescribeNasDirQuota(id string) (map[string]interface{}, error) {
	parts := strings.SplitN(id, ":", 2)
	fileSystemId := parts[0]
	path := parts[1]

	request := map[string]interface{}{
		"FileSystemId": fileSystemId,
		"Path":         path,
		"PageSize":     100,
		"PageNumber":   1,
	}

	response, err := s.client.DoTeaRequest("GET", "Nas", "2017-06-26", "DescribeDirQuotas", "", nil, request, nil)
	if err != nil {
		return nil, err
	}

	dirQuotaInfos := response["DirQuotaInfos"].([]interface{})
	if len(dirQuotaInfos) == 0 {
		return nil, errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("Nas dir quota not found with id: %s", id))
	}
	for _, v := range dirQuotaInfos {
		dirQuotaInfo := v.(map[string]interface{})
		if dirQuotaInfo["Path"] == path {
			return dirQuotaInfo, nil
		}
	}

	return nil, errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("Nas dir quota not found with id: %s", id))
}

func (s *NasService) NasDirQuotaStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeNasDirQuota(id)
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

func (s *NasService) DescribeNasNamespace(id string) (map[string]interface{}, error) {
	request := map[string]interface{}{
		"NasNamespaceId": id,
	}
	response, err := s.client.DoTeaRequest("GET", "Nas", "2017-06-26", "DescribeNamespaces", "", nil, request, nil)
	addDebug("DescribeNamespaces", response, request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidNasNamespace.NotFound", "Forbidden.NasNotFound", "Resource.NotFound"}) {
			err = errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("NasNamespace:%s Not found!", id))
			return nil, err
		}
		return nil, err
	}
	nasNamespaces, ok := response["NasNamespaces"].([]interface{})
	if !ok {
		return nil, errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("NasNamespace:%s Not found!", id))
	}
	for _, v := range nasNamespaces {
		nasNamespace := v.(map[string]interface{})
		if nasNamespace["NasNamespaceId"].(string) == id {
			return nasNamespace, nil
		}
	}
	return nil, errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("NasNamespace:%s Not found!", id))
}
func (s *NasService) DescribeNasNamespaceFilesystemAttachment(id string) (map[string]interface{}, error) {
	parts := strings.SplitN(id, ":", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid resource id format, expected NasNamespaceId:FileSystemId")
	}

	nasNamespaceId := parts[0]
	fileSystemId := parts[1]

	reqQuery := map[string]interface{}{
		"NasNamespaceId": nasNamespaceId,
	}

	response, err := s.client.DoTeaRequest("GET", "Nas", "2017-06-26", "ListFileSystemsInNamespace", "", nil, reqQuery, nil)
	if err != nil {
		return nil, err
	}

	nsMembers, ok := response["NSMembers"].([]interface{})
	if !ok {
		return nil, errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("No filesystem found in namespace %s", nasNamespaceId))
	}
	for _, v := range nsMembers {
		member := v.(map[string]interface{})
		if member["FileSystemId"] == fileSystemId {
			return member, nil
		}
	}

	return nil, errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("Filesystem %s not found in namespace %s", fileSystemId, nasNamespaceId))
}

func (s *NasService) DescribeNasNamespaceMountTarget(id string) (map[string]interface{}, error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		err = errmsgs.WrapError(err)
		return nil, err
	}

	request := map[string]interface{}{
		"NasNamespaceId":    parts[0],
		"MountTargetDomain": parts[1],
	}
	response, err := s.client.DoTeaRequest("POST", "Nas", "2017-06-26", "DescribeNamespaceMountTargets", "", nil, nil, request)
	addDebug("DescribeNamespaceMountTargets", response, request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"Forbidden.NasNotFound", "InvalidNasNamespace.NotFound", "InvalidMountTarget.NotFound"}) {
			err = errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("NasNamespaceMountTarget:%s Not found!", id))
			return nil, err
		}
		return nil, err
	}
	mountTargets, ok := response["MountTargets"].([]interface{})
	if !ok {
		return nil, errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("NasNamespaceMountTarget:%s Not found!", id))
	}
	for _, v := range mountTargets {
		if v.(map[string]interface{})["MountTargetDomain"] == parts[1] {
			return v.(map[string]interface{}), nil
		}
	}
	return nil, errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("NasNamespaceMountTarget:%s Not found!", id))
}

func (s *NasService) NasNamespaceMountTargetStateRefreshFunc(id string, failStates []string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		object, err := s.DescribeNasNamespaceMountTarget(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				// Set this to nil as if we didn't find anything.
				return nil, "", nil
			}
			return nil, "", errmsgs.WrapError(err)
		}

		status := object["Status"].(string)
		for _, failState := range failStates {
			if status == failState {
				return object, status, errmsgs.WrapError(errmsgs.Error(errmsgs.FailedToReachTargetStatus, status))
			}
		}
		return object, status, nil
	}
}

func (s *NasService) DescribeNasNamespaceGroup(id string) (map[string]interface{}, error) {
	request := map[string]interface{}{
		"MountTargetDomain": id,
	}
	response, err := s.client.DoTeaRequest("GET", "Nas", "2017-06-26", "DescribeNamespaceGroup", "", nil, request, nil)
	addDebug("DescribeNamespaceGroup", response, request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidNasNamespace.NotFound", "Forbidden.NasNotFound", "Resource.NotFound"}) {
			err = errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("NasNamespaceGroup:%s Not found!", id))
			return nil, err
		}
		return nil, err
	}
	nasNamespaces, ok := response["NasNamespaces"].([]interface{})
	if !ok || len(nasNamespaces) == 0 {
		return nil, errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("NasNamespaceGroup:%s Not found!", id))
	}
	// for _, v := range nasNamespaces {

	// }
	return nasNamespaces[0].(map[string]interface{}), nil
}
