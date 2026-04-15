package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
		if errmsgs.IsExpectedErrors(err, "ErrorLoginPolicyNotFound") {
			return resp, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return resp, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, "ListLoginPolicy", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("LoginPolicy", bresponse, request, request.QueryParams)

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, errmsgs.WrapError(err)
	}

	if len(resp.Data) < 1 || resp.Code == "200" {
		return resp, errmsgs.WrapError(err)
	}
	return resp, nil
}

func (s *AscmService) DescribeAscmResourceGroup(id string) (result *ResourceGroupData, err error) {
	did := strings.Split(id, COLON_SEPARATED)
	request := s.client.NewCommonRequest("POST", "ascm", "2019-05-10", "ListResourceGroup", "/ascm/auth/resource_group/list_resource_group")
	if len(did) < 2 {
		return nil, fmt.Errorf("invalid id format, expected at least 2 parts separated by colon")
	}
	request.QueryParams["OrganizationId"] = did[0]
	request.QueryParams["Department"] = did[0]
	var resp = &ResourceGroup{}
	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug("ListResourceGroup", bresponse, request, request.QueryParams)
	if err != nil {
		errmsg := ""
		if bresponse != nil {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		} else {
			return nil, err
		}
		if errmsgs.IsExpectedErrors(err, "ErrorResourceGroupNotFound") {
			return nil, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, did[0], "ListResourceGroup", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("ListResourceGroup", bresponse, request, request.QueryParams)

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return nil, errmsgs.WrapError(err)
	}

	if len(resp.Data) < 1 || resp.Code != "200" {
		return nil, errmsgs.GetNotFoundErrorFromString("ResourceGroup not found.")
	}
	for _, v := range resp.Data {
		if did[1] != "" && fmt.Sprint(v.ID) == did[1] {
			return &v, nil
		}
	}
	return nil, errmsgs.GetNotFoundErrorFromString("ResourceGroup not found.")
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
		if errmsgs.IsExpectedErrors(err, "ErrorRoleNotFound") {
			return resp, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return resp, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, "ListRoles", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("ListRoles", bresponse, request, request.QueryParams)

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, errmsgs.WrapError(err)
	}

	if resp.AsapiErrorCode == "200" {
		return resp, errmsgs.WrapError(err)
	}

	return resp, nil
}

func (s *AscmService) DescribeAscmRamRoleForRoleid(id string) (role *AscmRoleDataForGet, err error) {
	did := strings.Split(id, COLON_SEPARATED)
	request := s.client.NewCommonRequest("POST", "ascm", "2019-05-10", "GetRole", "")
	request.QueryParams["roleId"] = did[1]
	response := AscmGetRoleResponse{}
	request.SetDomain(s.client.Config.Endpoints[connectivity.ASAPICode])
	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_ram_role", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		return nil, errmsgs.WrapError(err)
	}
	return &response.Data, nil
}

func (s *AscmService) DescribeAscmRamRole(id string) (role *AscmRoleData, err error) {
	did := strings.Split(id, COLON_SEPARATED)
	request := s.client.NewCommonRequest("POST", "ascm", "2019-05-10", "ListRoles", "/ascm/auth/role/listRoles")
	request.QueryParams["roleName"] = did[0]
	pageSize := 100
	request.QueryParams["pageSize"] = strconv.Itoa(pageSize)
	currentPage := 1
	data := []AscmRoleData{}
	response := ListAscmRolesResponse{}

	for {
		request.QueryParams["currentPage"] = strconv.Itoa(currentPage)
		bresponse, err := s.client.ProcessCommonRequest(request)
		if err != nil {
			if bresponse == nil {
				return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_ram_role", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}

		err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
		if err != nil {
			return nil, errmsgs.WrapError(err)
		}
		data = append(data, response.Data...)
		if response.AsapiErrorCode != "" || response.PageInfo.TotalPage <= currentPage || len(response.Data) < pageSize {
			break
		}
		currentPage += 1
	}

	for _, rg := range data {
		if rg.RoleName == did[0] {
			return &rg, nil
		}
	}

	return nil, errmsgs.GetNotFoundErrorFromString("role " + did[0] + " not found")
}

func (s *AscmService) DescribeAscmRamServiceRole(id string) (response *RamRole, err error) {
	request := s.client.NewCommonRequest("POST", "ascm", "2019-05-10", "ListRAMServiceRoles", "/ascm/auth/role/listRAMServiceRoles")
	request.QueryParams["id"] = id
	var resp = &RamRole{}
	bresponse, err := s.client.ProcessCommonRequest(request)

	if err != nil {
		errmsg := ""
		if bresponse != nil {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		} else {
			return nil, err
		}
		if errmsgs.IsExpectedErrors(err, "ErrorRamServiceRoleNotFound") {
			return resp, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return resp, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, "ListRAMServiceRoles", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("ListRAMServiceRoles", bresponse, request, request.QueryParams)

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

type ListAscmUsersResponse struct {
	Code     string       `json:"code"`
	Cost     int          `json:"cost"`
	Data     []AscmUserv2 `json:"data"`
	Success  bool         `json:"success"`
	PageInfo PageInfo     `json:"pageInfo"`
	Message  string       `json:"message"`
}

type AscmUserv2 struct {
	ID                 int            `json:"id"`
	LoginName          string         `json:"loginName"`
	DisplayName        string         `json:"displayName"`
	Email              string         `json:"email"`
	CellphoneNum       string         `json:"cellphoneNum"`
	Status             string         `json:"status"`
	Deleted            bool           `json:"deleted"`
	Organization       Organizationv2 `json:"organization"`
	DefaultRole        Rolev2         `json:"defaultRole"`
	Roles              []Rolev2       `json:"roles"`
	UserRoles          []Rolev2       `json:"userRoles"`
	LoginPolicy        LoginPolicyv2  `json:"loginPolicy"`
	PrimaryKey         string         `json:"primaryKey"`
	ParentPk           string         `json:"parentPk"`
	EnableEmail        bool           `json:"enableEmail"`
	EnableShortMessage bool           `json:"enableShortMessage"`
	EnableDingTalk     bool           `json:"enableDingTalk"`
	MobileNationCode   string         `json:"mobileNationCode"`
	UserGroups         []interface{}  `json:"userGroups"`
	UserGroupRoles     []interface{}  `json:"userGroupRoles"`
}

type Organizationv2 struct {
	ID                int      `json:"id"`
	Name              string   `json:"name"`
	Alias             string   `json:"alias"`
	Level             string   `json:"level"`
	ParentID          int      `json:"parentId"`
	UUID              string   `json:"uuid"`
	MTime             int64    `json:"mtime"`
	CTime             int64    `json:"ctime"`
	MUserID           string   `json:"muserId"`
	CUserID           string   `json:"cuserId"`
	Internal          bool     `json:"internal"`
	MultiCloudStatus  string   `json:"multiCloudStatus"`
	SupportRegions    string   `json:"supportRegions"`
	SupportRegionList []string `json:"supportRegionList"`
}

type Rolev2 struct {
	ID                     int    `json:"id"`
	Name                   string `json:"roleName"`
	Code                   string `json:"code"`
	Description            string `json:"description"`
	RoleType               string `json:"roleType"`
	RoleRange              string `json:"roleRange"`
	OrganizationVisibility string `json:"organizationVisibility"`
	OwnerOrganizationID    int    `json:"ownerOrganizationId"`
	Level                  int    `json:"roleLevel"`
	Active                 bool   `json:"active"`
	Enable                 bool   `json:"enable"`
	Default                bool   `json:"default"`
	ArID                   string `json:"arId"`
	RAMRole                bool   `json:"rAMRole"`
}

type LoginPolicyv2 struct {
	ID                     int         `json:"id"`
	Name                   string      `json:"name"`
	Rule                   string      `json:"rule"`
	Default                bool        `json:"default"`
	Enable                 bool        `json:"enable"`
	OwnerOrganizationID    int         `json:"ownerOrganizationId"`
	OrganizationVisibility string      `json:"organizationVisibility"`
	IPRanges               []IPRange   `json:"ipRanges"`
	TimeRanges             []TimeRange `json:"timeRanges"`
	MUserID                string      `json:"muserId"`
	CUserID                string      `json:"cuserId"`
	LPID                   string      `json:"lpId"`
}

type IPRange struct {
	ID       int    `json:"loginPolicyId"`
	Protocol string `json:"protocol"`
	IPRange  string `json:"ipRange"`
}

type TimeRange struct {
	ID        int    `json:"loginPolicyId"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

type PageInfo struct {
	Total       int `json:"total"`
	TotalPage   int `json:"totalPage"`
	PageSize    int `json:"pageSize"`
	CurrentPage int `json:"currentPage"`
}

func (s *AscmService) DescribeAscmResourceGroupUserAttachment(rgId string) (*ListAscmUsersResponse, error) {
	client := s.client
	request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "ListAscmUsersInsideResourceGroup", "/ascm/auth/resource_group/list_ascm_users")
	if strings.Contains(rgId, ":") {
		parts := strings.Split(rgId, ":")
		rgId = parts[0]
	}
	request.QueryParams["resourceGroupId"] = rgId

	bresponse, err := client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if bresponse == nil {
			return nil, errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_resource_group_user_attachment", "ListAscmUsersInsideResourceGroup", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	var response ListAscmUsersResponse
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.DataDefaultErrorMsg, "ListAscmUsersInsideResourceGroup", "ListAscmUsersInsideResourceGroup", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	if response.Code != "200" {
		return nil, errmsgs.Error(response.Message)
	}
	return &response, nil
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
		if errmsgs.IsExpectedErrors(err, "ErrorResourceGroupNotFound") {
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
	request.QueryParams["activeOnly"] = "false"
	delete(request.QueryParams, "ResourceGroup")
	delete(request.QueryParams, "OrganizationId")
	delete(request.QueryParams, "Department")

	var resp = &MembersInsideResourceSet{}
	bresponse, err := s.client.ProcessCommonRequest(request)

	if err != nil {
		errmsg := ""
		if bresponse != nil {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		} else {
			return nil, err
		}
		if errmsgs.IsExpectedErrors(err, "ErrorUserGroupNotFound") {
			return nil, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, "ListUserGroups", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("ListUserGroups", bresponse, request, request.QueryParams)

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return nil, errmsgs.WrapError(err)
	}

	if len(resp.Data) < 1 || resp.Code != "200" {
		return nil, errmsgs.WrapError(err)
	}

	for _, data := range resp.Data {
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
	delete(request.QueryParams, "ResourceGroup")
	delete(request.QueryParams, "OrganizationId")
	delete(request.QueryParams, "Department")
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
		if errmsgs.IsExpectedErrors(err, "ErrorUserNotFound") {
			return resp, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return resp, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, "ListUsers", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, errmsgs.WrapError(err)
	}

	if len(resp.Data) < 1 {
		return resp, errmsgs.GetNotFoundErrorFromString("Ascm User not found!")
	}

	return resp, nil
}

func (s *AscmService) DescribeAscmUserGroup(id string) (response *UserGroup, err error) {
	request := s.client.NewCommonRequest("POST", "ascm", "2019-05-10", "ListUserGroups", "/ascm/auth/user/listUserGroups")
	if id != "" {
		request.QueryParams["userGroupName"] = id
		delete(request.QueryParams, "ResourceGroup")
		delete(request.QueryParams, "OrganizationId")
		delete(request.QueryParams, "Department")
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
		if errmsgs.IsExpectedErrors(err, "ErrorUserGroupNotFound") {
			return resp, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return resp, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, "ListUserGroups", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("ListUserGroups", bresponse, request, request.QueryParams)

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, errmsgs.WrapError(err)
	}

	if len(resp.Data) < 1 {
		return resp, errmsgs.GetNotFoundErrorFromString("Ascm usergroup not found!")
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
		if errmsgs.IsExpectedErrors(err, "ErrorUserGroupNotFound") {
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
	parts := strings.Split(id, ":")
	if len(parts) > 2 {
		return nil, fmt.Errorf("Error id format")
	}
	request := s.client.NewCommonRequest("POST", "ascm", "2019-05-10", "ListUsers", "/ascm/auth/user/listUsers")
	request.QueryParams["loginName"] = parts[0]
	var resp = &User{}
	bresponse, err := s.client.ProcessCommonRequest(request)

	if err != nil {
		errmsg := ""
		if bresponse != nil {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		} else {
			return nil, err
		}
		if errmsgs.IsExpectedErrors(err, "ErrorUserNotFound") {
			return resp, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return resp, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, "ListUsers", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("ListUsers", bresponse, request, request.QueryParams)

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, errmsgs.WrapError(err)
	}

	if len(resp.Data) < 1 {
		return nil, errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("Not Found binding for user %s", parts[0]))
	}

	if len(parts) == 2 {
		roleId, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, err
		}
		for _, role := range resp.Data[0].Roles {
			if roleId == role.ID {
				return resp, nil
			}
		}
		return nil, errmsgs.GetNotFoundErrorFromString(fmt.Sprintf("Not Found role %d binding for user %s", roleId, parts[0]))
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
		if errmsgs.IsExpectedErrors(err, "ErrorUserNotFound") {
			return resp, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return resp, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, "ListDeletedUsers", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("ListDeletedUsers", bresponse, request, request.QueryParams)

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, errmsgs.WrapError(err)
	}
	if resp.Data != nil {
		return resp, errmsgs.WrapError(err)
	}

	return resp, nil
}

func (s *AscmService) DescribeAscmOrganizationByName(parentid string, name string) (response *Organization, err error) {
	request := s.client.NewCommonRequest("POST", "ascm", "2019-05-10", "GetOrganizationList", "/ascm/auth/organization/queryList")
	request.QueryParams["id"] = parentid
	var resp = &ListOrganizationResponse{}
	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug("GetOrganization", bresponse, request, request.QueryParams)

	if err != nil {
		errmsg := ""
		if bresponse != nil {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, "ErrorOrganizationNotFound") {
			return nil, errmsgs.GetNotFoundErrorFromString("The Organization is not found")
		}
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, name, "GetOrganization", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return nil, errmsgs.WrapError(err)
	}
	for _, o := range resp.Data {
		if o.Alias == name {
			return &o, nil
		}
	}
	return nil, errmsgs.GetNotFoundErrorFromString("The Organization is not found")
}

func (s *AscmService) DescribeAscmOrganization(id string) (response *Organization, err error) {
	request := s.client.NewCommonRequest("POST", "ascm", "2019-05-10", "GetOrganization", "/ascm/auth/organization/query")
	request.QueryParams["id"] = id
	var resp = &GetOrganizationResponse{}
	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug("GetOrganization", bresponse, request, request.QueryParams)

	if err != nil {
		errmsg := ""
		if bresponse != nil {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, "ErrorOrganizationNotFound") {
			return nil, errmsgs.GetNotFoundErrorFromString("The Organization is not found")
		}
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, "GetOrganization", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return nil, errmsgs.WrapError(err)
	}
	return &resp.Data, nil
}

func (s *AscmService) DescribeAscmRamPolicyForRoleId(id string) (response *RamPolicies, err error) {
	request := s.client.NewCommonRequest("POST", "ascm", "2019-05-10", "GetRole", "")
	request.QueryParams["roleId"] = id
	var resp = &RamPolicies{}
	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug("ListRamPolicies", bresponse, request, request.QueryParams)

	if err != nil {
		errmsg := ""
		if bresponse != nil {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		} else {
			return nil, err
		}
		if errmsgs.IsExpectedErrors(err, "ErrorRamPolicyNotFound") {
			return resp, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return resp, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, "ListRamPolicies", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
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
	request := s.client.NewCommonRequest("POST", "ascm", "2019-05-10", "ListRamPolicies", "/ascm/auth/role/listRAMPolicies")
	request.QueryParams["policyName"] = did[0]
	var resp = &RamPolicies{}
	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug("ListRamPolicies", bresponse, request, request.QueryParams)

	if err != nil {
		errmsg := ""
		if bresponse != nil {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		} else {
			return nil, err
		}
		if errmsgs.IsExpectedErrors(err, "ErrorRamPolicyNotFound") {
			return resp, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return resp, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, "ListRamPolicies", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
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

func (s *AscmService) DescribeAscmRamPolicyForRole(id string) (response *RamPolicies, err error) {
	did := strings.Split(id, COLON_SEPARATED)
	request := s.client.NewCommonRequest("POST", "ascm", "2019-05-10", "ListRamPolicies", "/ascm/auth/role/listRAMPolicies")
	request.QueryParams["RamPolicyId"] = did[0]
	var resp = &RamPolicies{}
	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug("ListRamPolicies", bresponse, request, request.QueryParams)

	if err != nil {
		errmsg := ""
		if bresponse != nil {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		} else {
			return nil, err
		}
		if errmsgs.IsExpectedErrors(err, "ErrorRamPolicyNotFound") {
			return resp, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return resp, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, "ListRamPolicies", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("ListRamPolicies", response, request, request.QueryParams)

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
	request := s.client.NewCommonRequest("GET", "ascm", "2019-05-10", "GetQuota", "/ascm/manage/quota/query")
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
		if errmsgs.IsExpectedErrors(err, "ErrorQuotaNotFound") {
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
		if errmsgs.IsExpectedErrors(err, "ErrorOrganizationNotFound") {
			return resp, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return resp, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, "GetPasswordPolicy", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("GetPasswordPolicy", bresponse, request, request.QueryParams)

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
	delete(request.QueryParams, "ResourceGroup")
	delete(request.QueryParams, "OrganizationId")
	delete(request.QueryParams, "Department")
	body := map[string]interface{}{
		"userGroupId": id,
	}
	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, errmsgs.WrapError(fmt.Errorf("Error marshaling to JSON: %v", err))
	}
	request.SetContentType(requests.Json)
	request.SetContent(jsonData)
	var resp = &User{}
	bresponse, err := s.client.ProcessCommonRequest(request)
	addDebug("ListUsersInUserGroup", bresponse, request, request.QueryParams)
	if err != nil {
		errmsg := ""
		if bresponse != nil {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		if errmsgs.IsExpectedErrors(err, "ErrorUserNotFound") {
			return resp, errmsgs.GetNotFoundErrorFromString("ascm usergroup user not found!")
		}
		return resp, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, id, "ListUsersInUserGroup", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, errmsgs.WrapError(err)
	}

	if len(resp.Data) < 1 {
		return resp, errmsgs.GetNotFoundErrorFromString("ascm usergroup user not found!")
	}

	return resp, nil
}

func (s *AscmService) ExportInitPasswordByLoginName(loginname string) (initPassword string, err error) {
	// This interface does not support pop gateway
	var loginnamelist []string
	loginnamelist = append(loginnamelist, loginname)
	request := s.client.NewCommonRequest("POST", "ascm", "2019-05-10", "ExportInitPasswordByLoginNameList", "/ascm/auth/user/exportInitPasswordByLoginNameList")
	loginnamestring, _ := json.Marshal(loginnamelist)
	request.QueryParams["LoginNameList"] = fmt.Sprint(loginnamestring)
	var response InitPasswordListResponse
	request.SetDomain(s.client.Config.Endpoints[connectivity.ASAPICode])
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

type RAMServiceRoleData struct {
	ID                       int    `json:"id"`
	Arn                      string `json:"arn"`
	Region                   string `json:"region"`
	Product                  string `json:"product"`
	RoleId                   string `json:"roleId"`
	RoleName                 string `json:"roleName"`
	RoleType                 string `json:"roleType"`
	Description              string `json:"description"`
	AliyunUserId             int    `json:"aliyunUserId"`
	AssumeRolePolicyDocument string `json:"assumeRolePolicyDocument"`
	OrganizationId           int    `json:"organizationId"`
	OrganizationName         string `json:"organizationName"`
	Policies                 []struct {
		ID              int    `json:"id"`
		Region          string `json:"region"`
		PolicyName      string `json:"policyName"`
		Description     string `json:"description"`
		PolicyDocument  string `json:"policyDocument"`
		PolicyType      string `json:"policyType"`
		DefaultVersion  string `json:"defaultVersion"`
		AliyunUserId    int    `json:"aliyunUserId"`
		RamGroupId      int    `json:"ramGroupId"`
		AscmRamPolicyId int    `json:"ascmRamPolicyId"`
		AttachDate      int64  `json:"attachDate"`
		ResourceSetId   int    `json:"resourceSetId"`
		PrivilegeId     int    `json:"privilegeId"`
		RamRoleId       int    `json:"ramRoleId"`
	} `json:"policies"`
}

type ListRAMServiceRolesResponse struct {
	Code         string               `json:"code"`
	Cost         int                  `json:"cost"`
	Message      string               `json:"message"`
	PureListData bool                 `json:"pureListData"`
	Redirect     bool                 `json:"redirect"`
	Success      bool                 `json:"success"`
	Data         []RAMServiceRoleData `json:"data"`
	PageInfo     struct {
		CurrentPage int `json:"currentPage"`
		PageSize    int `json:"pageSize"`
		Total       int `json:"total"`
		TotalPage   int `json:"totalPage"`
	} `json:"pageInfo"`
}

func (s *AscmService) ListRAMServiceRoles(id string) (*ListRAMServiceRolesResponse, error) {

	request := s.client.NewCommonRequest("POST", "ascm", "2019-05-10", "ListRAMServiceRoles", "/ascm/auth/role/listRAMServiceRoles")
	params := strings.Split(id, ":")
	request.QueryParams["organizationId"] = params[0]
	request.QueryParams["OrganizationId"] = params[0] // Due to the case-sensitivity strategy of the POP gateway, both orgid and OrgId must be set.
	request.QueryParams["productName"] = params[1]
	var response ListRAMServiceRolesResponse
	bresponse, err := s.client.ProcessCommonRequest(request)
	if err != nil {
		errmsg := ""
		if bresponse != nil {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		} else {
			return nil, err
		}
		log.Printf("ListRAMServiceRoles err:%v", err)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "", "ListRAMServiceRoles", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("ListRAMServiceRoles", bresponse, request, request.QueryParams)
	e := json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if e != nil {
		log.Printf("ListRAMServiceRoles err:%v", e)
		return nil, errmsgs.WrapErrorf(e, errmsgs.DefaultErrorMsg, "", "ListRAMServiceRoles", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	return &response, nil
}

func (s *AscmService) SetResourceTags(d *schema.ResourceData, resourceType string) error {

	if d.HasChange("tags") {
		added, removed := parsingTags(d)

		removedTagKeys := make([]string, 0)
		for _, v := range removed {
			if !ignoredTags(v, "") {
				removedTagKeys = append(removedTagKeys, v)
			}
		}
		if len(removedTagKeys) > 0 {
			action := "UnTagResources"
			request := map[string]interface{}{
				"ResourceType": resourceType,
				"ResourceId.1": d.Id(),
			}
			for i, key := range removedTagKeys {
				request[fmt.Sprintf("TagKey.%d", i+1)] = key
			}
			_, err := s.client.DoTeaRequest("POST", "ascm", "2019-05-10", action, "/ascm/manage/tag_manage/unTagResources", nil, nil, request)
			if err != nil {
				return err
			}
		}
		if len(added) > 0 {
			action := "TagResources"
			request := map[string]interface{}{
				"ResourceType":     resourceType,
				"ResourceId.1":     d.Id(),
				"RegionId":         s.client.RegionId,
				"Department":       s.client.Department,
				"AscmPlatformCode": "default",
			}
			count := 1
			// tagmap := make([]map[string]interface{}, 0)
			for key, value := range added {
				// tagmap = append(tagmap, map[string]interface{}{
				// 	"Key":   key,
				// 	"Value": value,
				// })
				request[fmt.Sprintf("Tag.%d.Key", count)] = key
				request[fmt.Sprintf("Tag.%d.Value", count)] = value
				count++
			}
			// request["Tag"] = tagmap
			_, err := s.client.DoTeaRequest("POST", "ascm", "2019-05-10", action, "/ascm/manage/tag_manage/tagResources", nil, nil, request)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
