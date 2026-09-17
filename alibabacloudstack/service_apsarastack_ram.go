package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ram"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

type Effect string

const (
	Allow Effect = "Allow"
	Deny  Effect = "Deny"
)

type Principal struct {
	Service []string
	RAM     []string
}

type RolePolicyStatement struct {
	Effect    Effect
	Action    string
	Principal Principal
}

type RolePolicy struct {
	Statement []RolePolicyStatement
	Version   string
}

type PolicyStatement struct {
	Effect   Effect
	Action   interface{}
	Resource interface{}
}

type Policy struct {
	Statement []PolicyStatement
	Version   string
}

type RamService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *RamService) ParseRolePolicyDocument(policyDocument string) (RolePolicy, error) {
	var policy RolePolicy
	err := json.Unmarshal([]byte(policyDocument), &policy)
	if err != nil {
		return RolePolicy{}, errmsgs.WrapError(err)
	}
	return policy, nil
}


func (s *RamService) DescribeRamUser(id string) (*ram.User, error) {
	user := &ram.User{}
	listUsersRequest := ram.CreateListUsersRequest()
	s.client.InitRpcRequest(*listUsersRequest.RpcRequest)
	listUsersRequest.MaxItems = requests.NewInteger(100)
	var userName string

	for {
		raw, err := s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
			return ramClient.ListUsers(listUsersRequest)
		})
		response, ok := raw.(*ram.ListUsersResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return user, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, listUsersRequest.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(listUsersRequest.GetActionName(), raw, listUsersRequest.RegionId, listUsersRequest)
		for _, user := range response.Users.User {
			if user.UserId == id {
				userName = user.UserName
				break
			}
		}
		if userName != "" || !response.IsTruncated {
			break
		}
		listUsersRequest.Marker = response.Marker
	}

	if userName == "" {
		// the d.Id() has changed from userName to userId since v1.44.0, add the logic for backward compatibility.
		userName = id
	}
	getUserRequest := ram.CreateGetUserRequest()
	s.client.InitRpcRequest(*getUserRequest.RpcRequest)
	getUserRequest.UserName = userName
	raw, err := s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.GetUser(getUserRequest)
	})
	response, ok := raw.(*ram.GetUserResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, "EntityNotExist.User") {
			return user, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return user, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, getUserRequest.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(getUserRequest.GetActionName(), raw, getUserRequest.RpcRequest, getUserRequest)

	return &response.User, nil
}


func (s *RamService) DescribeRamGroupMembership(id string) (*ram.ListUsersForGroupResponse, error) {
	response := &ram.ListUsersForGroupResponse{}
	request := ram.CreateListUsersForGroupRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.GroupName = id
	raw, err := s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.ListUsersForGroup(request)
	})
	response, ok := raw.(*ram.ListUsersForGroupResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, "EntityNotExist") {
			return response, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return response, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	
	if len(response.Users.User) > 0 {
		return response, nil
	}
	return response, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
}


func (s *RamService) DescribeRamLoginProfile(id string) (*ram.GetLoginProfileResponse, error) {
	response := &ram.GetLoginProfileResponse{}
	request := ram.CreateGetLoginProfileRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.UserName = id

	raw, err := s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.GetLoginProfile(request)
	})
	response, ok := raw.(*ram.GetLoginProfileResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, "EntityNotExist.User.LoginProfile", "EntityNotExist.User") {
			return response, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return response, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return response, nil
}


func (s *RamService) DescribeRamGroupPolicyAttachment(id string) (*ram.Policy, error) {
	response := &ram.Policy{}
	request := ram.CreateListPoliciesForGroupRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	parts, err := ParseResourceId(id, 4)
	if err != nil {
		return response, errmsgs.WrapError(err)
	}
	request.GroupName = parts[3]
	raw, err := s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.ListPoliciesForGroup(request)
	})
	listPoliciesForGroupResponse, ok := raw.(*ram.ListPoliciesForGroupResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(listPoliciesForGroupResponse.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, "EntityNotExist.Group") {
			return response, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return response, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if len(listPoliciesForGroupResponse.Policies.Policy) > 0 {
		for _, v := range listPoliciesForGroupResponse.Policies.Policy {
			if v.PolicyName == parts[1] && v.PolicyType == parts[2] {
				return &v, nil
			}
		}
	}
	return response, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
}


func (s *RamService) DescribeRamAccountAlias(id string) (*ram.GetAccountAliasResponse, error) {
	response := &ram.GetAccountAliasResponse{}
	request := ram.CreateGetAccountAliasRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	raw, err := s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.GetAccountAlias(request)
	})
	response, ok := raw.(*ram.GetAccountAliasResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, "EntityNotExist") {
			return response, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return response, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	return response, nil
}

func (s *RamService) DescribeRamAccessKey(id, userName string) (*ram.AccessKey, error) {
	key := &ram.AccessKey{}
	request := ram.CreateListAccessKeysRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.UserName = userName
	raw, err := s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.ListAccessKeys(request)
	})

	response, ok := raw.(*ram.ListAccessKeysResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, "EntityNotExist") {
			return key, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("RamAccessKey", id)), errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return key, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	
	for _, accessKey := range response.AccessKeys.AccessKey {
		if accessKey.AccessKeyId == id {
			return &accessKey, nil
		}
	}
	return key, errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("RamAccessKey", id)), errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
}


func (s *RamService) DescribeRamPolicy(id string) (*ram.GetPolicyResponse, error) {
	response := &ram.GetPolicyResponse{}
	request := ram.CreateGetPolicyRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.PolicyName = id
	request.PolicyType = "Custom"

	raw, err := s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.GetPolicy(request)
	})
	response, ok := raw.(*ram.GetPolicyResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, "EntityNotExist.Policy") {
			return response, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return response, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	
	return response, nil
}


func (s *RamService) DescribeRamRoleAttachment(id string) (*ecs.DescribeInstanceRamRoleResponse, error) {
	response := &ecs.DescribeInstanceRamRoleResponse{}
	parts, err := ParseResourceId(id, 2)
	if err != nil {
		return response, errmsgs.WrapError(err)
	}
	request := ecs.CreateDescribeInstanceRamRoleRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.InstanceIds = fmt.Sprintf("[\"%s\"]", parts[1])
	var raw interface{}
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err = s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeInstanceRamRole(request)
		})
		if err != nil {
			if errmsgs.IsExpectedErrors(err, "unexpected end of JSON input") {
				return resource.RetryableError(errmsgs.WrapError(err))
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	response, ok := raw.(*ecs.DescribeInstanceRamRoleResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, "InvalidRamRole.NotFound") {
			return response, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return response, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	
	instRoleSets := response.InstanceRamRoleSets.InstanceRamRoleSet
	if len(instRoleSets) > 0 {
		var instIds []string
		for _, item := range instRoleSets {
			if item.RamRoleName == parts[0] {
				instIds = append(instIds, item.InstanceId)
			}
		}
		ids := strings.Split(strings.TrimRight(strings.TrimLeft(strings.Replace(strings.Split(id, ":")[1], "\"", "", -1), "["), "]"), ",")
		sort.Strings(instIds)
		sort.Strings(ids)
		if reflect.DeepEqual(instIds, ids) {
			return response, nil
		}
	}
	return response, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
}

func (s *RamService) WaitForRamRoleAttachment(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeRamRoleAttachment(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		if status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, strconv.Itoa(object.TotalCount), status, errmsgs.ProviderERROR)
		}
	}
}

func (s *RamService) DescribeRamRole(id string) (*ram.GetRoleResponse, error) {
	response := &ram.GetRoleResponse{}
	request := ram.CreateGetRoleRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.RoleName = id
	raw, err := s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.GetRole(request)
	})
	response, ok := raw.(*ram.GetRoleResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(raw.(*ram.GetRoleResponse).BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, "EntityNotExist.Role") {
			return response, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return response, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return response, nil
}


func (s *RamService) DescribeRamUserPolicyAttachment(id string) (*ram.Policy, error) {
	response := &ram.Policy{}
	request := ram.CreateListPoliciesForUserRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	parts, err := ParseResourceId(id, 4)
	if err != nil {
		return response, errmsgs.WrapError(err)
	}
	request.UserName = parts[3]
	raw, err := s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.ListPoliciesForUser(request)
	})
	listPoliciesForUserResponse, ok := raw.(*ram.ListPoliciesForUserResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(raw.(*ram.ListPoliciesForUserResponse).BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, "EntityNotExist") {
			return response, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return response, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	
	if len(listPoliciesForUserResponse.Policies.Policy) > 0 {
		for _, v := range listPoliciesForUserResponse.Policies.Policy {
			if v.PolicyName == parts[1] && v.PolicyType == parts[2] {
				return &v, nil
			}
		}
	}
	return response, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
}


func (s *RamService) DescribeRamRolePolicyAttachment(id string) (*ram.Policy, error) {
	response := &ram.Policy{}
	request := ram.CreateListPoliciesForRoleRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	parts, err := ParseResourceId(id, 4)
	if err != nil {
		return response, errmsgs.WrapError(err)
	}
	request.RoleName = parts[3]
	raw, err := s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.ListPoliciesForRole(request)
	})
	listPoliciesForRoleResponse, ok := raw.(*ram.ListPoliciesForRoleResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(raw.(*ram.ListPoliciesForRoleResponse).BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, "EntityNotExist.Role") {
			return response, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return response, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	if len(listPoliciesForRoleResponse.Policies.Policy) > 0 {
		for _, v := range listPoliciesForRoleResponse.Policies.Policy {
			if v.PolicyName == parts[1] && v.PolicyType == parts[2] {
				return &v, nil
			}
		}
	}
	return response, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
}


func (s *RamService) DescribeRamGroup(id string) (*ram.GetGroupResponse, error) {
	response := &ram.GetGroupResponse{}
	request := ram.CreateGetGroupRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.GroupName = id
	raw, err := s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.GetGroup(request)
	})
	response, ok := raw.(*ram.GetGroupResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, "EntityNotExist.Group") {
			return response, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return response, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	

	if response.Group.GroupName != id {
		return response, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
	}
	return response, nil
}
