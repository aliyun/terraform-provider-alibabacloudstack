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

func (s *MaxcomputeService) DescribeMaxcomputeCu(id string) (object map[string]interface{}, err error) {
	request := map[string]interface{}{
		"Region":      s.client.RegionId,
		"Action":      "ListOdpsCusForAscm",
		"AccessKeyId": s.client.AccessKey,
		"CuId":        id,
	}

	response, err := s.client.DoTeaRequest("GET", "dataworks-private-cloud", "2019-01-17", "ListOdpsCusForAscm", "", nil, request, nil)
	addDebug("ListOdpsCusForAscm", response, request)
	if err != nil {
		err = errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_maxcompute_cu", "ListOdpsCusForAscm", errmsgs.AlibabacloudStackSdkGoERROR)
		return
	}
	if fmt.Sprintf(`%v`, response["HttpStatusCode"]) != "200" {
		err = errmsgs.Error("ListOdpsCusForAscm failed for " + response["asapiErrorMessage"].(string))
		return object, err
	}

	v, err := jsonpath.Get("$.Data.data", response)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$", response)
	}
	datas := v.([]interface{})
	for _, element := range datas {
		if element.(map[string]interface{})["id"].(string) == id {
			return element.(map[string]interface{}), nil
		}
	}
	return nil, errmsgs.Error(errmsgs.GetNotFoundMessage("Maxcompute", id))
}

func (s *MaxcomputeService) DescribeMaxcomputeUser(name string) (response *OdpsUser, err error) {
	roleId, err := s.client.RoleIds()
	if err != nil {
		err = errmsgs.WrapErrorf(errmsgs.Error(errmsgs.GetNotFoundMessage("ASCM User", "defaultRoleId")), errmsgs.NotFoundMsg, errmsgs.ProviderERROR)
		return nil, err
	}

	request := make(map[string]interface{})
	request["UserName"] = name
	request["Region"] = s.client.RegionId
	request["Action"] = "GetOdpsUserList"
	request["AccessKeyId"] = s.client.AccessKey
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
