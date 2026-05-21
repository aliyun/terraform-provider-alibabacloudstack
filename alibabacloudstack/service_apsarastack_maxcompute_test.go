package alibabacloudstack

import (
	"encoding/json"
	"fmt"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
)

func (s *MaxcomputeService) DescribeMaxcomputeUsers() ([]interface{}, error) {
	action := "GetOdpsUserList"
	request := s.client.NewCommonRequest("POST", "ascm", "2019-05-10", action, "")
	request.SetDomain(s.client.Config.Endpoints[connectivity.ASAPICode])

	response := make(map[string]interface{})
	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(action, bresponse, request, request.QueryParams)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_maxcompute_user", action, errmsgs.AlibabacloudStackSdkGoERROR)
	}
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_maxcompute_user", action, errmsgs.AlibabacloudStackSdkGoERROR)
	}

	datas, err := jsonpath.Get("$.data", response)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, "$.data", response)
	}
	if datas != nil {
		return datas.([]interface{}), nil
	}
	return nil, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
}

func (s *MaxcomputeService) DescribeMaxcomputeUser(id string) (map[string]interface{}, error) {
	action := "GetOdpsUserList"
	request := s.client.NewCommonRequest("POST", "ascm", "2019-05-10", action, "")
	request.SetDomain(s.client.Config.Endpoints[connectivity.ASAPICode])
	request.QueryParams["Id"] = id

	response := make(map[string]interface{})
	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(action, bresponse, request, request.QueryParams)
	if err != nil {
		return response, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_maxcompute_user", action, errmsgs.AlibabacloudStackSdkGoERROR)
	}
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		return response, errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_maxcompute_user", action, errmsgs.AlibabacloudStackSdkGoERROR)
	}

	datas, err := jsonpath.Get("$.data", response)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.data", response)
	}
	for _, v := range datas.([]interface{}) {
		data := v.(map[string]interface{})
		if fmt.Sprint(data["id"]) == id {
			return data, nil
		}
	}
	return nil, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
}

// DescribeMaxcomputeUserForName queries user information by username
func (s *MaxcomputeService) DescribeMaxcomputeUserForName(username string) (map[string]interface{}, error) {
	users, err := s.DescribeMaxcomputeUsers()
	if err != nil {
		return nil, err
	}
	for _, user := range users {
		userMap := user.(map[string]interface{})
		if userMap["userName"].(string) == username {
			return userMap, nil
		}
	}
	return nil, errmsgs.Error(errmsgs.GetNotFoundMessage("Maxcompute User", username))
}

// GetOrCreateMaxcomputeUser retrieves or creates a user by username, returning userId and aasPk
func (s *MaxcomputeService) GetOrCreateMaxcomputeUser(username string) (userId string, aasPk string, err error) {
	// Try to check if the user exists
	userInfo, err := s.DescribeMaxcomputeUserForName(username)
	if err == nil && userInfo != nil {
		// User exists, return id and aasPk
		userId = fmt.Sprintf("%v", userInfo["id"])
		aasPk = fmt.Sprintf("%v", userInfo["aasPk"])
		return userId, aasPk, nil
	}
	// User does not exist, create a new user
	action := "CreateOdpsUser"
	request := s.client.NewCommonRequest("POST", "ascm", "2019-05-10", action, "")
	request.SetDomain(s.client.Config.Endpoints[connectivity.ASAPICode])
	mergeMaps(request.QueryParams, map[string]string{
		"UserName":    username,
		"Description": "Auto created for test",
	})

	response := make(map[string]interface{})
	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(action, bresponse, request, request.QueryParams)
	if err != nil {
		return "", "", errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_maxcompute_user", action, errmsgs.AlibabacloudStackSdkGoERROR)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		return "", "", errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_maxcompute_user", action, errmsgs.AlibabacloudStackSdkGoERROR)
	}

	// Re-query the user list to retrieve the newly created user information
	users, err := s.DescribeMaxcomputeUsers()
	if err != nil || len(users) == 0 {
		return "", "", errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_maxcompute_user", action, errmsgs.AlibabacloudStackSdkGoERROR)
	}

	for _, user := range users {
		userMap := user.(map[string]interface{})
		if userMap["userName"].(string) == username {
			userId = fmt.Sprintf("%v", userMap["id"])
			aasPk = fmt.Sprintf("%v", userMap["aasPk"])
			return userId, aasPk, nil
		}
	}

	return "", "", errmsgs.Error(errmsgs.GetNotFoundMessage("Maxcompute User after creation", username))
}
