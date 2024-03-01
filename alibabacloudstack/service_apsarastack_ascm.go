package alibabacloudstack

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
)

type AscmService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *AscmService) DescribeAscmLogonPolicy(id string) (response *LoginPolicy, err error) {
	var requestInfo *ecs.Client
	request := requests.NewCommonRequest()
	if s.client.Config.Insecure {
		request.SetHTTPSInsecure(s.client.Config.Insecure)
	}
	request.Method = "POST"
	request.Product = "ascm"
	request.Version = "2019-05-10"
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "ListLoginPolicies"
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.QueryParams = map[string]string{
		"AccessKeySecret": s.client.SecretKey,
		"Product":         "ascm",
		"Department":      s.client.Department,
		"ResourceGroup":   s.client.ResourceGroup,
		"RegionId":        s.client.RegionId,
		"Action":          "ListLoginPolicies",
		"Version":         "2019-05-10",
		"Name":            id,
	}
	var resp = &LoginPolicy{}
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ErrorLoginPolicyNotFound"}) {
			return resp, WrapErrorf(err, NotFoundMsg, AlibabacloudStackSdkGoERROR)
		}
		return resp, WrapErrorf(err, DefaultErrorMsg, id, "ListLoginPolicy", AlibabacloudStackSdkGoERROR)

	}
	addDebug("LoginPolicy", response, requestInfo, request)

	bresponse, _ := raw.(*responses.CommonResponse)
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, WrapError(err)
	}

	if len(resp.Data) < 1 || resp.Code == "200" {
		return resp, WrapError(err)
	}
	return resp, nil
}
func (s *AscmService) DescribeAscmResourceGroup(id string) (response *ResourceGroup, err error) {
	var requestInfo *ecs.Client
	did := strings.Split(id, COLON_SEPARATED)

	request := requests.NewCommonRequest()

	if s.client.Config.Insecure {
		request.SetHTTPSInsecure(s.client.Config.Insecure)
	}
	request.QueryParams = map[string]string{
		"RegionId":          s.client.RegionId,
		"AccessKeySecret":   s.client.SecretKey,
		"Product":           "ascm",
		"Action":            "ListResourceGroup",
		"Version":           "2019-05-10",
		"resourceGroupName": did[0],
	}
	request.Method = "POST"
	request.Product = "Ascm"
	request.Version = "2019-05-10"
	request.ServiceCode = "ascm"
	request.Domain = s.client.Domain
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "ListResourceGroup"
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.RegionId = s.client.RegionId
	var resp = &ResourceGroup{}
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ErrorResourceGroupNotFound"}) {
			return resp, WrapErrorf(err, NotFoundMsg, AlibabacloudStackSdkGoERROR)
		}
		return resp, WrapErrorf(err, DefaultErrorMsg, did[0], "ListResourceGroup", AlibabacloudStackSdkGoERROR)

	}
	addDebug("ListResourceGroup", response, requestInfo, request)

	bresponse, _ := raw.(*responses.CommonResponse)
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, WrapError(err)
	}

	if len(resp.Data) < 1 || resp.Code == "200" {
		return resp, WrapError(err)
	}
	return resp, nil
}
func (s *AscmService) DescribeAscmCustomRole(id string) (response *AscmCustomRole, err error) {
	var requestInfo *ecs.Client
	did := strings.Split(id, COLON_SEPARATED)

	request := requests.NewCommonRequest()

	if s.client.Config.Insecure {
		request.SetHTTPSInsecure(s.client.Config.Insecure)
	}
	request.QueryParams = map[string]string{
		"RegionId":        s.client.RegionId,
		"AccessKeySecret": s.client.SecretKey,
		"Product":         "ascm",
		"Action":          "ListRoles",
		"Version":         "2019-05-10",
		"roleName":        did[0],
		"roleType":        "ROLETYPE_ASCM",
	}
	request.Method = "POST"
	request.Product = "Ascm"
	request.Version = "2019-05-10"
	request.ServiceCode = "ascm"
	request.Domain = s.client.Domain
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "ListRoles"
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.RegionId = s.client.RegionId
	var resp = &AscmCustomRole{}
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ErrorRoleNotFound"}) {
			return resp, WrapErrorf(err, NotFoundMsg, AlibabacloudStackSdkGoERROR)
		}
		return resp, WrapErrorf(err, DefaultErrorMsg, id, "ListRoles", AlibabacloudStackSdkGoERROR)

	}
	addDebug("ListRoles", response, requestInfo, request)

	bresponse, _ := raw.(*responses.CommonResponse)
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, WrapError(err)
	}

	if resp.AsapiErrorCode == "200" {
		return resp, WrapError(err)
	}

	return resp, nil
}
func (s *AscmService) DescribeAscmRamRole(id string) (response *AscmRoles, err error) {
	var requestInfo *ecs.Client
	did := strings.Split(id, COLON_SEPARATED)

	request := requests.NewCommonRequest()

	if s.client.Config.Insecure {
		request.SetHTTPSInsecure(s.client.Config.Insecure)
	}
	request.QueryParams = map[string]string{
		"RegionId":        s.client.RegionId,
		"AccessKeySecret": s.client.SecretKey,
		"Department":      s.client.Department,
		"ResourceGroup":   s.client.ResourceGroup,
		"Product":         "ascm",
		"Action":          "ListRoles",
		"Version":         "2019-05-10",
		"roleName":        did[0],
		"roleType":        "ROLETYPE_RAM",
	}
	request.Method = "POST"
	request.Product = "Ascm"
	request.Version = "2019-05-10"
	request.ServiceCode = "ascm"
	request.Domain = s.client.Domain
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "ListRoles"
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.RegionId = s.client.RegionId
	var resp = &AscmRoles{}
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ErrorRamRoleNotFound"}) {
			return resp, WrapErrorf(err, NotFoundMsg, AlibabacloudStackSdkGoERROR)
		}
		return resp, WrapErrorf(err, DefaultErrorMsg, id, "ListRoles", AlibabacloudStackSdkGoERROR)

	}
	addDebug("ListRoles", response, requestInfo, request)

	bresponse, _ := raw.(*responses.CommonResponse)
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, WrapError(err)
	}

	if resp.AsapiErrorCode == "200" {
		return resp, WrapError(err)
	}

	return resp, nil
}

func (s *AscmService) DescribeAscmRamServiceRole(id string) (response *RamRole, err error) {
	var requestInfo *ecs.Client

	request := requests.NewCommonRequest()

	if s.client.Config.Insecure {
		request.SetHTTPSInsecure(s.client.Config.Insecure)
	}
	request.QueryParams = map[string]string{
		"RegionId":        s.client.RegionId,
		"AccessKeySecret": s.client.SecretKey,
		"Department":      s.client.Department,
		"ResourceGroup":   s.client.ResourceGroup,
		"Product":         "ascm",
		"id":              id,
		"Action":          "GetRAMServiceRole",
		"Version":         "2019-05-10",
		"roleType":        "ROLETYPE_RAM",
	}
	request.Method = "POST"
	request.Product = "Ascm"
	request.Version = "2019-05-10"
	request.ServiceCode = "ascm"
	request.Domain = s.client.Domain
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "ListRAMServiceRoles"
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.RegionId = s.client.RegionId
	var resp = &RamRole{}
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ErrorRamServiceRoleNotFound"}) {
			return resp, WrapErrorf(err, NotFoundMsg, AlibabacloudStackSdkGoERROR)
		}
		return resp, WrapErrorf(err, DefaultErrorMsg, id, "ListRAMServiceRoles", AlibabacloudStackSdkGoERROR)

	}
	addDebug("ListRAMServiceRoles", response, requestInfo, request)

	bresponse, _ := raw.(*responses.CommonResponse)
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, WrapError(err)
	}

	if len(resp.Data) < 1 || resp.Code == "200" {
		return resp, WrapError(err)
	}

	return resp, nil
}

type AscmResourceGroupUser struct {
	CurrentPage     int    `json:"currentPage"`
	PageSize        int    `json:"pageSize"`
	ResourceGroupID int    `json:"resourceGroupId"`
	RGID            int    `json:"resource_group_id"`
	AscmUserIds     string `json:"ascm_user_ids"`
}
type BindResourceAndUsers struct {
	ResourceGroupID int    `json:"resource_group_id"`
	AscmUserIds     string `json:"ascm_user_ids"`
}

func (s *AscmService) DescribeAscmResourceGroupUserAttachment(id string) (response *AscmResourceGroupUser, err error) {
	var requestInfo *ecs.Client
	request := requests.NewCommonRequest()

	if s.client.Config.Insecure {
		request.SetHTTPSInsecure(s.client.Config.Insecure)
	}
	request.QueryParams = map[string]string{
		"RegionId":        s.client.RegionId,
		"AccessKeySecret": s.client.SecretKey,
		"Product":         "ascm",
		"Action":          "ListAscmUsersInsideResourceGroup",
		"Version":         "2019-05-10",
		"resourceGroupId": id,
	}
	request.Method = "POST"
	request.Product = "Ascm"
	request.Version = "2019-05-10"
	request.ServiceCode = "ascm"
	request.Domain = s.client.Domain
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "ListAscmUsersInsideResourceGroup"
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.RegionId = s.client.RegionId
	var resp = &AscmResourceGroupUser{}
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ErrorListAscmUsersInsideResourceGroupNotFound"}) {
			return resp, WrapErrorf(err, NotFoundMsg, AlibabacloudStackSdkGoERROR)
		}
		return resp, WrapErrorf(err, DefaultErrorMsg, id, "ListAscmUsersInsideResourceGroup", AlibabacloudStackSdkGoERROR)

	}
	addDebug("ListAscmUsersInsideResourceGroup", response, requestInfo, request)

	bresponse, _ := raw.(*responses.CommonResponse)
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, WrapError(err)
	}

	if resp.ResourceGroupID != 0 {
		return resp, WrapError(err)
	}

	return resp, nil
}

func (s *AscmService) DescribeAscmUserGroupResourceSet(id string) (response *ListResourceGroup, err error) {
	var requestInfo *ecs.Client
	did := strings.Split(id, COLON_SEPARATED)

	request := requests.NewCommonRequest()

	if s.client.Config.Insecure {
		request.SetHTTPSInsecure(s.client.Config.Insecure)
	}
	if id == "" {
		request.QueryParams = map[string]string{
			"RegionId":        s.client.RegionId,
			"AccessKeySecret": s.client.SecretKey,
			"Product":         "ascm",
			"Action":          "ListResourceGroup",
			"Version":         "2019-05-10",
			"pageSize":        "1000",
		}
	} else {
		request.QueryParams = map[string]string{
			"RegionId":          s.client.RegionId,
			"AccessKeySecret":   s.client.SecretKey,
			"Product":           "ascm",
			"Action":            "ListResourceGroup",
			"Version":           "2019-05-10",
			"resourceGroupName": did[0],
		}
	}

	request.Method = "POST"
	request.Product = "Ascm"
	request.Version = "2019-05-10"
	request.ServiceCode = "ascm"
	request.Domain = s.client.Domain
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "ListResourceGroup"
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.RegionId = s.client.RegionId
	var resp = &ListResourceGroup{}
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ErrorResourceGroupNotFound"}) {
			return resp, WrapErrorf(err, NotFoundMsg, AlibabacloudStackSdkGoERROR)
		}
		return resp, WrapErrorf(err, DefaultErrorMsg, did[0], "ListResourceGroup", AlibabacloudStackSdkGoERROR)

	}
	addDebug("ListResourceGroup", response, requestInfo, request)

	bresponse, _ := raw.(*responses.CommonResponse)
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, WrapError(err)
	}

	if len(resp.Data) < 1 || resp.Code == "200" {
		return resp, WrapError(err)
	}
	return resp, nil
}

func (s *AscmService) DescribeAscmUserGroupResourceSetBinding(id string) (response *ListResourceGroup, err error) {
	var requestInfo *ecs.Client
	request := requests.NewCommonRequest()

	if s.client.Config.Insecure {
		request.SetHTTPSInsecure(s.client.Config.Insecure)
	}
	request.QueryParams = map[string]string{
		"RegionId":        s.client.RegionId,
		"AccessKeySecret": s.client.SecretKey,
		"Product":         "ascm",
		"Action":          "ListResourceGroup",
		"Version":         "2019-05-10",
		"pageSize":        "1000",
	}
	request.Method = "POST"
	request.Product = "Ascm"
	request.Version = "2019-05-10"
	request.ServiceCode = "ascm"
	request.Domain = s.client.Domain
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "ListResourceGroup"
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.RegionId = s.client.RegionId
	var resp = &ListResourceGroup{}
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ErrorListResourceGroupNotFound"}) {
			return resp, WrapErrorf(err, NotFoundMsg, AlibabacloudStackSdkGoERROR)
		}
		return resp, WrapErrorf(err, DefaultErrorMsg, id, "ListResourceGroup", AlibabacloudStackSdkGoERROR)

	}
	addDebug("ListResourceGroup", response, requestInfo, request)

	bresponse, _ := raw.(*responses.CommonResponse)
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, WrapError(err)
	}

	if len(resp.Data) < 1 || resp.Code != "200" {
		return resp, WrapError(err)
	}

	var rgname string
	for i := range resp.Data {
		if strconv.Itoa(resp.Data[i].Id) == id {
			rgname = resp.Data[i].ResourceGroupName
			break
		}
	}
	res, err := s.DescribeAscmUserGroupResourceSet(rgname)

	return res, nil
}

func (s *AscmService) DescribeAscmUser(id string) (response *User, err error) {
	var requestInfo *ecs.Client
	request := requests.NewCommonRequest()
	if s.client.Config.Insecure {
		request.SetHTTPSInsecure(s.client.Config.Insecure)
	}
	request.QueryParams = map[string]string{
		"RegionId":        s.client.RegionId,
		"AccessKeySecret": s.client.SecretKey,
		"Product":         "ascm",
		"Action":          "ListUsers",
		"Version":         "2019-05-10",
		"loginName":       id,
	}
	request.Method = "POST"
	request.Product = "Ascm"
	request.Version = "2019-05-10"
	request.ServiceCode = "ascm"
	request.Domain = s.client.Domain
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "ListUsers"
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.RegionId = s.client.RegionId
	var resp = &User{}
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ErrorUserNotFound"}) {
			return resp, WrapErrorf(err, NotFoundMsg, AlibabacloudStackSdkGoERROR)
		}
		return resp, WrapErrorf(err, DefaultErrorMsg, id, "ListUsers", AlibabacloudStackSdkGoERROR)

	}
	addDebug("ListUsers", response, requestInfo, request)

	bresponse, _ := raw.(*responses.CommonResponse)
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, WrapError(err)
	}

	if len(resp.Data) < 1 || resp.Code == "200" {
		return resp, WrapError(err)
	}

	return resp, nil
}

func (s *AscmService) DescribeAscmUserGroup(id string) (response *UserGroup, err error) {
	var requestInfo *ecs.Client
	request := requests.NewCommonRequest()
	if s.client.Config.Insecure {
		request.SetHTTPSInsecure(s.client.Config.Insecure)
	}
	if id == "" {
		request.QueryParams = map[string]string{
			"RegionId":         s.client.RegionId,
			"AccessKeySecret":  s.client.SecretKey,
			"Product":          "ascm",
			"Action":           "ListUserGroups",
			"Version":          "2019-05-10",
			"SecurityToken":    s.client.Config.SecurityToken,
			"SignatureVersion": "1.0",
			"SignatureMethod":  "HMAC-SHA1",
		}
	} else {
		request.QueryParams = map[string]string{
			"RegionId":         s.client.RegionId,
			"AccessKeySecret":  s.client.SecretKey,
			"Product":          "ascm",
			"Action":           "ListUserGroups",
			"Version":          "2019-05-10",
			"userGroupName":    id,
			"SecurityToken":    s.client.Config.SecurityToken,
			"SignatureVersion": "1.0",
			"SignatureMethod":  "HMAC-SHA1",
		}
	}

	request.Method = "POST"
	request.Product = "Ascm"
	request.Version = "2019-05-10"
	request.ServiceCode = "ascm"
	request.Domain = s.client.Domain
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "ListUserGroups"
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.RegionId = s.client.RegionId
	var resp = &UserGroup{}
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ErrorUserGroupNotFound"}) {
			return resp, WrapErrorf(err, NotFoundMsg, AlibabacloudStackSdkGoERROR)
		}
		return resp, WrapErrorf(err, DefaultErrorMsg, id, "ListUserGroups", AlibabacloudStackSdkGoERROR)

	}
	addDebug("ListUserGroups", response, requestInfo, request)

	bresponse, _ := raw.(*responses.CommonResponse)
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, WrapError(err)
	}

	if len(resp.Data) < 1 || resp.Code != "200" {
		return resp, WrapError(err)
	}

	return resp, nil
}

func (s *AscmService) DescribeAscmUserGroupRoleBinding(id string) (response *UserGroup, err error) {
	var requestInfo *ecs.Client
	request := requests.NewCommonRequest()
	if s.client.Config.Insecure {
		request.SetHTTPSInsecure(s.client.Config.Insecure)
	}
	request.QueryParams = map[string]string{
		"RegionId":         s.client.RegionId,
		"AccessKeySecret":  s.client.SecretKey,
		"Product":          "ascm",
		"Action":           "ListUserGroups",
		"Version":          "2019-05-10",
		"pageSize":         "1000",
		"SecurityToken":    s.client.Config.SecurityToken,
		"SignatureVersion": "1.0",
		"SignatureMethod":  "HMAC-SHA1",
	}
	request.Method = "POST"
	request.Product = "Ascm"
	request.Version = "2019-05-10"
	request.ServiceCode = "ascm"
	request.Domain = s.client.Domain
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "ListUserGroups"
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.RegionId = s.client.RegionId
	var resp = &UserGroup{}
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ErrorUserGroupNotFound"}) {
			return resp, WrapErrorf(err, NotFoundMsg, AlibabacloudStackSdkGoERROR)
		}
		return resp, WrapErrorf(err, DefaultErrorMsg, id, "ListUserGroups", AlibabacloudStackSdkGoERROR)

	}
	addDebug("ListUserGroups", response, requestInfo, request)

	bresponse, _ := raw.(*responses.CommonResponse)
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, WrapError(err)
	}

	if len(resp.Data) < 1 || resp.Code != "200" {
		return resp, WrapError(err)
	}
	var gname string
	for i := range resp.Data {
		if strconv.Itoa(resp.Data[i].Id) == id {
			gname = resp.Data[i].GroupName
			break
		}
	}
	res, err := s.DescribeAscmUserGroup(gname)

	return res, nil
}

func (s *AscmService) DescribeAscmUserRoleBinding(id string) (response *User, err error) {
	var requestInfo *ecs.Client
	request := requests.NewCommonRequest()
	if s.client.Config.Insecure {
		request.SetHTTPSInsecure(s.client.Config.Insecure)
	}
	request.QueryParams = map[string]string{
		"RegionId":        s.client.RegionId,
		"AccessKeySecret": s.client.SecretKey,
		"Product":         "ascm",
		"Action":          "ListUsers",
		"Version":         "2019-05-10",
		"loginName":       id,
	}
	request.Method = "POST"
	request.Product = "Ascm"
	request.Version = "2019-05-10"
	request.ServiceCode = "ascm"
	request.Domain = s.client.Domain
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "ListUsers"
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.RegionId = s.client.RegionId
	var resp = &User{}
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ErrorUserNotFound"}) {
			return resp, WrapErrorf(err, NotFoundMsg, AlibabacloudStackSdkGoERROR)
		}
		return resp, WrapErrorf(err, DefaultErrorMsg, id, "ListUsers", AlibabacloudStackSdkGoERROR)

	}
	addDebug("ListUsers", response, requestInfo, request)

	bresponse, _ := raw.(*responses.CommonResponse)
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, WrapError(err)
	}

	if len(resp.Data) < 1 || resp.Code == "200" {
		return resp, WrapError(err)
	}

	return resp, nil
}
func (s *AscmService) DescribeAscmDeletedUser(id string) (response *DeletedUser, err error) {
	var requestInfo *ecs.Client
	request := requests.NewCommonRequest()
	if s.client.Config.Insecure {
		request.SetHTTPSInsecure(s.client.Config.Insecure)
	}
	request.QueryParams = map[string]string{
		"RegionId":        s.client.RegionId,
		"AccessKeySecret": s.client.SecretKey,
		"Product":         "ascm",
		"Action":          "ListDeletedUsers",
		"Version":         "2019-05-10",
		"loginName":       id,
	}
	request.Method = "POST"
	request.Product = "Ascm"
	request.Version = "2019-05-10"
	request.ServiceCode = "ascm"
	request.Domain = s.client.Domain
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "ListDeletedUsers"
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.RegionId = s.client.RegionId
	var resp = &DeletedUser{}
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ErrorUserNotFound"}) {
			return resp, WrapErrorf(err, NotFoundMsg, AlibabacloudStackSdkGoERROR)
		}
		return resp, WrapErrorf(err, DefaultErrorMsg, id, "ListDeletedUsers", AlibabacloudStackSdkGoERROR)

	}
	addDebug("ListDeletedUsers", response, requestInfo, request)

	bresponse, _ := raw.(*responses.CommonResponse)
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, WrapError(err)
	}
	if resp.Data != nil {
		return resp, WrapError(err)
	}

	return resp, nil
}

func (s *AscmService) DescribeAscmOrganization(id string) (response *Organization, err error) {
	var requestInfo *ecs.Client
	did := strings.Split(id, COLON_SEPARATED)
	request := requests.NewCommonRequest()
	if s.client.Config.Insecure {
		request.SetHTTPSInsecure(s.client.Config.Insecure)
	}
	request.QueryParams = map[string]string{
		"RegionId":        s.client.RegionId,
		"AccessKeySecret": s.client.SecretKey,
		"Department":      s.client.Department,
		"ResourceGroup":   s.client.ResourceGroup,
		"Product":         "ascm",
		"Action":          "GetOrganizationList",
		"Version":         "2019-05-10",
		"name":            did[0],
	}
	request.Method = "POST"
	request.Product = "ascm"
	request.Version = "2019-05-10"
	request.ServiceCode = "ascm"
	request.Domain = s.client.Domain
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "GetOrganizationList"
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.RegionId = s.client.RegionId
	var resp = &Organization{}
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ErrorOrganizationNotFound"}) {
			return resp, WrapErrorf(err, NotFoundMsg, AlibabacloudStackSdkGoERROR)
		}
		return resp, WrapErrorf(err, DefaultErrorMsg, id, "GetOrganization", AlibabacloudStackSdkGoERROR)

	}
	addDebug("GetOrganization", response, requestInfo, request)

	bresponse, _ := raw.(*responses.CommonResponse)
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, WrapError(err)
	}

	if resp.Code == "200" {
		return resp, WrapError(err)
	}

	return resp, nil
}
func (s *AscmService) DescribeAscmRamPolicy(id string) (response *RamPolicies, err error) {
	var requestInfo *ecs.Client
	did := strings.Split(id, COLON_SEPARATED)
	request := requests.NewCommonRequest()
	if s.client.Config.Insecure {
		request.SetHTTPSInsecure(s.client.Config.Insecure)
	}
	request.QueryParams = map[string]string{
		"RegionId":        s.client.RegionId,
		"AccessKeyId":     s.client.AccessKey,
		"AccessKeySecret": s.client.SecretKey,
		"Department":      s.client.Department,
		"ResourceGroup":   s.client.ResourceGroup,
		"Product":         "ascm",
		"Action":          "ListRAMPolicies",
		"Version":         "2019-05-10",
		"policyName":      did[0],
	}
	request.Method = "POST"
	request.Product = "ascm"
	request.Version = "2019-05-10"
	request.ServiceCode = "ascm"
	request.Domain = s.client.Domain
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "ListRAMPolicies"
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.RegionId = s.client.RegionId
	var resp = &RamPolicies{}
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ErrorRamPolicyNotFound"}) {
			return resp, WrapErrorf(err, NotFoundMsg, AlibabacloudStackSdkGoERROR)
		}
		return resp, WrapErrorf(err, DefaultErrorMsg, id, "ListRAMPolicies", AlibabacloudStackSdkGoERROR)

	}
	addDebug("ListRAMPolicies", response, requestInfo, request)

	bresponse, _ := raw.(*responses.CommonResponse)
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, WrapError(err)
	}

	if resp.Code == "200" {
		return resp, WrapError(err)
	}

	return resp, nil
}
func (s *AscmService) DescribeAscmRamPolicyForRole(id string) (response *RamPolicies, err error) {
	var requestInfo *ecs.Client
	did := strings.Split(id, COLON_SEPARATED)

	request := requests.NewCommonRequest()
	if s.client.Config.Insecure {
		request.SetHTTPSInsecure(s.client.Config.Insecure)
	}
	request.QueryParams = map[string]string{
		"RegionId":        s.client.RegionId,
		"AccessKeyId":     s.client.AccessKey,
		"AccessKeySecret": s.client.SecretKey,
		"Department":      s.client.Department,
		"ResourceGroup":   s.client.ResourceGroup,
		"Product":         "ascm",
		"Action":          "ListRAMPolicies",
		"Version":         "2019-05-10",
		"RamPolicyId":     did[0],
		//"roleId":     did[1],
	}
	request.Method = "POST"
	request.Product = "ascm"
	request.Version = "2019-05-10"
	request.ServiceCode = "ascm"
	request.Domain = s.client.Domain
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "ListRAMPolicies"
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.RegionId = s.client.RegionId
	var resp = &RamPolicies{}
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ErrorRamPolicyNotFound"}) {
			return resp, WrapErrorf(err, NotFoundMsg, AlibabacloudStackSdkGoERROR)
		}
		return resp, WrapErrorf(err, DefaultErrorMsg, id, "ListRAMPolicies", AlibabacloudStackSdkGoERROR)

	}
	addDebug("ListRAMPolicies", response, requestInfo, request)

	bresponse, _ := raw.(*responses.CommonResponse)
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, WrapError(err)
	}

	if resp.Code == "200" {
		return resp, WrapError(err)
	}

	return resp, nil
}
func (s *AscmService) DescribeAscmQuota(id string) (response *AscmQuota, err error) {
	var requestInfo *ecs.Client
	did := strings.Split(id, COLON_SEPARATED)
	var targetType string
	if did[0] == "RDS" {
		targetType = "MySql"
	} else if did[0] == "R-KVSTORE" {
		targetType = "redis"
	} else if did[0] == "DDS" {
		targetType = "mongodb"
	} else {
		targetType = ""
	}
	request := requests.NewCommonRequest()
	if s.client.Config.Insecure {
		request.SetHTTPSInsecure(s.client.Config.Insecure)
	}
	request.QueryParams = map[string]string{
		"RegionId":        s.client.RegionId,
		"regionName ":     s.client.RegionId,
		"AccessKeySecret": s.client.SecretKey,
		"Department":      s.client.Department,
		"ResourceGroup":   s.client.ResourceGroup,
		"Product":         "ascm",
		"Action":          "GetQuota",
		"Version":         "2019-05-10",
		"productName":     did[0],
		"quotaType":       did[1],
		"quotaTypeId":     did[2],
		"targetType":      targetType,
	}
	request.Method = "GET"
	request.Product = "Ascm"
	request.Version = "2019-05-10"
	request.ServiceCode = "ascm"
	request.ApiName = "GetQuota"
	request.Domain = s.client.Domain
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.RegionId = s.client.RegionId
	var resp = &AscmQuota{}
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})

	if err != nil {
		if IsExpectedErrors(err, []string{"ErrorQuotaNotFound"}) {
			return resp, WrapErrorf(err, NotFoundMsg, AlibabacloudStackSdkGoERROR)
		}
		return resp, WrapErrorf(err, DefaultErrorMsg, did[0], "GetQuota", AlibabacloudStackSdkGoERROR)

	}
	addDebug("GetQuota", response, requestInfo, request)

	bresponse, _ := raw.(*responses.CommonResponse)
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, WrapError(err)
	}
	if resp.Code == "200" {
		return resp, WrapError(err)
	}

	return resp, nil
}

func (s *AscmService) DescribeAscmPasswordPolicy(id string) (response *PasswordPolicy, err error) {
	var requestInfo *ecs.Client
	//	did := strings.Split(id, COLON_SEPARATED)
	request := requests.NewCommonRequest()
	if s.client.Config.Insecure {
		request.SetHTTPSInsecure(s.client.Config.Insecure)
	}
	request.QueryParams = map[string]string{
		"RegionId":        s.client.RegionId,
		"AccessKeySecret": s.client.SecretKey,
		"Department":      s.client.Department,
		"ResourceGroup":   s.client.ResourceGroup,
		"Product":         "ascm",
		"Action":          "GetPasswordPolicy",
		"Version":         "2019-05-10",
		"id":              id,
	}
	request.Method = "POST"
	request.Product = "ascm"
	request.Version = "2019-05-10"
	request.ServiceCode = "ascm"
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "GetPasswordPolicy"
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.RegionId = s.client.RegionId
	request.Domain = s.client.Domain
	var resp = &PasswordPolicy{}
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ErrorOrganizationNotFound"}) {
			return resp, WrapErrorf(err, NotFoundMsg, AlibabacloudStackSdkGoERROR)
		}
		return resp, WrapErrorf(err, DefaultErrorMsg, id, "GetPasswordPolicy", AlibabacloudStackSdkGoERROR)

	}
	addDebug("GetPasswordPolicy", response, requestInfo, request)

	bresponse, _ := raw.(*responses.CommonResponse)
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, WrapError(err)
	}

	if resp.Code == "200" {
		return resp, WrapError(err)
	}

	return resp, nil
}

func (s *AscmService) DescribeAscmUsergroupUser(id string) (response *User, err error) {
	var requestInfo *ecs.Client
	request := requests.NewCommonRequest()
	if s.client.Config.Insecure {
		request.SetHTTPSInsecure(s.client.Config.Insecure)
	}
	request.QueryParams = map[string]string{
		"AccessKeySecret":  s.client.SecretKey,
		"AccessKeyId":      s.client.AccessKey,
		"Department":       s.client.Department,
		"ResourceGroup":    s.client.ResourceGroup,
		"RegionId":         s.client.RegionId,
		"Product":          "ascm",
		"Action":           "ListUsersInUserGroup",
		"Version":          "2019-05-10",
		"userGroupId":      id,
		"SecurityToken":    s.client.Config.SecurityToken,
		"SignatureVersion": "1.0",
		"SignatureMethod":  "HMAC-SHA1",
	}
	request.Method = "POST"
	request.Product = "Ascm"
	request.Version = "2019-05-10"
	request.ServiceCode = "ascm"
	request.Domain = s.client.Domain
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "ListUsersInUserGroup"
	request.Headers = map[string]string{"RegionId": s.client.RegionId}
	request.RegionId = s.client.RegionId
	var resp = &User{}
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"ErrorUserNotFound"}) {
			return resp, WrapErrorf(err, NotFoundMsg, AlibabacloudStackSdkGoERROR)
		}
		return resp, WrapErrorf(err, DefaultErrorMsg, id, "ListUsersInUserGroup", AlibabacloudStackSdkGoERROR)

	}
	addDebug("ListUsersInUserGroup", response, requestInfo, request)

	bresponse, _ := raw.(*responses.CommonResponse)
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, WrapError(err)
	}

	if len(resp.Data) < 1 || resp.Code == "200" {
		return resp, WrapError(err)
	}

	return resp, nil
}

func (s *AscmService) ExportInitPasswordByLoginName(loginname string) (initPassword string, err error) {
	request := requests.NewCommonRequest()
	if s.client.Config.Insecure {
		request.SetHTTPSInsecure(s.client.Config.Insecure)
	}
	var loginnamelist []string
	loginnamelist = append(loginnamelist, loginname)
	QueryParams := map[string]interface{}{
		"AccessKeySecret":  s.client.SecretKey,
		"AccessKeyId":      s.client.AccessKey,
		"Department":       s.client.Department,
		"ResourceGroup":    s.client.ResourceGroup,
		"RegionId":         s.client.RegionId,
		"Product":          "ascm",
		"Action":           "ExportInitPasswordByLoginNameList",
		"Version":          "2019-05-10",
		"SecurityToken":    s.client.Config.SecurityToken,
		"SignatureVersion": "1.0",
		"SignatureMethod":  "HMAC-SHA1",
		"LoginNameList":    loginnamelist,
	}
	request.Method = "POST"
	request.Product = "Ascm"
	request.Version = "2019-05-10"
	request.ServiceCode = "ascm"
	request.Domain = s.client.Domain
	requeststring, jsonerr := json.Marshal(QueryParams)
	log.Printf("=========================  ExportInitPasswordByLoginNameList jsonerr:%v", jsonerr)
	request.SetContent(requeststring)
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "ExportInitPasswordByLoginNameList"
	request.Headers = map[string]string{
		"RegionId":               s.client.RegionId,
		"x-ascm-product-name":    "ascm",
		"x-ascm-product-version": "2019-05-10",
	}
	request.RegionId = s.client.RegionId
	request.SetContentType("application/json")
	request.PathPattern = "/roa/ascm/auth/user/exportInitPasswordByLoginNameList"
	log.Printf("ExportInitPasswordByLoginNameList loginname:%v", loginname)
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	if err != nil {
		log.Printf("ExportInitPasswordByLoginNameList err:%v", err)
		return initPassword, WrapErrorf(err, DefaultErrorMsg, "", "ExportInitPasswordByLoginNameList", AlibabacloudStackSdkGoERROR)
	}
	bresponse, _ := raw.(*responses.CommonResponse)
	addDebug("ExportInitPasswordByLoginNameList", bresponse, request, loginname)
	var response InitPasswordListResponse
	e := json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	log.Printf("ExportInitPasswordByLoginNameList response:%v", response)
	if e != nil {
		log.Printf("ExportInitPasswordByLoginNameList err:%v", e)
		return initPassword, WrapErrorf(e, DefaultErrorMsg, "", "ExportInitPasswordByLoginNameList", AlibabacloudStackSdkGoERROR)
	}
	if len(response.Data) > 0 {
		initPassword = response.Data[0].Password
	}
	log.Printf("ExportInitPasswordByLoginNameList initPassword:%v", initPassword)
	return initPassword, err
}
