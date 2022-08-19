package apsarastack

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/PaesslerAG/jsonpath"

	"strconv"
	"strings"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
)

type MaxcomputeService struct {
	client *connectivity.ApsaraStackClient
}

func (s *MaxcomputeService) DescribeMaxcomputeProject(name string) (object *MaxComputeProject, err error) {
	client := s.client
	var requestInfo *ecs.Client

	request := requests.NewCommonRequest()
	if client.Config.Insecure {
		request.SetHTTPSInsecure(client.Config.Insecure)
	}

	roleId, err := client.RoleIds()
	if err != nil {
		err = WrapErrorf(Error(GetNotFoundMessage("ASCM User", "defaultRoleId")), NotFoundMsg, ProviderERROR)
		return nil, err
	}

	request.QueryParams = map[string]string{
		"Action":          "ListCalcEnginesForAscm",
		"ResourceGroupId": client.ResourceGroup,
		"Product":         "dataworks-private-cloud",
		"CalcEngineType":  "ODPS", // 固定值
		"OrganizationId":  client.Department,
		"Department":      client.Department,
		"ResourceGroup":   client.ResourceGroup,
		"RegionId":        client.RegionId,
		"CurrentRoleId":   strconv.Itoa(roleId),
	}

	if strings.Trim(name, " ") != "" {
		request.QueryParams["Name"] = name
	}

	request.Method = "POST"
	request.Product = "dataworks-private-cloud"
	request.Version = "2019-01-17"
	request.ServiceCode = "dataworks-private-cloud"
	request.Domain = client.Domain
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "ListCalcEnginesForAscm"
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	log.Printf("response of raw create maxcomputecluster is : %s", raw)

	if err != nil {
		return nil, WrapErrorf(err, DefaultErrorMsg, "apsarastack_maxcompute_project", "List", raw)
	}

	addDebug("MaxcomputeProjectCreate", raw, requestInfo, request)

	bresponse, _ := raw.(*responses.CommonResponse)
	if bresponse.GetHttpStatus() != 200 {
		return nil, WrapErrorf(err, DefaultErrorMsg, "apsarastack_maxcompute_project", "List", ApsaraStackSdkGoERROR)
	}

	var resp = &MaxComputeProject{}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), resp)
	if err != nil {
		return resp, WrapError(err)
	}

	if resp.TotalCount < 1 || resp.Code == "200" {
		return resp, WrapError(err)
	}
	return resp, nil
}

func (s *MaxcomputeService) DescribeMaxcomputeCu(name string) (object map[string]interface{}, err error) {
	var response map[string]interface{}
	conn, err := s.client.NewAscmClient()
	if err != nil {
		return nil, WrapError(err)
	}
	action := "ListOdpsCus"
	request := map[string]interface{}{
		"RegionName": s.client.RegionId,
		//"Type":            "cuName",
		//"CuName":          name,
		"Product":         "ascm",
		"OrganizationId":  s.client.Department,
		"ResourceGroupId": s.client.ResourceGroup,
		"Department":      s.client.Department,
	}
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequestWithOrg(StringPointer(action), nil, StringPointer("POST"), StringPointer("2019-05-10"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		err = WrapErrorf(err, DefaultErrorMsg, name, action, ApsaraStackSdkGoERROR)
		return
	}
	addDebug(action, response, request)
	if IsExpectedErrorCodes(fmt.Sprintf("%v", response["code"]), []string{"102", "403"}) {
		err = WrapErrorf(Error(GetNotFoundMessage("MaxcomputeProject", name)), NotFoundMsg, ProviderERROR)
		return object, err
	}
	if fmt.Sprintf(`%v`, response["code"]) != "200" {
		err = Error("ListOdpsCus failed for " + response["asapiErrorMessage"].(string))
		return object, err
	}
	v, err := jsonpath.Get("$", response)
	if err != nil {
		return object, WrapErrorf(err, FailedGetAttributeMsg, name, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
}

func (s *MaxcomputeService) DescribeMaxcomputeUser(name string) (response *OdpsUser, err error) {
	var requestInfo *ecs.Client
	request := requests.NewCommonRequest()
	if s.client.Config.Insecure {
		request.SetHTTPSInsecure(s.client.Config.Insecure)
	}

	roleId, err := s.client.RoleIds()
	if err != nil {
		err = WrapErrorf(Error(GetNotFoundMessage("ASCM User", "defaultRoleId")), NotFoundMsg, ProviderERROR)
		return nil, err
	}

	request.Method = "POST"
	request.Product = "ascm"
	request.Version = "2019-05-10"
	if strings.ToLower(s.client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "GetOdpsUserList"
	request.Headers = map[string]string{
		"RegionId":              s.client.RegionId,
		"x-acs-roleid":          strconv.Itoa(roleId),
		"x-acs-resourcegroupid": s.client.ResourceGroup,
		"x-acs-regionid":        s.client.RegionId,
		"x-acs-organizationid":  s.client.Department,
	}
	request.QueryParams = map[string]string{
		"RegionId":        s.client.RegionId,
		"UserName":        name,
		"Product":         "ascm",
		"OrganizationId":  s.client.Department,
		"ResourceGroupId": s.client.ResourceGroup,
	}
	raw, err := s.client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	if err != nil {
		if IsExpectedErrors(err, []string{"Error OdpsUser Not Found"}) {
			return nil, WrapErrorf(err, NotFoundMsg, ApsaraStackSdkGoERROR)
		}
		return nil, WrapErrorf(err, DefaultErrorMsg, name, "GetOdpsUserList", ApsaraStackSdkGoERROR)

	}
	addDebug("GetOdpsUserList", raw, requestInfo, request)

	var resp = &OdpsUser{}

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
