package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
)

type MaxcomputeService struct {
	client *connectivity.AlibabacloudStackClient
}

func (s *MaxcomputeService) DescribeMaxcomputeProject(id string) (object *MaxComputeProject, err error) {
	client := s.client
	request := client.NewCommonRequest("POST", "dataworks-private-cloud", "2019-01-17", "ListCalcEnginesForAscm", "")
	bresponse, err := client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		errmsg := ""
		errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_maxcompute_project", "ListCalcEnginesForAscm", errmsg)
	}

	response := &MaxComputeProjectDetailResponse{}
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		return nil, errmsgs.WrapError(err)
	}
	for _, v := range response.Data.CalcEngines {
		// EngineStatus == 0 is running   1 is deleted
		if fmt.Sprint(v.EngineId) == id && v.EngineStatus == 0 {
			return &v, nil
		}
	}
	return nil, errmsgs.Error(errmsgs.GetNotFoundMessage("Maxcompute Project", id))
}

func (s *MaxcomputeService) DescribeMaxcomputeProjectEngine(id string) (object *MaxComputeProjectEngineData, err error) {
	client := s.client
	request := client.NewCommonRequest("GET", "dataworks-private-cloud", "2019-01-17", "GetCalcEngineForAscm", "")
	request.QueryParams["EngineId"] = id
	bresponse, err := client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		errmsg := ""
		errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_maxcompute_project", "GetCalcEngineForAscm", errmsg)
	}
	response := MaxComputeProjectEngineDetailResponse{}
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		return nil, errmsgs.WrapError(err)
	}
	return &response.Data, nil
}

func (s *MaxcomputeService) ListOdpsEngineQuotaForAscm(id, name string) (object *OdpsEngineQuotaData, err error) {
	client := s.client
	request := client.NewCommonRequest("GET", "dataworks-private-cloud", "2019-01-17", "ListOdpsEngineQuotaForAscm", "")
	request.QueryParams = map[string]string{
		"Id":          id,
		"ProjectName": name,
	}
	bresponse, err := client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		errmsg := ""
		errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_maxcompute_project", "ListOdpsEngineQuotaForAscm", errmsg)
	}
	response := ListOdpsEngineQuotaForAscmResponse{}
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		return nil, errmsgs.WrapError(err)
	}
	return &response.Data[0], nil
}

func (s *MaxcomputeService) DescribeMaxProjectPropertiesForAscm(id, name string) (map[string]interface{}, error) {
	client := s.client
	request := client.NewCommonRequest("GET", "dataworks-private-cloud", "2019-01-17", "GetOdpsProjectPropertiesForAscm", "")
	request.QueryParams = map[string]string{
		"ProjectName":        name,
		"Properties":         "ENCRYPTION,odps.security.vpc.whitelist",
		"EngineId":           id,
		"DoReplaceTunnelIds": "true",
	}
	bresponse, err := client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		errmsg := ""
		errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_maxcompute_project", "GetOdpsProjectPropertiesForAscm", errmsg)
	}
	response := make(map[string]interface{})
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		return nil, errmsgs.WrapError(err)
	}
	v, err := jsonpath.Get("$.Data", response)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$", response)
	}
	return v.(map[string]interface{}), nil
}

func (s *MaxcomputeService) DescribeMaxcomputeQuota(cluster_id string) (object map[string]interface{}, err error) {
	request := map[string]interface{}{
		"Region":      s.client.RegionId,
		"Action":      "GetOdpsQuotaForAscm",
		"AccessKeyId": s.client.AccessKey,
		"Cluster":     cluster_id,
		"Project":     "odps",
	}

	response, err := s.client.DoTeaRequest("GET", "dataworks-private-cloud", "2019-01-17", "GetOdpsQuotaForAscm", "", nil, request, nil)
	addDebug("ListOdpsCusForAscm", response, request)
	if err != nil {
		err = errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_maxcompute_project", "GetOdpsQuotaForAscm", errmsgs.AlibabacloudStackSdkGoERROR)
		return
	}
	if fmt.Sprintf(`%v`, response["HttpStatusCode"]) != "200" {
		err = errmsgs.Error("ListOdpsCusForAscm failed for " + response["asapiErrorMessage"].(string))
		return object, err
	}

	v, err := jsonpath.Get("$.Data.data", response)
	if err != nil {
		return object, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, cluster_id, "$", response)
	}
	object = v.(map[string]interface{})
	return object, nil
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

func (s *MaxcomputeService) DescribeMaxcomputeUsers(organization_id string) (response []interface{}, err error) {
	request := make(map[string]interface{})
	request["Department"] = organization_id
	request["Region"] = s.client.RegionId
	request["Action"] = "GetOdpsUserList"
	request["AccessKeyId"] = s.client.AccessKey

	responseData, err := s.client.DoTeaRequest("GET", "ascm", "2019-05-10", "GetOdpsUserList", "", nil, request, nil)
	addDebug("GetOdpsUserList", responseData, request, request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"Error OdpsUser Not Found"}) {
			return nil, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return nil, err
	}

	datas, err := jsonpath.Get("$.data", responseData)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, organization_id, "$.data", response)
	}
	if datas != nil {
		return datas.([]interface{}), nil
	}
	return nil, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
}

func (s *MaxcomputeService) DescribeMaxcomputeUser(id string) (response map[string]interface{}, err error) {
	params := strings.Split(id, ":")
	request := make(map[string]interface{})
	request["Department"] = params[0]
	request["Id"] = params[1]
	request["Region"] = s.client.RegionId
	request["Action"] = "GetOdpsUserList"
	request["AccessKeyId"] = s.client.AccessKey

	responseData, err := s.client.DoTeaRequest("GET", "ascm", "2019-05-10", "GetOdpsUserList", "", nil, request, nil)
	addDebug("GetOdpsUserList", responseData, request, request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"Error OdpsUser Not Found"}) {
			return nil, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
		}
		return nil, err
	}

	datas, err := jsonpath.Get("$.data", responseData)
	if err != nil {
		return nil, errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, id, "$.data", response)
	}
	for _, v := range datas.([]interface{}) {
		data := v.(map[string]interface{})
		if fmt.Sprint(data["id"]) == params[1] {
			return data, nil
		}
	}
	return nil, errmsgs.WrapErrorf(err, errmsgs.NotFoundMsg, errmsgs.AlibabacloudStackSdkGoERROR)
}
