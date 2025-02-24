package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
)

type MaxcomputeService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *MaxcomputeService) DescribeMaxcomputeProject(name string) (object *MaxComputeProject, err error) {
	client := s.client

	roleId, err := client.RoleIds()
	if err != nil {
		err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("ASCM User", "defaultRoleId")), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
		return nil, err
	}

	request := make(map[string]interface{})
	request["ResourceGroupId"] = client.ResourceGroup
	request["CalcEngineType"] = "ODPS" // 固定值
	request["OrganizationId"] = client.Department
	request["Department"] = client.Department
	request["ResourceGroup"] = client.ResourceGroup
	request["CurrentRoleId"] = strconv.Itoa(roleId)

	if strings.Trim(name, " ") != "" {
		request["Name"] = name
	}

	response, err := client.DoTeaRequest("POST", "dataworks-private-cloud", "2019-01-17", "ListCalcEnginesForAscm", "", nil, nil, request)
	addDebug("ListCalcEnginesForAscm", response, request)
	if err != nil {
		return nil, err
	}

	resp := &MaxComputeProject{}
	body, ok := response["Body"].(string)
	if !ok {
		return resp, errmsgs.WrapError(err)
	}
	err = json.Unmarshal([]byte(body), resp)
	if err != nil {
		return resp, errmsgs.WrapError(err)
	}

	if resp.TotalCount < 1 || resp.Code == "200" {
		return resp, errmsgs.WrapError(err)
	}
	return resp, nil
}

func (s *MaxcomputeService) DescribeMaxcomputeCu(name string) (object map[string]interface{}, err error) {
	request := make(map[string]interface{})
	request["RegionName"] = s.client.RegionId
	request["Product"] = "ascm"
	request["OrganizationId"] = s.client.Department
	request["ResourceGroupId"] = s.client.ResourceGroup
	request["Department"] = s.client.Department

	response, err := s.client.DoTeaRequest("POST", "ascm", "2019-05-10", "ListOdpsCus", "/ascm/manage/odps/list_cus", nil, nil, request)
	addDebug("ListOdpsCus", response, request)
	if err != nil {
		err = errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, name, "ListOdpsCus", errmsgs.AlibabacloudStackSdkGoERROR)
		return
	}

	if errmsgs.IsExpectedErrorCodes(fmt.Sprintf("%v", response["code"]), []string{"102", "403"}) {
		err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("MaxcomputeProject", name)), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
		return object, err
	}
	if fmt.Sprintf(`%v`, response["code"]) != "200" {
		err = errmsgs.Error("ListOdpsCus failed for " + response["asapiErrorMessage"].(string))
		return object, err
	}

	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, name, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *MaxcomputeService) DescribeMaxcomputeUser(name string) (response *OdpsUser, err error) {
	roleId, err := s.client.RoleIds()
	if err != nil {
		err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("ASCM User", "defaultRoleId")), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
		return nil, err
	}

	request := make(map[string]interface{})
	request["UserName"] = name
	request["x-acs-roleid"] = strconv.Itoa(roleId)

	responseData, err := s.client.DoTeaRequest("POST", "ascm", "2019-05-10", "GetOdpsUserList", "/ascm/manage/resource_mgmt/listOdpsUser", nil, nil, request)
	addDebug("GetOdpsUserList", responseData, request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"Error OdpsUser Not Found"}) {
			return nil, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return nil, err
	}

	resp := &OdpsUser{}
	body, ok := responseData["Body"].(string)
	if !ok {
		return resp, errmsgs.WrapError(err)
	}
	err = json.Unmarshal([]byte(body), resp)
	if err != nil {
		return resp, errmsgs.WrapError(err)
	}

	if len(resp.Data) < 1 || resp.Code == "200" {
		return resp, errmsgs.WrapError(err)
	}
	return resp, nil
}
