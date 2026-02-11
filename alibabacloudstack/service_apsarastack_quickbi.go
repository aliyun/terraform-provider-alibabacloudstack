package alibabacloudstack

import (
	"fmt"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
)

type QuickbiService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *QuickbiService) DescribeQuickBiUser(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	request := map[string]interface{}{
		"UserId": id,
	}
	response, err = s.client.DoTeaRequest("POST", "quickbi-public", "2022-03-01", "QueryUserInfoByUserId", "", nil, nil, request)
	addDebug("QueryUserInfoByUserId", response, request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"User.Not.In.Organization"}) {
			return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("QuickBI:User", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, err
	}
	v, err := jsonpath.Get("$.Result", response)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.Result", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *QuickbiService) QueryUserInfoByUserId(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	request := map[string]interface{}{
		"UserId": id,
	}
	response, err = s.client.DoTeaRequest("POST", "quickbi-public", "2022-03-01", "QueryUserInfoByUserId", "", nil, nil, request)
	addDebug("QueryUserInfoByUserId", response, request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"User.Not.In.Organization"}) {
			return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("QuickBI:User", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR, fmt.Sprint(response["RequestId"]))
		}
		return object, err
	}
	v, err := jsonpath.Get("$.Result", response)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.Result", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *QuickbiService) DescribeQuickBiUserGroup(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	request := map[string]interface{}{
		"UserGroupIds": id,
	}
	response, err = s.client.DoTeaRequest("POST", "quickbi-public", "2022-03-01", "ListByUserGroupId", "", nil, nil, request)
	addDebug("ListByUserGroupId", response, request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"User.Not.In.Organization"}) {
			return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("QuickBI:User", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
		}
		return object, err
	}
	v, err := jsonpath.Get("$.Result", response)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.Result", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *QuickbiService) DescribeQuickBiWorkspace(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	request := map[string]interface{}{
		"UserId": id,
	}
	response, err = s.client.DoTeaRequest("POST", "quickbi-public", "2022-03-01", "QueryWorkspaceUserList", "", nil, nil, request)
	addDebug("QueryWorkspaceUserList", response, request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"User.Not.In.Organization"}) {
			return object, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("QuickBI:User", id)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
		}
		return object, err
	}
	v, err := jsonpath.Get("$.Result", response)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.Result", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *QuickbiService) DescribeQuickBiUserGroupUser(id string) (object map[string]interface{}, err error) {
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return nil, errmsgs.WrapError(err)
	}
	action := "QueryUserGroupMember"
	request := map[string]interface{}{
		"UserGroupId": parts[0],
	}

	response, err := s.client.DoTeaRequest("POST", "quickbi-public", "2022-03-01", action, "", nil, request, nil)
	if err != nil {
		return object, err
	}
	for _, v := range response["Result"].([]interface{}) {
		user := v.(map[string]interface{})
		if user["Id"].(string) == parts[1] {
			return user, nil
		}
	}
	return object, errmsgs.GetNotFoundErrorFromString("Not found UserGroupUser " + id)
}
