package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
)

type AscmService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *AscmService) DescribeAscmLogonPolicy(id string) (response *LoginPolicy, err error) {
	request := s.client.NewCommonRequest("POST", "ascm", "2019-05-10", "ListLoginPolicies", "")
	request.QueryParams["name"] = id
	var resp = &LoginPolicy{}
	raw, err := s.client.WithAscmClient(func(ascmClient *sdk.Client) (interface{}, error) {
		return ascmClient.ProcessCommonRequest(request)
	})
	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, []string{"ErrorLoginPolicyNotFound"}) {
			return resp, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return resp, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, "ListLoginPolicy", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("LoginPolicy", response, request, request)

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, errmsgs.WrapError(err)
	}

	if len(resp.Data) < 1 || resp.Code == "200" {
		return resp, errmsgs.WrapError(err)
	}
	return resp, nil
}

func (s *AscmService) DescribeAscmResourceGroup(id string) (response *ResourceGroup, err error) {
	did := strings.Split(id, COLON_SEPARATED)
	request := s.client.NewCommonRequest("POST", "Ascm", "2019-05-10", "ListResourceGroup", "")
	request.QueryParams["resourceGroupName"] = did[0]

	var resp = &ResourceGroup{}
	raw, err := s.client.WithAscmClient(func(ascmClient *sdk.Client) (interface{}, error) {
		return ascmClient.ProcessCommonRequest(request)
	})
	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, []string{"ErrorResourceGroupNotFound"}) {
			return resp, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return resp, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, did[0], "ListResourceGroup", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("ListResourceGroup", response, request, request)

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, errmsgs.WrapError(err)
	}

	if len(resp.Data) < 1 || resp.Code == "200" {
		return resp, errmsgs.WrapError(err)
	}
	return resp, nil
}

func (s *AscmService) DescribeAscmCustomRole(id string) (response *AscmCustomRole, err error) {
	did := strings.Split(id, COLON_SEPARATED)
	request := s.client.NewCommonRequest("POST", "Ascm", "2019-05-10", "ListRoles", "")
	request.QueryParams["roleName"] = did[0]
	request.QueryParams["roleType"] = "ROLETYPE_ASCM"

	var resp = &AscmCustomRole{}
	raw, err := s.client.WithAscmClient(func(ascmClient *sdk.Client) (interface{}, error) {
		return ascmClient.ProcessCommonRequest(request)
	})
	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, []string{"ErrorRoleNotFound"}) {
			return resp, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return resp, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, "ListRoles", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("ListRoles", response, request, request)

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, errmsgs.WrapError(err)
	}

	if resp.AsapiErrorCode == "200" {
		return resp, errmsgs.WrapError(err)
	}

	return resp, nil
}

func (s *AscmService) DescribeAscmRamRole(id string) (response *AscmRoles, err error) {
	did := strings.Split(id, COLON_SEPARATED)
	request := s.client.NewCommonRequest("POST", "Ascm", "2019-05-10", "ListRoles", "")
	request.QueryParams["roleName"] = did[0]
	request.QueryParams["roleType"] = "ROLETYPE_RAM"
	var resp = &AscmRoles{}
	raw, err := s.client.WithAscmClient(func(ascmClient *sdk.Client) (interface{}, error) {
		return ascmClient.ProcessCommonRequest(request)
	})
	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, []string{"ErrorRamRoleNotFound"}) {
			return resp, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return resp, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, "ListRoles", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("ListRoles", response, request, request)

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, errmsgs.WrapError(err)
	}

	if resp.AsapiErrorCode == "200" {
		return resp, errmsgs.WrapError(err)
	}

	return resp, nil
}

func (s *AscmService) DescribeAscmRamServiceRole(id string) (response *RamRole, err error) {
	request := s.client.NewCommonRequest("POST", "Ascm", "2019-05-10", "ListRAMServiceRoles", "")
	request.QueryParams["id"] = id
	request.QueryParams["roleType"] = "ROLETYPE_RAM"
	var resp = &RamRole{}
	raw, err := s.client.WithAscmClient(func(ascmClient *sdk.Client) (interface{}, error) {
		return ascmClient.ProcessCommonRequest(request)
	})
	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, []string{"ErrorRamServiceRoleNotFound"}) {
			return resp, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return resp, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, "ListRAMServiceRoles", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("ListRAMServiceRoles", response, request, request)

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, errmsgs.WrapError(err)
	}

	if len(resp.Data) < 1 || resp.Code == "200" {
		return resp, errmsgs.WrapError(err)
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
	request := s.client.NewCommonRequest("POST", "Ascm", "2019-05-10", "ListAscmUsersInsideResourceGroup", "")
	request.QueryParams["resourceGroupId"] = id
	var resp = &AscmResourceGroupUser{}
	raw, err := s.client.WithAscmClient(func(ascmClient *sdk.Client) (interface{}, error) {
		return ascmClient.ProcessCommonRequest(request)
	})
	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, []string{"ErrorListAscmUsersInsideResourceGroupNotFound"}) {
			return resp, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return resp, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, "ListAscmUsersInsideResourceGroup", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("ListAscmUsersInsideResourceGroup", response, request, request)

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, errmsgs.WrapError(err)
	}

	if resp.ResourceGroupID != 0 {
		return resp, errmsgs.WrapError(err)
	}

	return resp, nil
}

func (s *AscmService) DescribeAscmUserGroupResourceSet(id string) (response *ListResourceGroup, err error) {
	did := strings.Split(id, COLON_SEPARATED)
	request := s.client.NewCommonRequest("POST", "Ascm", "2019-05-10", "ListResourceGroup", "")
	if id == "" {
		request.QueryParams["pageSize"] = "1000"
	} else {
		request.QueryParams["resourceGroupName"] = did[0]
	}
	var resp = &ListResourceGroup{}
	raw, err := s.client.WithAscmClient(func(ascmClient *sdk.Client) (interface{}, error) {
		return ascmClient.ProcessCommonRequest(request)
	})
	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, []string{"ErrorResourceGroupNotFound"}) {
			return resp, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return resp, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, did[0], "ListResourceGroup", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("ListResourceGroup", response, request, request)

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, errmsgs.WrapError(err)
	}

	if len(resp.Data) < 1 || resp.Code == "200" {
		return resp, errmsgs.WrapError(err)
	}
	return resp, nil
}

func (s *AscmService) DescribeAscmUserGroupResourceSetBinding(id string) (response *ListResourceGroup, err error) {
	request := s.client.NewCommonRequest("POST", "Ascm", "2019-05-10", "ListResourceGroup", "")
	request.QueryParams["pageSize"] = "1000"
	var resp = &ListResourceGroup{}
	raw, err := s.client.WithAscmClient(func(ascmClient *sdk.Client) (interface{}, error) {
		return ascmClient.ProcessCommonRequest(request)
	})
	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, []string{"ErrorListResourceGroupNotFound"}) {
			return resp, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return resp, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, "ListResourceGroup", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("ListResourceGroup", response, request, request)

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, errmsgs.WrapError(err)
	}

	if len(resp.Data) < 1 || resp.Code != "200" {
		return resp, errmsgs.WrapError(err)
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
	request := s.client.NewCommonRequest("POST", "Ascm", "2019-05-10", "ListUsers", "")
	request.QueryParams["loginName"] = id
	var resp = &User{}
	raw, err := s.client.WithAscmClient(func(ascmClient *sdk.Client) (interface{}, error) {
		return ascmClient.ProcessCommonRequest(request)
	})
	addDebug("ListUsers", raw, request, request.QueryParams)
	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, []string{"ErrorUserNotFound"}) {
			return resp, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return resp, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, "ListUsers", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, errmsgs.WrapError(err)
	}

	if len(resp.Data) < 1 || resp.Code == "200" {
		return resp, errmsgs.WrapError(err)
	}

	return resp, nil
}

func (s *AscmService) DescribeAscmUserGroup(id string) (response *UserGroup, err error) {
	request := s.client.NewCommonRequest("POST", "Ascm", "2019-05-10", "ListUserGroups", "")
	if id != "" {
		request.QueryParams["userGroupName"] = id
	}
	var resp = &UserGroup{}
	raw, err := s.client.WithAscmClient(func(ascmClient *sdk.Client) (interface{}, error) {
		return ascmClient.ProcessCommonRequest(request)
	})
	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, []string{"ErrorUserGroupNotFound"}) {
			return resp, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return resp, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, "ListUserGroups", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("ListUserGroups", response, request, request)

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, errmsgs.WrapError(err)
	}

	if len(resp.Data) < 1 || resp.Code != "200" {
		return resp, errmsgs.WrapError(err)
	}

	return resp, nil
}

func (s *AscmService) DescribeAscmUserGroupRoleBinding(id string) (response *UserGroup, err error) {
	request := s.client.NewCommonRequest("POST", "Ascm", "2019-05-10", "ListUserGroups", "")
	request.QueryParams["pageSize"] = "1000"
	var resp = &UserGroup{}
	raw, err := s.client.WithAscmClient(func(ascmClient *sdk.Client) (interface{}, error) {
		return ascmClient.ProcessCommonRequest(request)
	})
	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, []string{"ErrorUserGroupNotFound"}) {
			return resp, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return resp, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, "ListUserGroups", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("ListUserGroups", response, request, request)

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, errmsgs.WrapError(err)
	}

	if len(resp.Data) < 1 || resp.Code != "200" {
		return resp, errmsgs.WrapError(err)
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
	request := s.client.NewCommonRequest("POST", "Ascm", "2019-05-10", "ListUsers", "")
	request.QueryParams["loginName"] = id
	var resp = &User{}
	raw, err := s.client.WithAscmClient(func(ascmClient *sdk.Client) (interface{}, error) {
		return ascmClient.ProcessCommonRequest(request)
	})
	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, []string{"ErrorUserNotFound"}) {
			return resp, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return resp, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, "ListUsers", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("ListUsers", response, request, request)

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, errmsgs.WrapError(err)
	}

	if len(resp.Data) < 1 || resp.Code == "200" {
		return resp, errmsgs.WrapError(err)
	}

	return resp, nil
}

func (s *AscmService) DescribeAscmDeletedUser(id string) (response *DeletedUser, err error) {
	request := s.client.NewCommonRequest("POST", "Ascm", "2019-05-10", "ListDeletedUsers", "")
	request.QueryParams["loginName"] = id
	var resp = &DeletedUser{}
	raw, err := s.client.WithAscmClient(func(ascmClient *sdk.Client) (interface{}, error) {
		return ascmClient.ProcessCommonRequest(request)
	})
	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, []string{"ErrorUserNotFound"}) {
			return resp, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return resp, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, "ListDeletedUsers", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("ListDeletedUsers", response, request, request)

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, errmsgs.WrapError(err)
	}
	if resp.Data != nil {
		return resp, errmsgs.WrapError(err)
	}

	return resp, nil
}

func (s *AscmService) DescribeAscmOrganization(id string) (response *Organization, err error) {
	did := strings.Split(id, COLON_SEPARATED)
	request := s.client.NewCommonRequest("POST", "Ascm", "2019-05-10", "GetOrganizationList", "")
	request.QueryParams["name"] = did[0]
	var resp = &Organization{}
	raw, err := s.client.WithAscmClient(func(ascmClient *sdk.Client) (interface{}, error) {
		return ascmClient.ProcessCommonRequest(request)
	})
	addDebug("GetOrganization", raw, request, request.QueryParams)
	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, []string{"ErrorOrganizationNotFound"}) {
			return resp, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return resp, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, "GetOrganization", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, errmsgs.WrapError(err)
	}

	if resp.Code == "200" {
		return resp, errmsgs.WrapError(err)
	}

	return resp, nil
}

func (s *AscmService) DescribeAscmRamPolicy(id string) (response *RamPolicies, err error) {
	did := strings.Split(id, COLON_SEPARATED)
	request := s.client.NewCommonRequest("POST", "Ascm", "2019-05-10", "ListRAMPolicies", "")
	request.QueryParams["policyName"] = did[0]
	var resp = &RamPolicies{}
	raw, err := s.client.WithAscmClient(func(ascmClient *sdk.Client) (interface{}, error) {
		return ascmClient.ProcessCommonRequest(request)
	})
	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, []string{"ErrorRamPolicyNotFound"}) {
			return resp, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return resp, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, "ListRAMPolicies", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("ListRAMPolicies", response, request, request)

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, errmsgs.WrapError(err)
	}

	if resp.Code == "200" {
		return resp, errmsgs.WrapError(err)
	}

	return resp, nil
}

func (s *AscmService) DescribeAscmRamPolicyForRole(id string) (response *RamPolicies, err error) {
	did := strings.Split(id, COLON_SEPARATED)
	request := s.client.NewCommonRequest("POST", "Ascm", "2019-05-10", "ListRAMPolicies", "")
	request.QueryParams["RamPolicyId"] = did[0]
	var resp = &RamPolicies{}
	raw, err := s.client.WithAscmClient(func(ascmClient *sdk.Client) (interface{}, error) {
		return ascmClient.ProcessCommonRequest(request)
	})
	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, []string{"ErrorRamPolicyNotFound"}) {
			return resp, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return resp, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, "ListRAMPolicies", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("ListRAMPolicies", response, request, request)

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, errmsgs.WrapError(err)
	}

	if resp.Code == "200" {
		return resp, errmsgs.WrapError(err)
	}

	return resp, nil
}

func (s *AscmService) DescribeAscmQuota(id string) (response *AscmQuota, err error) {
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
	request := s.client.NewCommonRequest("GET", "Ascm", "2019-05-10", "GetQuota", "")
	mergeMaps(request.QueryParams, map[string]string{
		"productName": did[0],
		"quotaType":   did[1],
		"quotaTypeId": did[2],
		"targetType":  targetType,
		"regionName":  s.client.RegionId,
	})
	var resp = &AscmQuota{}
	raw, err := s.client.WithAscmClient(func(ascmClient *sdk.Client) (interface{}, error) {
		return ascmClient.ProcessCommonRequest(request)
	})
	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, []string{"ErrorQuotaNotFound"}) {
			return resp, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return resp, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, did[0], "GetQuota", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("GetQuota", response, request, request)

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, errmsgs.WrapError(err)
	}
	if resp.Code == "200" {
		return resp, errmsgs.WrapError(err)
	}

	return resp, nil
}

func (s *AscmService) DescribeAscmPasswordPolicy(id string) (response *PasswordPolicy, err error) {
	request := s.client.NewCommonRequest("POST", "Ascm", "2019-05-10", "GetPasswordPolicy", "")
	request.QueryParams["id"] = id
	var resp = &PasswordPolicy{}
	raw, err := s.client.WithAscmClient(func(ascmClient *sdk.Client) (interface{}, error) {
		return ascmClient.ProcessCommonRequest(request)
	})
	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, []string{"ErrorOrganizationNotFound"}) {
			return resp, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return resp, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, "GetPasswordPolicy", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("GetPasswordPolicy", response, request, request)

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, errmsgs.WrapError(err)
	}

	if resp.Code == "200" {
		return resp, errmsgs.WrapError(err)
	}

	return resp, nil
}

func (s *AscmService) DescribeAscmUsergroupUser(id string) (response *User, err error) {
	request := s.client.NewCommonRequest("POST", "Ascm", "2019-05-10", "ListUsersInUserGroup", "")
	request.QueryParams["userGroupId"] = id
	var resp = &User{}
	raw, err := s.client.WithAscmClient(func(ascmClient *sdk.Client) (interface{}, error) {
		return ascmClient.ProcessCommonRequest(request)
	})
	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, []string{"ErrorUserNotFound"}) {
			return resp, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return resp, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, "ListUsersInUserGroup", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("ListUsersInUserGroup", response, request, request)

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, errmsgs.WrapError(err)
	}

	if len(resp.Data) < 1 || resp.Code == "200" {
		return resp, errmsgs.WrapError(err)
	}

	return resp, nil
}

func (s *AscmService) ExportInitPasswordByLoginName(loginname string) (initPassword string, err error) {
	var loginnamelist []string
	loginnamelist = append(loginnamelist, loginname)
	request := s.client.NewCommonRequest("POST", "Ascm", "2019-05-10", "ExportInitPasswordByLoginNameList", "")
	loginnamestring, _ := json.Marshal(loginnamelist)
	request.QueryParams["LoginNameList"] = fmt.Sprint(loginnamestring)
	var response InitPasswordListResponse
	raw, err := s.client.WithAscmClient(func(ascmClient *sdk.Client) (interface{}, error) {
		return ascmClient.ProcessCommonRequest(request)
	})
	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		log.Printf("ExportInitPasswordByLoginNameList err:%v", err)
		return initPassword, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "ExportInitPasswordByLoginNameList", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("ExportInitPasswordByLoginNameList", bresponse, request, loginname)
	e := json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	log.Printf("ExportInitPasswordByLoginNameList response:%v", response)
	if e != nil {
		log.Printf("ExportInitPasswordByLoginNameList err:%v", e)
		return initPassword, errmsgs.WrapErrorf(e, errmsgs.DefaultErrorMsg, "", "ExportInitPasswordByLoginNameList", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	if len(response.Data) > 0 {
		initPassword = response.Data[0].Password
	}
	log.Printf("ExportInitPasswordByLoginNameList initPassword:%v", initPassword)
	return initPassword, err
}
