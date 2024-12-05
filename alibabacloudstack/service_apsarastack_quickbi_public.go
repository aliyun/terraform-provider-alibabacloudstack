package alibabacloudstack

import (
	"fmt"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
)

type QuickbiPublicService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *QuickbiPublicService) DescribeQuickBiUser(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	request := map[string]interface{}{
		"UserId": id,
	}
	response, err = s.client.DoTeaRequest("POST", "Quickbi", "2022-03-01", "QueryUserInfoByUserId", "", nil, request)
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

func (s *QuickbiPublicService) QueryUserInfoByUserId(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	request := map[string]interface{}{
		"UserId": id,
	}
	response, err = s.client.DoTeaRequest("POST", "Quickbi", "2022-03-01", "QueryUserInfoByUserId", "", nil, request)
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

func (s *QuickbiPublicService) DescribeQuickBiUserGroup(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	request := map[string]interface{}{
		"UserGroupIds": id,
	}
	response, err = s.client.DoTeaRequest("POST", "Quickbi", "2022-03-01", "ListByUserGroupId", "", nil, request)
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

func (s *QuickbiPublicService) DescribeQuickBiWorkspace(id string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	request := map[string]interface{}{
		"UserId": id,
	}
	response, err = s.client.DoTeaRequest("POST", "Quickbi", "2022-03-01", "QueryWorkspaceUserList", "", nil, request)
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
