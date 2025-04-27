package alibabacloudstack

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"fmt"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
)

type AscmService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *AscmService) DescribeAscmLogonPolicy(id string) (response *LoginPolicy, err error) {
	request := s.client.NewCommonRequest("POST", "ascm", "2019-05-10", "ListLoginPolicies", "/ascm/auth/loginPolicy/listLoginPolicies")
	request.QueryParams["name"] = id
	var resp = &LoginPolicy{}
	bresponse, err := s.client.ProcessCommonRequest(request)

	if err != nil {
		errmsg := ""
		if bresponse != nil {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		} else {
			return nil, err
		}
		if errmsgs.IsExpectedErrors(err, []string{"ErrorLoginPolicyNotFound"}) {
			return resp, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return resp, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, "ListLoginPolicy", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("LoginPolicy", response, request, request.QueryParams)

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
	request := s.client.NewCommonRequest("POST", "ascm", "2019-05-10", "ListResourceGroup", "/ascm/auth/resource_group/list_resource_group")
	request.QueryParams["resourceGroupName"] = did[0]

	var resp = &ResourceGroup{}
	bresponse, err := s.client.ProcessCommonRequest(request)

	if err != nil {
		errmsg := ""
		if bresponse != nil {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		} else {
			return nil, err
		}
		if errmsgs.IsExpectedErrors(err, []string{"ErrorResourceGroupNotFound"}) {
			return resp, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return resp, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, did[0], "ListResourceGroup", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("ListResourceGroup", bresponse, request, request.QueryParams)

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
	request := s.client.NewCommonRequest("POST", "ascm", "2019-05-10", "ListRoles", "/ascm/auth/role/listRoles")
	request.QueryParams["roleName"] = did[0]
	request.QueryParams["roleType"] = "ROLETYPE_ASCM"

	var resp = &AscmCustomRole{}
	bresponse, err := s.client.ProcessCommonRequest(request)

	if err != nil {
		errmsg := ""
		if bresponse != nil {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		} else {
			return nil, err
		}
		if errmsgs.IsExpectedErrors(err, []string{"ErrorRoleNotFound"}) {
			return resp, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return resp, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, "ListRoles", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("ListRoles", response, request, request.QueryParams)

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
	request := s.client.NewCommonRequest("POST", "ascm", "2019-05-10", "ListRoles", "/ascm/auth/role/listRoles")
	request.QueryParams["roleName"] = did[0]
	request.QueryParams["roleType"] = "ROLETYPE_RAM"
	var resp = &AscmRoles{}
	bresponse, err := s.client.ProcessCommonRequest(request)

	if err != nil {
		errmsg := ""
		if bresponse != nil {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		} else {
			return nil, err
		}
		if errmsgs.IsExpectedErrors(err, []string{"ErrorRamRoleNotFound"}) {
			return resp, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return resp, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, "ListRoles", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("ListRoles", response, request, request.QueryParams)

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
	request := s.client.NewCommonRequest("POST", "ascm", "2019-05-10", "ListRAMServiceRoles", "/ascm/auth/role/listRAMServiceRoles")
	request.QueryParams["id"] = id
	request.QueryParams["roleType"] = "ROLETYPE_RAM"
	var resp = &RamRole{}
	bresponse, err := s.client.ProcessCommonRequest(request)

	if err != nil {
		errmsg := ""
		if bresponse != nil {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		} else {
			return nil, err
		}
		if errmsgs.IsExpectedErrors(err, []string{"ErrorRamServiceRoleNotFound"}) {
			return resp, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return resp, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, "ListRAMServiceRoles", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("ListRAMServiceRoles", response, request, request.QueryParams)

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
	request := s.client.NewCommonRequest("POST", "ascm", "2019-05-10", "ListAscmUsersInsideResourceGroup", "/ascm/auth/resource_group/list_ascm_users")
	request.QueryParams["resourceGroupId"] = id
	var resp = &AscmResourceGroupUser{}
	bresponse, err := s.client.ProcessCommonRequest(request)

	if err != nil {
		errmsg := ""
		if bresponse != nil {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		} else {
			return nil, err
		}
		if errmsgs.IsExpectedErrors(err, []string{"ErrorListAscmUsersInsideResourceGroupNotFound"}) {
			return resp, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return resp, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, "ListAscmUsersInsideResourceGroup", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("ListAscmUsersInsideResourceGroup", response, request, request.QueryParams)

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
	request := s.client.NewCommonRequest("POST", "ascm", "2019-05-10", "ListResourceGroup", "/ascm/auth/resource_group/list_resource_group")
	if id == "" {
		request.QueryParams["pageSize"] = "1000"
	} else {
		request.QueryParams["resourceGroupName"] = did[0]
	}
	var resp = &ListResourceGroup{}
	bresponse, err := s.client.ProcessCommonRequest(request)

	if err != nil {
		errmsg := ""
		if bresponse != nil {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		} else {
			return nil, err
		}
		if errmsgs.IsExpectedErrors(err, []string{"ErrorResourceGroupNotFound"}) {
			return resp, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return resp, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, did[0], "ListResourceGroup", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("ListResourceGroup", response, request, request.QueryParams)

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, errmsgs.WrapError(err)
	}

	if len(resp.Data) < 1 || resp.Code == "200" {
		return resp, errmsgs.WrapError(err)
	}
	return resp, nil
}

func (s *AscmService) DescribeAscmUserGroupResourceSetBinding(id string) (*MembersInsideResourceSet, error) {

	var err error
	var resourceSetId, userGroupId string 
	id_infos := strings.Split(id, ":")
	if len(id_infos) == 3 {
		resourceSetId = id_infos[0]
		userGroupId = id_infos[1]
	} else {
		return nil, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
	}
	
	request := s.client.NewCommonRequest("POST", "ascm", "2019-05-10", "ListMembersInsideResourceSet", "/ascm/auth/user/listMembersInsideResourceGroup")
	request.QueryParams["resourceSetId"] = resourceSetId
	
	var resp = &MembersInsideResourceSet{}
	bresponse, err := s.client.ProcessCommonRequest(request)

	if err != nil {
		errmsg := ""
		if bresponse != nil {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		} else {
			return nil,err
		}
		if errmsgs.IsExpectedErrors(err, []string{"ErrorUserGroupNotFound"}) {
			return nil,errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return nil,errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, "ListUserGroups", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("ListUserGroups", bresponse, request, request.QueryParams)

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return nil,errmsgs.WrapError(err)
	}

	if len(resp.Data) < 1 || resp.Code != "200" {
		return nil,errmsgs.WrapError(err)
	}

	for _, data := range(resp.Data) {
		if data.AuthorizedType != "UserGroup" {
			continue
		}
		if strconv.Itoa(data.AuthorizedId) != userGroupId {
			continue
		}
		
		resp.Data = []MembersInsideResourceData{data}
		
		return resp, nil
	}
	return resp, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
}

func (s *AscmService) DescribeAscmUser(id string) (response *User, err error) {
	request := s.client.NewCommonRequest("POST", "ascm", "2019-05-10", "ListUsers", "/ascm/auth/user/listUsers")
	request.QueryParams["loginName"] = id
	var resp = &User{}
	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug("ListUsers", bresponse, request, request.QueryParams)

	if err != nil {
		errmsg := ""
		if bresponse != nil {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		} else {
			return nil, err
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
	request := s.client.NewCommonRequest("POST", "ascm", "2019-05-10", "ListUserGroups", "/ascm/auth/user/listUserGroups")
	if id != "" {
		request.QueryParams["userGroupName"] = id
	}
	var resp = &UserGroup{}
	bresponse, err := s.client.ProcessCommonRequest(request)

	if err != nil {
		errmsg := ""
		if bresponse != nil {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		} else {
			return nil, err
		}
		if errmsgs.IsExpectedErrors(err, []string{"ErrorUserGroupNotFound"}) {
			return resp, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return resp, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, "ListUserGroups", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("ListUserGroups", bresponse, request, request.QueryParams)

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
	request := s.client.NewCommonRequest("POST", "ascm", "2019-05-10", "ListUserGroups", "/ascm/auth/user/listUserGroups")
	request.QueryParams["pageSize"] = "1000"
	var resp = &UserGroup{}
	bresponse, err := s.client.ProcessCommonRequest(request)

	if err != nil {
		errmsg := ""
		if bresponse != nil {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		} else {
			return nil, err
		}
		if errmsgs.IsExpectedErrors(err, []string{"ErrorUserGroupNotFound"}) {
			return resp, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return resp, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, "ListUserGroups", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("ListUserGroups", response, request, request.QueryParams)

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
	request := s.client.NewCommonRequest("POST", "ascm", "2019-05-10", "ListUsers", "/ascm/auth/user/listUsers")
	request.QueryParams["loginName"] = id
	var resp = &User{}
	bresponse, err := s.client.ProcessCommonRequest(request)

	if err != nil {
		errmsg := ""
		if bresponse != nil {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		} else {
			return nil, err
		}
		if errmsgs.IsExpectedErrors(err, []string{"ErrorUserNotFound"}) {
			return resp, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return resp, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, "ListUsers", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("ListUsers", response, request, request.QueryParams)

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
	request := s.client.NewCommonRequest("POST", "ascm", "2019-05-10", "ListDeletedUsers", "/ascm/auth/user/listDeletedUsers")
	request.QueryParams["loginName"] = id
	var resp = &DeletedUser{}
	bresponse, err := s.client.ProcessCommonRequest(request)

	if err != nil {
		errmsg := ""
		if bresponse != nil {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		} else {
			return nil, err
		}
		if errmsgs.IsExpectedErrors(err, []string{"ErrorUserNotFound"}) {
			return resp, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return resp, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, "ListDeletedUsers", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("ListDeletedUsers", response, request, request.QueryParams)

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
	request := s.client.NewCommonRequest("POST", "ascm", "2019-05-10", "GetOrganizationList", "/ascm/auth/organization/queryList")
	request.QueryParams["name"] = did[0]
	var resp = &Organization{}
	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug("GetOrganization", bresponse, request, request.QueryParams)

	if err != nil {
		errmsg := ""
		if bresponse != nil {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		} else {
			return nil, err
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
	request := s.client.NewCommonRequest("POST", "ascm", "2019-05-10", "ListRAMPolicies", "/ascm/auth/role/listRAMPolicies")
	request.QueryParams["policyName"] = did[0]
	var resp = &RamPolicies{}
	bresponse, err := s.client.ProcessCommonRequest(request)

	if err != nil {
		errmsg := ""
		if bresponse != nil {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		} else {
			return nil, err
		}
		if errmsgs.IsExpectedErrors(err, []string{"ErrorRamPolicyNotFound"}) {
			return resp, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return resp, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, "ListRAMPolicies", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("ListRAMPolicies", response, request, request.QueryParams)

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
	request := s.client.NewCommonRequest("POST", "ascm", "2019-05-10", "ListRAMPolicies", "/ascm/auth/role/listRAMPolicies")
	request.QueryParams["RamPolicyId"] = did[0]
	var resp = &RamPolicies{}
	bresponse, err := s.client.ProcessCommonRequest(request)

	if err != nil {
		errmsg := ""
		if bresponse != nil {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		} else {
			return nil, err
		}
		if errmsgs.IsExpectedErrors(err, []string{"ErrorRamPolicyNotFound"}) {
			return resp, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return resp, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, "ListRAMPolicies", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("ListRAMPolicies", response, request, request.QueryParams)

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
	bresponse, err := s.client.ProcessCommonRequest(request)

	if err != nil {
		errmsg := ""
		if bresponse != nil {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		} else {
			return nil, err
		}
		if errmsgs.IsExpectedErrors(err, []string{"ErrorQuotaNotFound"}) {
			return resp, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return resp, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, did[0], "GetQuota", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("GetQuota", response, request, request.QueryParams)

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
	request := s.client.NewCommonRequest("POST", "ascm", "2019-05-10", "GetPasswordPolicy", "/ascm/auth/user/getPasswordPolicy")
	request.QueryParams["id"] = id
	var resp = &PasswordPolicy{}
	bresponse, err := s.client.ProcessCommonRequest(request)

	if err != nil {
		errmsg := ""
		if bresponse != nil {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		} else {
			return nil, err
		}
		if errmsgs.IsExpectedErrors(err, []string{"ErrorOrganizationNotFound"}) {
			return resp, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return resp, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, "GetPasswordPolicy", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("GetPasswordPolicy", response, request, request.QueryParams)

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
	request := s.client.NewCommonRequest("POST", "ascm", "2019-05-10", "ListUsersInUserGroup", "/ascm/auth/user/listUsersInUserGroup")
	request.QueryParams["userGroupId"] = id
	var resp = &User{}
	bresponse, err := s.client.ProcessCommonRequest(request)

	if err != nil {
		errmsg := ""
		if bresponse != nil {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		} else {
			return nil, err
		}
		if errmsgs.IsExpectedErrors(err, []string{"ErrorUserNotFound"}) {
			return resp, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return resp, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, "ListUsersInUserGroup", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("ListUsersInUserGroup", response, request, request.QueryParams)

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
	// 该接口不支持pop网关
	var loginnamelist []string
	loginnamelist = append(loginnamelist, loginname)
	request := s.client.NewCommonRequest("POST", "ascm", "2019-05-10", "ExportInitPasswordByLoginNameList", "/ascm/auth/user/exportInitPasswordByLoginNameList")
	loginnamestring, _ := json.Marshal(loginnamelist)
	request.QueryParams["LoginNameList"] = fmt.Sprint(loginnamestring)
	var response InitPasswordListResponse
	bresponse, err := s.client.ProcessCommonRequest(request)

	if err != nil {
		errmsg := ""
		if bresponse != nil {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		} else {
			return initPassword, err
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
