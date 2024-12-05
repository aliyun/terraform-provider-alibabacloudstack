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

func (s *RamService) ParsePolicyDocument(policyDocument string) (statement []map[string]interface{}, version string, err error) {
	policy := Policy{}
	err = json.Unmarshal([]byte(policyDocument), &policy)
	if err != nil {
		err = errmsgs.WrapError(err)
		return
	}

	version = policy.Version
	statement = make([]map[string]interface{}, 0, len(policy.Statement))
	for _, v := range policy.Statement {
		item := make(map[string]interface{})

		item["effect"] = v.Effect
		if val, ok := v.Action.([]interface{}); ok {
			item["action"] = val
		} else {
			item["action"] = []interface{}{v.Action}
		}

		if val, ok := v.Resource.([]interface{}); ok {
			item["resource"] = val
		} else {
			item["resource"] = []interface{}{v.Resource}
		}
		statement = append(statement, item)
	}
	return
}

func (s *RamService) AssembleRolePolicyDocument(ramUser, service []interface{}, version string) (string, error) {
	services := expandStringList(service)
	users := expandStringList(ramUser)

	statement := RolePolicyStatement{
		Effect: Allow,
		Action: "sts:AssumeRole",
		Principal: Principal{
			RAM:     users,
			Service: services,
		},
	}

	policy := RolePolicy{
		Version:   version,
		Statement: []RolePolicyStatement{statement},
	}

	data, err := json.Marshal(policy)
	if err != nil {
		return "", errmsgs.WrapError(err)
	}
	return string(data), nil
}

func (s *RamService) AssemblePolicyDocument(document []interface{}, version string) (string, error) {
	var statements []PolicyStatement

	for _, v := range document {
		doc := v.(map[string]interface{})

		actions := expandStringList(doc["action"].([]interface{}))
		resources := expandStringList(doc["resource"].([]interface{}))

		statement := PolicyStatement{
			Effect:   Effect(doc["effect"].(string)),
			Action:   actions,
			Resource: resources,
		}
		statements = append(statements, statement)
	}

	policy := Policy{
		Version:   version,
		Statement: statements,
	}

	data, err := json.Marshal(policy)
	if err != nil {
		return "", errmsgs.WrapError(err)
	}
	return string(data), nil
}

func (s *RamService) JudgeRolePolicyPrincipal(roleName string) error {
	request := ram.CreateGetRoleRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	request.RoleName = roleName
	raw, err := s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.GetRole(request)
	})
	resp, ok := raw.(*ram.GetRoleResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(resp.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, roleName, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	policy, err := s.ParseRolePolicyDocument(resp.Role.AssumeRolePolicyDocument)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	for _, v := range policy.Statement {
		for _, val := range v.Principal.Service {
			if strings.Trim(val, " ") == "ecs.aliyuncs.com" {
				return nil
			}
		}
	}
	return errmsgs.WrapError(fmt.Errorf("Role policy services must contains 'ecs.aliyuncs.com', Now is \n%v.", resp.Role.AssumeRolePolicyDocument))
}

func (s *RamService) GetIntersection(dataMap []map[string]interface{}, allDataMap map[string]interface{}) (allData []interface{}) {
	for _, v := range dataMap {
		if len(v) > 0 {
			for key := range allDataMap {
				if _, ok := v[key]; !ok {
					allDataMap[key] = nil
				}
			}
		}
	}

	for _, v := range allDataMap {
		if v != nil {
			allData = append(allData, v)
		}
	}
	return
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
		if errmsgs.IsExpectedErrors(err, []string{"EntityNotExist.User"}) {
			return user, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return user, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, getUserRequest.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(getUserRequest.GetActionName(), raw, getUserRequest.RpcRequest, getUserRequest)

	return &response.User, nil
}

func (s *RamService) WaitForRamUser(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeRamUser(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			}
			return errmsgs.WrapError(err)
		}
		if object.UserId == id {
			break
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultTimeoutMsg, id, GetFunc(1), errmsgs.ProviderERROR)
		}
	}
	return nil
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
		if errmsgs.IsExpectedErrors(err, []string{"EntityNotExist"}) {
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

func (s *RamService) WaitForRamGroupMembership(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeRamGroupMembership(id)
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
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, strconv.Itoa(len(object.Users.User)), status, errmsgs.ProviderERROR)
		}
	}
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
		if errmsgs.IsExpectedErrors(err, []string{"EntityNotExist.User.LoginProfile", "EntityNotExist.User"}) {
			return response, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return response, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return response, nil
}

func (s *RamService) WaitForRamLoginProfile(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeRamLoginProfile(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		if object.LoginProfile.UserName == id {
			return nil
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.LoginProfile.UserName, id, errmsgs.ProviderERROR)
		}
	}
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
		if errmsgs.IsExpectedErrors(err, []string{"EntityNotExist.Group"}) {
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

func (s *RamService) WaitForRamGroupPolicyAttachment(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := ParseResourceId(id, 4)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	for {
		object, err := s.DescribeRamGroupPolicyAttachment(id)
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
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.PolicyName, parts[1], errmsgs.ProviderERROR)
		}
	}
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
		if errmsgs.IsExpectedErrors(err, []string{"EntityNotExist"}) {
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
		if errmsgs.IsExpectedErrors(err, []string{"EntityNotExist"}) {
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

func (s *RamService) WaitForRamAccessKey(id, useName string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeRamAccessKey(id, useName)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		if string(status) == object.Status {
			return nil
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.Status, status, errmsgs.ProviderERROR)
		}
	}
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
		if errmsgs.IsExpectedErrors(err, []string{"EntityNotExist.Policy"}) {
			return response, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return response, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	
	return response, nil
}

func (s *RamService) WaitForRamPolicy(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeRamPolicy(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		if object.Policy.PolicyName == id {
			return nil
		}

		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.Policy.PolicyName, id, errmsgs.ProviderERROR)
		}
	}
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
			if errmsgs.IsExpectedErrors(err, []string{"unexpected end of JSON input"}) {
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
		if errmsgs.IsExpectedErrors(err, []string{"InvalidRamRole.NotFound"}) {
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
		if errmsgs.IsExpectedErrors(err, []string{"EntityNotExist.Role"}) {
			return response, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return response, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	return response, nil
}

func (s *RamService) WaitForRamRole(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeRamRole(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		if object.Role.RoleName == id && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.Role.RoleName, id, errmsgs.ProviderERROR)
		}
	}
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
		if errmsgs.IsExpectedErrors(err, []string{"EntityNotExist"}) {
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

func (s *RamService) WaitForRamUserPolicyAttachment(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := ParseResourceId(id, 4)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	for {
		object, err := s.DescribeRamUserPolicyAttachment(id)
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
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.PolicyName, parts[1], errmsgs.ProviderERROR)
		}
	}
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
		if errmsgs.IsExpectedErrors(err, []string{"EntityNotExist.Role"}) {
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

func (s *RamService) WaitForRamRolePolicyAttachment(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	parts, err := ParseResourceId(id, 4)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	for {
		object, err := s.DescribeRamRolePolicyAttachment(id)
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
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.PolicyName, parts[1], errmsgs.ProviderERROR)
		}
	}
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
		if errmsgs.IsExpectedErrors(err, []string{"EntityNotExist.Group"}) {
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

func (s *RamService) WaitForRamGroup(id string, status Status, timeout int) error {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)
	for {
		object, err := s.DescribeRamGroup(id)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				if status == Deleted {
					return nil
				}
			} else {
				return errmsgs.WrapError(err)
			}
		}
		if object.Group.GroupName == id && status != Deleted {
			return nil
		}
		if time.Now().After(deadline) {
			return errmsgs.WrapErrorf(err, errmsgs.WaitTimeoutMsg, id, GetFunc(1), timeout, object.Group.GroupName, id, errmsgs.ProviderERROR)
		}
	}
}

func (s *RamService) DescribeRamAccountPasswordPolicy(id string) (*ram.GetPasswordPolicyResponse, error) {
	response := &ram.GetPasswordPolicyResponse{}
	request := ram.CreateGetPasswordPolicyRequest()
	s.client.InitRpcRequest(*request.RpcRequest)
	raw, err := s.client.WithRamClient(func(ramClient *ram.Client) (interface{}, error) {
		return ramClient.GetPasswordPolicy(request)
	})
	response, ok := raw.(*ram.GetPasswordPolicyResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return response, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	

	return response, nil
}
